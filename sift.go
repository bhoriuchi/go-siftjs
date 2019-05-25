package siftjs

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// in is true if the array contains a value
func in(val interface{}, slice interface{}) bool {
	if !IsArrayLike(slice) {
		return false
	}
	a, err := Arrayify(slice)
	if err != nil {
		return false
	}
	for _, v := range a {
		if reflect.DeepEqual(val, v) {
			return true
		}
	}

	return false
}

// all returns true if the document array contains all
// the values in the query array
func all(val interface{}, querySlice interface{}) bool {
	if !IsArrayLike(querySlice) {
		return false
	}
	if !IsArrayLike(val) {
		return false
	}
	v, err := Arrayify(val)
	if err != nil {
		return false
	}
	q, err := Arrayify(querySlice)
	if err != nil {
		return false
	}

	// loop through the query slice
	for _, qv := range q {
		found := false
		// loop through the value slice
		for _, vv := range v {
			// if the value was found, break this loop
			if reflect.DeepEqual(qv, vv) {
				found = true
				break
			}
		}
		// if not found return false
		if !found {
			return false
		}
	}

	// if we made it to the end all the values were there
	return true
}

// lt returns true if the document value is less than the query value
func lt(left interface{}, right interface{}) bool {
	if IsNumberLike(left) && IsNumberLike(right) {
		leftFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", left), 64)
		rightFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", right), 64)
		return leftFloat < rightFloat
	}
	return false
}

// lte returns true if the document value is less than or equal to the query value
func lte(left interface{}, right interface{}) bool {
	return reflect.DeepEqual(left, right) || lt(left, right)
}

// tt returns true if the document value is greater than the query value
func gt(left interface{}, right interface{}) bool {
	if IsNumberLike(left) && IsNumberLike(right) {
		leftFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", left), 64)
		rightFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", right), 64)
		return leftFloat > rightFloat
	}
	return false
}

// gte returns true if the document value is greater than or equal to the query value
func gte(left interface{}, right interface{}) bool {
	return reflect.DeepEqual(left, right) || gt(left, right)
}

// and returns true if all sub-queries return true
func and(doc interface{}, querys interface{}) bool {
	qs, err := Arrayify(querys)
	if err != nil {
		return false
	}
	for _, q := range qs {
		if !compare(q, doc) {
			return false
		}
	}
	return true
}

// regex returns true if the regex matches the document value
func regex(val interface{}, rx interface{}) bool {
	// both sides must be a string to do an rx comparison
	if !IsString(val) || !IsString(rx) {
		return false
	}
	// extract the javascript style regex
	rxRx := regexp.MustCompile(`^/(.+)/([igms]+)?`)
	rxMatch := rxRx.FindAllStringSubmatch(rx.(string), -1)
	if len(rxMatch) == 0 {
		return false
	}
	options := ""
	if rxMatch[0][2] != "" {
		options = fmt.Sprintf("(?%s)", rxMatch[0][2])
	}

	// build a new regexp from the javascript style regex string
	qRx, err := regexp.Compile(fmt.Sprintf("%s%s", options, rxMatch[0][1]))
	if err != nil {
		return false
	}

	// finally match the value
	return qRx.MatchString(val.(string))
}

// or returns true if at least one sub-query returns true
func or(doc interface{}, querys interface{}) bool {
	qs, err := Arrayify(querys)
	if err != nil {
		return false
	}
	for _, q := range qs {
		if compare(q, doc) {
			return true
		}
	}
	return false
}

// size returns true if the document value is an array of the size specified
func size(doc interface{}, s interface{}) bool {
	if !IsArrayLike(doc) || !IsIntLike(s) {
		return false
	}
	a, err := Arrayify(doc)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(len(a), s)
}

// compare compares the document value to the query
func compare(query interface{}, doc interface{}) bool {
	// if the query is an array, treat it as an OR query
	if IsArrayLike(query) {
		a, err := Arrayify(query)
		if err != nil {
			return false
		}
		for _, q := range a {
			if compare(q, doc) {
				return true
			}
		}
	}

	// if the query is not a map or array/slice we are at the value
	// use deepequal to compare
	if !IsMap(query) {
		return reflect.DeepEqual(query, doc)
	}

	// convert query to a map of interfaces
	q, err := Mapify(query)
	if err != nil {
		return false
	}

	// check that there are keys
	if len(q) == 0 {
		return false
	}

	// loop through each query map string
	for k, v := range q {
		var res bool
		switch k {
		case "$eq":
			res = reflect.DeepEqual(doc, v)
		case "$ne":
			res = !reflect.DeepEqual(doc, v)
		case "$lt":
			res = lt(doc, v)
		case "$lte":
			res = lte(doc, v)
		case "$gt":
			res = gt(doc, v)
		case "$gte":
			res = gte(doc, v)
		case "$in":
			res = in(doc, v)
		case "$nin":
			res = !in(doc, v)
		case "$all":
			res = all(doc, v)
		case "$and":
			res = and(doc, v)
		case "$or":
			res = or(doc, v)
		case "$nor":
			res = !or(doc, v)
		case "$regex":
			res = regex(doc, v)
		case "$size":
			res = size(doc, v)
		case "$not":
			res = !compare(doc, v)
		default:
			if IsMap(doc) {
				m, err := Mapify(doc)
				if err != nil {
					return false
				}
				val, ok := m[k]
				if !ok {
					return false
				}
				res = compare(v, val)
			}
		}
		// if any of the keys were false, return false
		if !res {
			return false
		}
	}

	// if we got here without returning false its true
	return true
}

// Sift sifts an interface returning an array/slice of matches
func Sift(query interface{}, docs interface{}) []interface{} {
	result := make([]interface{}, 0)
	d, err := Arrayify(docs)
	if err != nil {
		return result
	}

	// loop through each document and compare it to the query
	for _, doc := range d {
		if matched := compare(query, doc); matched {
			result = append(result, doc)
		}
	}

	// return the result
	return result
}
