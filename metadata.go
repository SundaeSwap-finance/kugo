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
	"strconv"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6"
)

type Metadatum struct {
	Hash   string
	Raw    string
	Schema json.RawMessage
}

func (c *Client) Metadata(ctx context.Context, slotNo int, txId string) (metadatum []Metadatum, err error) {
	start := time.Now()
	defer func() {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		c.options.logger.Info("Metadata() finished",
			ogmigo.KV("duration", time.Since(start).Round(time.Millisecond).String()),
			ogmigo.KV("err", errStr),
		)
	}()

	url, err := url.Parse(c.options.endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse endpoint %v: %w", c.options.endpoint, err)
	}
	url.Path = "/v1/metadata/" + strconv.Itoa(slotNo)

	if txId != "" {
		url.Path += "?transaction_id=" + txId
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Close = true
	req = req.WithContext(ctx)

	client := http.DefaultClient
	client.Timeout = 5 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch metadata: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	response := []Metadatum{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unable to parse body %v: %w", string(body), err)
	}
	return response, nil
}
