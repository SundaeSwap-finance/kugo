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
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
)

type matchesOptions struct {
	spent     bool
	unspent   bool
	pattern   string
	policyId  string
	assetName string
	txHash    string
	txIx      *int // This is a pointer so we can distinguish between null and 0
	// Pagination properties
	created_before uint64
	spent_before   uint64
	created_after  uint64
	spent_after    uint64
}

func (o matchesOptions) apply(url *url.URL) {
	qs := ""
	// Handle the mutually exclusive spent/unspent filters
	if o.spent && !o.unspent {
		qs += "spent"
	} else if o.unspent && !o.spent {
		qs += "unspent"
	}

	// Handle the pagination query params
	if o.created_before != 0 {
		if qs != "" {
			qs += "&"
		}
		qs += fmt.Sprintf("created_before=%v", o.created_before)
	}
	if o.created_after != 0 {
		if qs != "" {
			qs += "&"
		}
		qs += fmt.Sprintf("created_after=%v", o.created_after)
	}
	if o.spent_before != 0 {
		if qs != "" {
			qs += "&"
		}
		qs += fmt.Sprintf("spent_before=%v", o.spent_before)
	}
	if o.spent_after != 0 {
		if qs != "" {
			qs += "&"
		}
		qs += fmt.Sprintf("spent_after=%v", o.spent_after)
	}

	// Handle txHash / index parameters
	if o.txHash != "" {
		if qs != "" {
			qs += "&"
		}

		// If another pattern is specified, that pattern takes precedence
		// and we need to specify these as query string parameters
		if o.pattern != "" {
			qs += fmt.Sprintf("transaction_id=%v", o.txHash)
			if o.txIx != nil {
				qs += fmt.Sprintf("&output_index=%v", *o.txIx)
			}
		} else {
			// NOTE(pi): kugo uses 'idx@txHash' because # isn't safely url-encodable
			if o.txIx == nil {
				url.Path += fmt.Sprintf("/*@%v", o.txHash)
			} else {
				url.Path += fmt.Sprintf("/%v@%v", *o.txIx, o.txHash)
			}
		}
	}

	// Handle policy ID filters
	if o.policyId != "" {
		if qs != "" {
			qs += "&"
		}
		// If another pattern is specified, that pattern takes precedence
		// and we need to specify these as query string parameters
		if o.pattern != "" {
			qs += fmt.Sprintf("policy_id=%v", o.policyId)
			if o.assetName != "" {
				qs += fmt.Sprintf("&asset_name=%v", o.assetName)
			}
		} else {
			if o.assetName != "" {
				url.Path += fmt.Sprintf("/%v.%v", o.policyId, o.assetName)
			} else {
				url.Path += fmt.Sprintf("/%v.*", o.policyId)
			}
		}
	}

	// Handle explicit patterns
	if o.pattern != "" {
		url.Path += fmt.Sprintf("/%v", o.pattern)
	}

	url.RawQuery = qs
}

type MatchesFilter func(*matchesOptions)

func (c *Client) Matches(
	ctx context.Context,
	filters ...MatchesFilter,
) (matches []Match, err error) {
	start := time.Now()
	defer func() {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		c.options.logger.Debug(
			"Matches() finished",
			ogmigo.KV(
				"duration",
				time.Since(start).Round(time.Millisecond).String(),
			),
			ogmigo.KV("matched", fmt.Sprintf("%v", len(matches))),
			ogmigo.KV("err", errStr),
		)
	}()

	url, err := url.Parse(c.options.endpoint)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to parse endpoint %v: %w",
			c.options.endpoint,
			err,
		)
	}
	url.Path = "/v1/matches"

	o := matchesOptions{}
	for _, f := range filters {
		f(&o)
	}
	o.apply(url)

	c.logger.Debug("finding matches", ogmigo.KV("url", url.String()))

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Close = true
	req = req.WithContext(ctx)

	client := &http.Client{
		Timeout: c.options.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch matches: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("failed with a nil response")
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

func All() MatchesFilter {
	return func(o *matchesOptions) {
		o.spent = true
		o.unspent = true
	}
}

func OnlySpent() MatchesFilter {
	return func(o *matchesOptions) {
		o.spent = true
		o.unspent = false
	}
}

func OnlyUnspent() MatchesFilter {
	return func(o *matchesOptions) {
		o.unspent = true
		o.spent = false
	}
}

func Pattern(pattern string) MatchesFilter {
	return func(o *matchesOptions) {
		o.pattern = pattern
	}
}

func Address(address string) MatchesFilter {
	return func(o *matchesOptions) {
		o.pattern = address
	}
}

func Transaction(txHash string) MatchesFilter {
	return func(o *matchesOptions) {
		o.txHash = txHash
	}
}

// NOTE(pi): chainsync.TxID is named poorly, and we plan to rename it at some point
func TxOut(txOutId chainsync.TxID) MatchesFilter {
	return func(o *matchesOptions) {
		o.txHash = txOutId.TxHash()
		// store locally so we can take the address
		idx := txOutId.Index()
		o.txIx = &idx
	}
}

func PolicyID(policyId string) MatchesFilter {
	return func(o *matchesOptions) {
		o.policyId = policyId
	}
}

func AssetID(assetID shared.AssetID) MatchesFilter {
	return func(o *matchesOptions) {
		o.policyId = assetID.PolicyID()
		o.assetName = assetID.AssetName()
	}
}

func Overlapping(slot uint64) MatchesFilter {
	return func(o *matchesOptions) {
		o.created_before = slot
		o.spent_after = slot
	}
}

func CreatedBefore(slot uint64) MatchesFilter {
	return func(o *matchesOptions) {
		o.created_before = slot
	}
}
func CreatedAfter(slot uint64) MatchesFilter {
	return func(o *matchesOptions) {
		o.created_after = slot
	}
}
func SpentBefore(slot uint64) MatchesFilter {
	return func(o *matchesOptions) {
		o.spent_before = slot
	}
}
func SpentAfter(slot uint64) MatchesFilter {
	return func(o *matchesOptions) {
		o.spent_after = slot
	}
}
