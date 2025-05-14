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
	"encoding/json"

	"github.com/SundaeSwap-finance/ogmigo/ouroboros/chainsync/num"
	"github.com/SundaeSwap-finance/ogmigo/ouroboros/shared"
)

type Match struct {
	TransactionIndex int             `json:"transaction_index,omitempty"`
	TransactionID    string          `json:"transaction_id,omitempty"`
	OutputIndex      int             `json:"output_index,omitempty"`
	Address          string          `json:"address,omitempty"`
	DatumHash        string          `json:"datum_hash,omitempty"`
	DatumType        string          `json:"datum_type,omitempty"`
	Value            Value           `json:"value,omitempty"`
	CreatedAt        Point           `json:"created_at,omitempty"`
	SpentAt          Point           `json:"spent_at,omitempty"`
	ScriptHash       string          `json:"script_hash,omitempty"`
}

type Value shared.Value

func (c *Value) UnmarshalJSON(data []byte) error {
	var v shared.Value
	err := json.Unmarshal(data, &v)
	if err == nil {
		*c = Value(v)
		return nil
	}
	type ValueV5 struct {
		Coins  num.Int                    `json:"coins,omitempty"`
		Assets map[shared.AssetID]num.Int `json:"assets"`
	}
	var r5 ValueV5
	err = json.Unmarshal(data, &r5)
	if err != nil {
		return err
	}
	s := shared.Value{}
	if r5.Coins.BigInt().BitLen() != 0 {
		s.AddAsset(shared.CreateAdaCoin(r5.Coins))
	}
	for asset, coins := range r5.Assets {
		s.AddAsset(shared.Coin{AssetId: asset, Amount: coins})
	}
	*c = Value(s)
	return nil
}

type Point struct {
	SlotNo     int    `json:"slot_no,omitempty"`
	HeaderHash string `json:"header_hash,omitempty"`
}
