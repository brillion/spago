// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sinusoidalpositionalencoder

import (
	"encoding/gob"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/encoding/pe"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
)

var (
	_ nn.Model = &SinusoidalPositionalEncoder{}
)

// Config provides configuration settings for a SinusoidalPositionalEncoder Model.
type Config struct {
	NumEmbeddings int
	EmbeddingDim  int
}

// SinusoidalPositionalEncoder contains positional embeddings fine-tuned during
// the training phase.
type SinusoidalPositionalEncoder struct {
	nn.BaseModel
	Config   Config
	Delegate *pe.SinusoidalPositionalEncoder
}

func init() {
	gob.Register(&SinusoidalPositionalEncoder{})
}

// New returns a new SinusoidalPositionalEncoder.
func New(config Config) *SinusoidalPositionalEncoder {
	return &SinusoidalPositionalEncoder{
		Config:   config,
		Delegate: pe.NewSinusoidalPositionalEncoder(config.EmbeddingDim, config.NumEmbeddings),
	}
}

// Encode performs the forward step for each input and returns the result.
func (m *SinusoidalPositionalEncoder) Encode(positions []int) []ag.Node {
	embeddings := make([]ag.Node, len(positions))
	for i, vector := range m.Delegate.Encode(positions...) {
		embeddings[i] = m.Graph().NewVariable(vector, false)
	}
	return embeddings
}
