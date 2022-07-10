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
	"bytes"
	"log"
)

type KeyValue struct {
	Key   string
	Value string
}

func KV(key, value string) KeyValue {
	return KeyValue{
		Key:   key,
		Value: value,
	}
}

type Logger interface {
	Debug(message string, kvs ...KeyValue)
	Info(message string, kvs ...KeyValue)
	With(kvs ...KeyValue) Logger
}

// DefaultLogger logs via the log package
var DefaultLogger = defaultLogger{}

type defaultLogger struct {
	kvs []KeyValue
}

func (d defaultLogger) print(message string, kvs ...KeyValue) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(message)
	if len(kvs) > 0 {
		buf.WriteString(":")
	}
	for _, kv := range kvs {
		buf.WriteString(" ")
		buf.WriteString(kv.Key)
		buf.WriteString("=")
		buf.WriteString(kv.Value)
	}
	log.Println(buf)
}

func (d defaultLogger) Debug(message string, kvs ...KeyValue) {
	d.print(message, kvs...)
}

func (d defaultLogger) Info(message string, kvs ...KeyValue) {
	d.print(message, kvs...)
}

func (d defaultLogger) With(kvs ...KeyValue) Logger {
	return defaultLogger{
		kvs: append(d.kvs, kvs...),
	}
}

// NopLogger logs nothing
var NopLogger = nopLogger{}

type nopLogger struct {
}

func (n nopLogger) Debug(string, ...KeyValue) {}
func (n nopLogger) Info(string, ...KeyValue)  {}
func (n nopLogger) With(...KeyValue) Logger   { return n }
