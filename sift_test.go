package siftjs

import (
	"testing"
)

var docsString = `
[
	{
		"foo": "bar"
	},
	{
		"foo": "baz"
	},
	{
		"baz": "qux"
	}
]
`

func TestSiftEQ(t *testing.T) {
	docs := FromJSON(docsString)
	query1 := FromJSON(`{
		"foo": "bar"
	}`)
	query2 := FromJSON(`{
		"foo": {
			"$eq": "bar"
		}
	}`)

	result1 := Sift(query1, docs)
	result2 := Sift(query2, docs)

	if len(result1) != 1 {
		t.Errorf("basic equality test failed. got %d, expected: 1", len(result1))
	}
	if len(result2) != 1 {
		t.Errorf("$eq test failed. got %d, expected: 1", len(result1))
	}
}
