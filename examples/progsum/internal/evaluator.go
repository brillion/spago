// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/gosuri/uiprogress"
	"github.com/nlpodyssey/spago/pkg/mat/f64utils"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/ml/stats"
)

type Evaluator struct {
	model nn.Model
}

func NewEvaluator(model nn.Model) *Evaluator {
	return &Evaluator{
		model: model,
	}
}

// Predict performs the forward pass
func (t *Evaluator) Predict(example Sequence) int {
	g := ag.NewGraph()
	defer g.Clear()
	xs := make([]ag.Node, len(example))
	for i, x := range example {
		xs[i] = g.NewScalar(x.Input)
	}
	ys := t.model.NewProc(nn.Context{Graph: g, Mode: nn.Inference}).Forward(xs...)
	return f64utils.ArgMax(ys[len(example)-1].Value().Data())
}

func (t *Evaluator) Evaluate(dataset []Sequence) *stats.ClassMetrics {
	uip := uiprogress.New()
	bar := newTestBar(uip, dataset)
	uip.Start()
	defer uip.Stop()
	counter := stats.NewMetricCounter()
	for i := 0; i < len(dataset); i++ {
		sequence := dataset[i]
		if t.Predict(sequence) == sequence[len(sequence)-1].Target {
			counter.IncTruePos()
		} else {
			counter.IncFalsePos()
		}
		bar.Incr()
	}
	return counter
}

func newTestBar(p *uiprogress.Progress, dataset []Sequence) *uiprogress.Bar {
	bar := p.AddBar(len(dataset))
	bar.AppendCompleted().PrependElapsed()
	return bar
}
