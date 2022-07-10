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
	"net/url"
	"time"

	"github.com/SundaeSwap-finance/ogmigo"
	"github.com/SundaeSwap-finance/ogmigo/ouroboros/chainsync"
)

type options struct {
	spent     bool
	unspent   bool
	pattern   string
	policyId  string
	assetName string
}

func (o options) apply(url *url.URL) {
	qs := ""
	if o.spent && !o.unspent {
		qs += "spent"
	} else if o.unspent && !o.spent {
		qs += "unspent"
	}

	if o.policyId != "" {
		if qs != "" {
			qs += "&"
		}
		qs += fmt.Sprintf("policy_id=%v", o.policyId)
		if o.assetName != "" {
			qs += fmt.Sprintf("&asset_name=%v", o.assetName)
		}
	}
	if o.pattern != "" {
		url.Path += fmt.Sprintf("/%v", o.pattern)
	}
	url.RawQuery = qs
}

type Filter func(*options)

func (c *Client) Matches(ctx context.Context, filters ...Filter) (matches []Match, err error) {
	start := time.Now()
	defer func() {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		c.options.logger.Info("Matches() finished",
			ogmigo.KV("duration", time.Since(start).Round(time.Millisecond).String()),
			ogmigo.KV("matched", fmt.Sprintf("%v", len(matches))),
			ogmigo.KV("err", errStr),
		)
	}()

	url, err := url.Parse(c.options.endpoint)
	url.Path = "/v1/matches"
	if err != nil {
		return nil, fmt.Errorf("unable to parse endpoint %v: %w", c.options.endpoint, err)
	}

	o := options{}
	for _, f := range filters {
		f(&o)
	}
	o.apply(url)

	c.logger.Debug("finding matches", ogmigo.KV("url", url.String()))

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch matches: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	matches = []Match{}
	if err := json.Unmarshal(body, &matches); err != nil {
		return nil, fmt.Errorf("unable to parse body %v: %w", string(body), err)
	}
	return matches, nil
}

func All() Filter {
	return func(o *options) {
		o.spent = true
		o.unspent = true
	}
}

func OnlySpent() Filter {
	return func(o *options) {
		o.spent = true
		o.unspent = false
	}
}

func OnlyUnspent() Filter {
	return func(o *options) {
		o.unspent = true
		o.spent = false
	}
}

func Pattern(pattern string) Filter {
	return func(o *options) {
		o.pattern = pattern
	}
}

func PolicyID(policyId string) Filter {
	return func(o *options) {
		o.policyId = policyId
	}
}

func AssetID(assetID chainsync.AssetID) Filter {
	return func(o *options) {
		o.policyId = assetID.PolicyID()
		o.assetName = assetID.AssetName()
	}
}
