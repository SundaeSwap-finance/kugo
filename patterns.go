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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SundaeSwap-finance/ogmigo"
)

func (c *Client) Patterns(ctx context.Context) (matches []string, err error) {
	start := time.Now()
	defer func() {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		c.options.logger.Info("Patterns() finished",
			ogmigo.KV("duration", time.Since(start).Round(time.Millisecond).String()),
			ogmigo.KV("matched", fmt.Sprintf("%v", len(matches))),
			ogmigo.KV("err", errStr),
		)
	}()

	req, err := http.NewRequest("GET", fmt.Sprintf("%v/v1/patterns", c.options.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Close = true
	req = req.WithContext(ctx)

	client := http.DefaultClient
	client.Timeout = 5 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve patterns: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	matches = []string{}
	if err := json.Unmarshal(body, &matches); err != nil {
		return nil, fmt.Errorf("error parsing response %v: %w", string(body), err)
	}
	return matches, nil
}
