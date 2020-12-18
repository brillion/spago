// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
)

var _ Function = &Reshape{}

type Reshape struct {
	x    Operand
	rows int
	cols int
}

// NewReshape returns a new Reshape Function.
func NewReshape(x Operand, r, c int) *Reshape {
	return &Reshape{x: x, rows: r, cols: c}
}

// Forward computes the output of the node.
func (r *Reshape) Forward() mat.Matrix {
	if r.x.Value().Size() != r.rows*r.cols {
		panic("fn: incompatible sizes")
	}
	return r.x.Value().Reshape(r.rows, r.cols)
}

func (r *Reshape) Backward(gy mat.Matrix) {
	if gy.Columns() != r.cols && gy.Rows() != r.rows {
		panic("fn: matrices with not compatible size")
	}
	if r.x.RequiresGrad() {
		gx := gy.Reshape(r.x.Value().Dims())
		defer mat.ReleaseDense(gx.(*mat.Dense))
		r.x.PropagateGrad(gx)
	}
}
