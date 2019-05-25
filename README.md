# go-siftjs
Sift.js for golang

---
This package attempts to emulate the functionality of [sift.js](https://github.com/crcn/sift.js) which provides MongoDB document queries on objects.

This package relies heavily on reflect, turn back now if this makes you feel un-easy.

## Installation
```golang
go get github.com/bhoriuchi/go-siftjs
```

## Example

```golang
import (
  	"encoding/json"
    "github.com/bhoriuchi/go-siftjs"
)

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
```

# Supported Operators

* `$eq` - Equals
* `$ne` - Not equals
* `$lt` - Less than
* `$lte` - Less than or equal to
* `$gt` - Greater than
* `$gte` - Greater than or equal to
* `$in` - In
* `$nin` - Not in
* `$all` - All
* `$and` - And
* `$or` - Or
* `$nor` - Nor
* `$regex` - RegEx (specified in js format `/expression/options`)
* `$size` - Size
* `$not` - Not