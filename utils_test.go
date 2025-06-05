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
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type DatumResponse struct {
	Datum string
}

type MockKugoServer struct {
	datums   map[string]DatumResponse
	matches  map[string][]Match
	metadata map[int]map[string][]Metadatum
	patterns []string
	scripts  map[string]Script
}

func NewMockServer() *MockKugoServer {
	return &MockKugoServer{}
}

func (m *MockKugoServer) AddScripts(script ...Script) *MockKugoServer {
	if m.scripts == nil {
		m.scripts = make(map[string]Script)
	}

	for _, script := range script {
		m.scripts[hex.EncodeToString(script.Hash())] = script
	}
	return m
}

func (m *MockKugoServer) AddPatterns(patterns ...string) *MockKugoServer {
	m.patterns = patterns
	return m
}

func (m *MockKugoServer) AddMatches(pattern string, matches ...Match) *MockKugoServer {
	if m.matches == nil {
		m.matches = make(map[string][]Match)
	}
	m.matches[pattern] = append(m.matches[pattern], matches...)
	return m
}

func (m *MockKugoServer) AddDatum(hash string, datum string) *MockKugoServer {
	if m.datums == nil {
		m.datums = make(map[string]DatumResponse)
	}
	m.datums[hash] = DatumResponse{Datum: datum}
	return m
}

type MetadatumEntry struct {
	Slot      int
	Tx        string
	Metadatum Metadatum
}

func (m *MockKugoServer) AddMetadata(entries ...MetadatumEntry) *MockKugoServer {
	if m.metadata == nil {
		m.metadata = make(map[int]map[string][]Metadatum)
	}
	for _, entry := range entries {
		if _, ok := m.metadata[entry.Slot]; !ok {
			m.metadata[entry.Slot] = make(map[string][]Metadatum)
		}
		m.metadata[entry.Slot][entry.Tx] = append(m.metadata[entry.Slot][entry.Tx], entry.Metadatum)
	}
	return m
}

type ErrorResponse struct {
	Hint string `json:"hint"`
}

func writeError(w http.ResponseWriter, code int, hint string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Hint: hint})
}
func writeSuccess(w http.ResponseWriter, body interface{}) {
	respBody, _ := json.Marshal(body)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBody)
}

func (m *MockKugoServer) HTTP() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1/scripts/") {
				response, ok := m.scripts[strings.TrimPrefix(r.URL.Path, "/v1/scripts/")]
				if !ok {
					// If a script isn't found, Kugo responds with `null` instead of an error.
					writeSuccess(w, json.RawMessage("null"))
				} else {
					writeSuccess(w, &response)
				}
			} else if r.URL.Path == "/v1/patterns" {
				writeSuccess(w, m.patterns)
			} else if strings.HasPrefix(r.URL.Path, "/v1/metadata/") {
				slotStr := strings.TrimPrefix(r.URL.Path, "/v1/metadata/")
				slot, _ := strconv.Atoi(slotStr)
				tx := r.URL.Query().Get("transaction_id")
				var metadata []Metadatum
				if tx == "" {
					for _, m := range m.metadata[slot] {
						metadata = append(metadata, m...)
					}
				} else {
					metadata = append(metadata, m.metadata[slot][tx]...)
				}
				writeSuccess(w, &metadata)
			} else if strings.HasPrefix(r.URL.Path, "/v1/matches") {
				pattern := strings.TrimPrefix(r.URL.Path, "/v1/matches/")
				// TODO: filter by other query parameters
				matches, ok := m.matches[pattern]
				if !ok {
					writeError(w, http.StatusNotFound, "pattern not found")
				} else {
					writeSuccess(w, &matches)
				}
			} else if strings.HasPrefix(r.URL.Path, "/v1/datums/") {
				hash := strings.TrimPrefix(r.URL.Path, "/v1/datums/")
				datum, ok := m.datums[hash]
				if !ok {
					writeError(w, http.StatusNotFound, "datum not found")
				} else {
					writeSuccess(w, &datum)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)
}
