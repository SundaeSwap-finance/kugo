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
	"fmt"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/ouroboros/chainsync"
	"github.com/tj/assert"
)

func Test(t *testing.T) {
	t.SkipNow()
	c := New(WithEndpoint("http://localhost:1442"))
	matches, err := c.Matches(context.Background(),
		OnlyUnspent(),
		AssetID(chainsync.AssetID("4fc16c94d066e949e771c5581235f8090ad6aaffaf373a426445ca51.73636f6f70209a0a")),
		Pattern("addr_test1qpluezahtqdtwg4f7qewdvjvz806hsatqwr4u04yzcrk2m7pucvj7jyhq97rca9m0wul2fu3qnsayxvqdwlda8wngurqgyfepe"),
	)
	assert.Nil(t, err)
	assert.NotZero(t, len(matches))

	fmt.Printf("Matches: %v\n", len(matches))
}
