/*
Package siftjs performs document queries on interfaces

This package attempts to emulate the functionality of sift.js which provides MongoDB document queries on objects.

This package relies heavily on reflect, turn back now if this makes you feel un-easy.

Usage

	queryString := `
	{
		"foo": {
			"$eq": "bar"
		}
	}
	`
	docsString := `[
		{
			"foo": "bar"
		},
		{
			"foo": "baz"
		},
		{
			"baz": "qux"
		}
	]`

	var query interface{}
	var docs interface{}
	json.Unmarshal([]byte(queryString), &query)
	json.Unmarshal([]byte(docsString), &docs)

	results := siftjs.Sift(query, docs)
	// [{"foo":"bar"}]
*/
package siftjs
