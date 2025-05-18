package kugo

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tj/assert"
)

func Test_DecodeMatch(t *testing.T) {
	bytes, err := os.ReadFile("testdata/simple.json")
	assert.Nil(t, err)

	var match Match
	err = json.Unmarshal(bytes, &match)
	assert.Nil(t, err)

	assert.Equal(
		t,
		match.TransactionIndex,
		1,
	)
	assert.Equal(
		t,
		match.TransactionID,
		"2222222222222222222222222222222222222222222222222222222222222222",
	)
	assert.Equal(
		t,
		match.OutputIndex,
		3,
	)
}

func Test_DecodeMatchScriptHash(t *testing.T) {
	bytes, err := os.ReadFile("testdata/script_hash.json")
	assert.Nil(t, err)

	var match Match
	err = json.Unmarshal(bytes, &match)
	assert.Nil(t, err)

	assert.Equal(
		t,
		match.ScriptHash,
		"abc",
	)
}
