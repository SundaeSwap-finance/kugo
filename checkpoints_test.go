// Copyright 2022 SundaeSwap Labs, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:Licensed under the MIT License;
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://opensource.org/licenses/MIT
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package kugo

import (
	"context"
	"testing"

	"github.com/tj/assert"
)

func Test_Checkpoints(t *testing.T) {
	// Integration test that relies on a local instance of kupo, so skip for now
	t.SkipNow()
	c := New(WithEndpoint("http://localhost:1442"))
	points, err := c.Checkpoints(context.Background(), BySlot(51540727))
	assert.Nil(t, err)
	assert.Len(t, points, 1)
	point := points[0]

	assert.Equal(t, "fe5f9af58ab0511a77524f4d2a0b930213b3bb1353e11e3d69e83129b9fbe65a", point.HeaderHash)
	assert.Equal(t, 51540722, point.SlotNo)
}
