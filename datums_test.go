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

func TestClient_Datum(t *testing.T) {
	t.Run(
		"Successful request and unmarshaling of response",
		func(t *testing.T) {
			t.Parallel()

			type DatumResponse struct {
				Datum string
			}
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/v1/datums/34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059" {
						response := DatumResponse{
							Datum: "d87980",
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
			datumResponse, err := client.Datum(
				context.Background(),
				"34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059",
			)
			expectedResponse := "d87980"
			if err != nil {
				t.Fatalf("Expected no error, got %s", err)
			}
			if !reflect.DeepEqual(datumResponse, expectedResponse) {
				t.Errorf(
					"Expected response %v, got %v",
					expectedResponse,
					datumResponse,
				)
			}
		},
	)

	t.Run(
		"Successful request returning empty",
		func(t *testing.T) {
			t.Parallel()

			type DatumResponse struct {
				Datum string
			}
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					response := DatumResponse{}
					respBody, _ := json.Marshal(response)
					_, _ = w.Write(respBody)
				}),
			)
			defer server.Close()

			client := New(WithEndpoint(server.URL))
			datumResponse, err := client.Datum(
				context.Background(),
				"34215ad90b1ade84f5b4fe3c0a16cb3afeae468210535e0305efd93931f35059",
			)
			if err != nil {
				t.Fatalf("Expected no error, got %s", err)
			}
			if datumResponse != "" {
				t.Errorf("Expected empty response, got %v", datumResponse)
			}
		},
	)
}
