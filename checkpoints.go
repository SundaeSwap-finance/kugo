// Copyright 2022 SundaeSwap Labs, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software
// is furnished to do so, subject to the following conditions:
//
// Licensed under the MIT License;
// You may not use this file except in compliance with the License.
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
	"net/url"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6"
)

type checkpointsOptions struct {
	singular bool
	slot     uint64
}

func (c checkpointsOptions) apply(url *url.URL) {
	if c.slot != 0 {
		url.Path += fmt.Sprintf("/%v", c.slot)
	}
}

type CheckpointsFilter struct {
	before func(*checkpointsOptions)
	after  func(points []Point) []Point
}

func (c *Client) Checkpoints(ctx context.Context, filters ...CheckpointsFilter) (points []Point, err error) {
	start := time.Now()
	defer func() {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		c.options.logger.Info("Checkpoints() finished",
			ogmigo.KV("duration", time.Since(start).Round(time.Millisecond).String()),
			ogmigo.KV("err", errStr),
		)
	}()

	endpoint, err := url.Parse(c.options.endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse endpoint ")
	}

	endpoint.Path = "/v1/checkpoints"

	o := checkpointsOptions{}
	for _, f := range filters {
		if f.before != nil {
			f.before(&o)
		}
	}
	o.apply(endpoint)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Close = true
	req = req.WithContext(ctx)

	client := &http.Client{
		Timeout: c.options.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve checkpoint by slot: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got unexpected response: %v: %v", resp.StatusCode, string(body))
	}

	if o.singular {
		var point Point
		if err := json.Unmarshal(body, &point); err != nil {
			return nil, fmt.Errorf("error parsing response %v: %w", string(body), err)
		}
		points = []Point{point}
	} else {
		if err := json.Unmarshal(body, &points); err != nil {
			return nil, fmt.Errorf("error parsing response %v: %w", string(body), err)
		}
	}
	for _, f := range filters {
		if f.after != nil {
			points = f.after(points)
		}
	}

	return points, nil
}

// Return a recent sampling of kupo checkpoints
// NOTE: equivalent to providing no filters, but useful for documenting your purpose
func Recent() CheckpointsFilter {
	return CheckpointsFilter{
		before: nil,
		after:  nil,
	}
}

func Latest() CheckpointsFilter {
	return CheckpointsFilter{
		before: nil,
		after: func(points []Point) []Point {
			return []Point{points[0]}
		},
	}
}

func BySlot(slot uint64) CheckpointsFilter {
	return CheckpointsFilter{
		before: func(co *checkpointsOptions) {
			co.singular = true
			co.slot = slot
		},
		after: nil,
	}
}
