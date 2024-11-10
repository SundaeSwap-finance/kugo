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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6"
	"golang.org/x/crypto/blake2b"
)

type Script struct {
	Language string `json:"language"`
	Script   string `json:"script"`
}

func (s Script) Hash() []byte {
	scriptBytes, _ := hex.DecodeString(s.Script)
	blake, _ := blake2b.New(224/8, nil)
	switch s.Language {
	case "plutus:v1":
		blake.Write([]byte{0x01})
	case "plutus:v2":
		blake.Write([]byte{0x02})
	}
	blake.Write(scriptBytes)
	hashBytes := blake.Sum(nil)
	return hashBytes[:]
}

func (c *Client) Script(ctx context.Context, scriptHash string) (script *Script, err error) {
	start := time.Now()
	defer func() {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		c.options.logger.Info("Script() finished",
			ogmigo.KV("duration", time.Since(start).Round(time.Millisecond).String()),
			ogmigo.KV("err", errStr),
		)
	}()

	url, err := url.Parse(c.options.endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to parse endpoint %v: %w", c.options.endpoint, err)
	}
	url.Path = "/v1/script/" + scriptHash

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
		return nil, fmt.Errorf("unable to fetch script: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("failed with a nil response")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	response := &Script{}
	if err := json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("unable to parse body %s: %w", body, err)
	}
	return response, nil
}
