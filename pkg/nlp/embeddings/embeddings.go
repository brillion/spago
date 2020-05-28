// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embeddings

import (
	"bytes"
	"github.com/nlpodyssey/spago/pkg/mat"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/utils/kvdb"
	"log"
	"strings"
	"sync"
)

var (
	_ nn.Model     = &Model{}
	_ nn.Processor = &Processor{}
)

var allModels []*Model

// TODO: add dedicated embeddings for out-of-vocabulary (OOV) words and other special words
type Model struct {
	Config
	storage        kvdb.KeyValueDB
	mu             sync.Mutex
	UsedEmbeddings map[string]*nn.Param `type:"weights"`
	ZeroEmbedding  *nn.Param            `type:"weights"`
}

// TODO: add Dropout
type Config struct {
	// Size of the embedding vectors.
	Size int
	// Whether to return the `ZeroEmbedding` in case the word doesn't exist in the embeddings map.
	// If it is false, nil is returned instead, so the caller has more responsibility but also more control.
	UseZeroEmbedding bool
	// The path to DB on the drive
	DBPath string
	// Whether to use the map in read-only mode (embeddings are not updated during training).
	ReadOnly bool
	// Whether to force the deletion of any existing DB to start with an empty embeddings map.
	ForceNewDB bool
}

// New returns a new embedding model.
func New(config Config) *Model {
	m := &Model{
		Config: config,
		storage: kvdb.NewDefaultKeyValueDB(kvdb.Config{
			Path:     config.DBPath,
			ReadOnly: config.ReadOnly,
			ForceNew: config.ForceNewDB,
		}),
		UsedEmbeddings: map[string]*nn.Param{},
		ZeroEmbedding:  nn.NewParam(mat.NewEmptyVecDense(config.Size)),
	}
	nn.RequiresGrad(false)(m.ZeroEmbedding)
	allModels = append(allModels, m)
	return m
}

// Close closes the DB underlying the model of the embeddings map.
// It automatically clears the cache.
func (m *Model) Close() {
	_ = m.storage.Close() // explicitly ignore errors here
	m.ClearUsedEmbeddings()
}

// ClearUsedEmbeddings clears the cache of the used embeddings.
// Beware of any external references to the values of m.UsedEmbeddings. These are weak references!
func (m *Model) ClearUsedEmbeddings() {
	m.mu.Lock()
	for _, embedding := range m.UsedEmbeddings {
		mat.ReleaseDense(embedding.Value().(*mat.Dense))
	}
	m.UsedEmbeddings = map[string]*nn.Param{}
	m.mu.Unlock()
}

// Close closes the DBs underlying all instantiated embeddings models.
// It automatically clears the caches.
func Close() {
	for _, model := range allModels {
		model.Close()
	}
}

// ClearUsedEmbeddings clears the cache of the used embeddings of all instantiated embeddings models.
// Beware of any external references to the values of m.UsedEmbeddings. These are weak references!
func ClearUsedEmbeddings() {
	for _, model := range allModels {
		model.ClearUsedEmbeddings()
	}
}

func (m *Model) Count() int {
	keys, err := m.storage.Keys()
	if err != nil {
		log.Fatal(err)
	}
	return len(keys)
}

// SetEmbeddings inserts a new word embeddings.
// If the word is already on the map, overwrites the existing value with the new one.
func (m *Model) SetEmbedding(word string, value *mat.Dense) {
	if m.ReadOnly {
		log.Fatal("embedding: set operation not permitted in read-only mode")
	}
	embedding := nn.NewParam(value)
	embedding.SetPayload(nn.NewEmptySupport())
	var buf bytes.Buffer
	if _, err := (&nn.ParamSerializer{Param: embedding}).Serialize(&buf); err != nil {
		log.Fatal(err)
	}
	if err := m.storage.Put([]byte(word), buf.Bytes()); err != nil {
		log.Fatal(err)
	}
}

// GetEmbedding returns the parameter (the word embedding) associated with the given word.
// It first looks for the exact correspondence of the word. If there is no match, it tries the word lowercase.
//
// The returned embedding is also cached in m.UsedEmbeddings for two reasons:
//     - to allow a faster recovery;
//     - to keep track of used embeddings, should they be optimized.
//
// If no embedding is found, nil is returned.
// It panics in case of storage errors.
func (m *Model) GetEmbedding(word string) *nn.Param {
	if found := m.getEmbedding(word); found != nil {
		return found
	}
	if found := m.getEmbedding(strings.ToLower(word)); found != nil {
		return found
	}
	return nil
}

// getEmbedding returns the parameter (the word embedding) associated with the given word (exact correspondence).
// The returned embedding is also cached in m.UsedEmbeddings for two reasons:
//     - to allow a faster recovery;
//     - to keep track of used embeddings, should they be optimized.
// If no embedding is found, nil or the `ZeroEmbedding` is returned, depending on the configuration.
// It panics in case of storage errors.
func (m *Model) getEmbedding(word string) *nn.Param {
	if embedding, ok := m.UsedEmbeddings[word]; ok {
		return embedding
	}
	data, ok, err := m.storage.Get([]byte(word))
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		return nil // embedding not found
	}
	embedding := nn.NewParam(nil, nn.SetStorage(m.storage))
	if _, err := (&nn.ParamSerializer{Param: embedding}).Deserialize(bytes.NewReader(data)); err != nil {
		log.Fatal(err)
	}
	if m.ReadOnly {
		nn.RequiresGrad(false)(embedding)
	}
	embedding.SetName(word)
	m.mu.Lock()
	m.UsedEmbeddings[word] = embedding // important
	m.mu.Unlock()
	return embedding
}

type Processor struct {
	nn.BaseProcessor
	zeroEmbedding ag.Node
}

func (m *Model) NewProc(g *ag.Graph) nn.Processor {
	var zeroEmbedding ag.Node = nil
	if m.UseZeroEmbedding {
		zeroEmbedding = g.NewWrap(m.ZeroEmbedding)
	}
	return &Processor{
		BaseProcessor: nn.BaseProcessor{
			Model:             m,
			Mode:              nn.Training,
			Graph:             g,
			FullSeqProcessing: false,
		},
		zeroEmbedding: zeroEmbedding, // it can be nil
	}
}

// Encodes returns the embeddings associated with the input words.
// The embeddings are returned as Node(s) already inserted in the graph.
// To words that have no embeddings, the corresponding nodes are nil.
func (p *Processor) Encode(words []string) []ag.Node {
	encoding := make([]ag.Node, len(words))
	cache := make(map[string]ag.Node) // be smart, don't create two nodes for the same word!
	for i, word := range words {
		if item, ok := cache[word]; ok {
			encoding[i] = item
		} else {
			embedding := p.getEmbedding(word)
			encoding[i], cache[word] = embedding, embedding
		}
	}
	return encoding
}

func (p *Processor) getEmbedding(words string) ag.Node {
	model := p.Model.(*Model)
	switch param := model.GetEmbedding(words); {
	case param == nil:
		return p.zeroEmbedding // it can be nil
	default:
		return p.Graph.NewWrap(param)
	}
}

func (p *Processor) Forward(_ ...ag.Node) []ag.Node {
	panic("embeddings: p.Forward() not implemented. Use p.Encode() instead.")
}
