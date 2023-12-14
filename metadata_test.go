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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_Metadata(t *testing.T) {
	t.Run(
		"Successful request and unmarshaling of response",
		func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/v1/metadata/108923398" {
						response := []Metadatum{
							{
								Hash: "b64602eebf602e8bbce198e2a1d6bbb2a109ae87fa5316135d217110d6d94649",
								Raw:  "a11902a2a1636d736781781c4d696e737761703a205377617020457861637420496e204f72646572",
								Schema: json.RawMessage(`{"exampleKey":"exampleValue"}`),
							},
						}
						respBody, _ := json.Marshal(response)
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write(respBody)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				}),
			)
			defer server.Close()

			client := New(WithEndpoint(server.URL))
			metadataResponse, err := client.Metadata(
				context.Background(),
				108923398,
				"",
			)
			if err != nil {
				t.Fatalf("Expected no error, got %s", err)
			}
			expectedList := []Metadatum{
				{
					Hash:   "b64602eebf602e8bbce198e2a1d6bbb2a109ae87fa5316135d217110d6d94649",
					Raw:    "a11902a2a1636d736781781c4d696e737761703a205377617020457861637420496e204f72646572",
					Schema: json.RawMessage(`{"exampleKey":"exampleValue"}`),
				},
			}

			if !reflect.DeepEqual(metadataResponse, expectedList) {
				t.Errorf(
					"Expected response %v, got %v",
					expectedList,
					metadataResponse,
				)
			}
		},
	)

	t.Run(
		"Successful request returning empty",
		func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					response := []Metadatum{}
					respBody, _ := json.Marshal(response)
					_, _ = w.Write(respBody)
				}),
			)
			defer server.Close()

			client := New(WithEndpoint(server.URL))
			metadataResponse, err := client.Metadata(
				context.Background(),
				108923398,
				"",
			)
			if err != nil {
				t.Fatalf("Expected no error, got %s", err)
			}
			if len(metadataResponse) != 0 {
				t.Errorf("Expected empty response, got %v", metadataResponse)
			}
		},
	)
}
