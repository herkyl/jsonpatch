package jsonpatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var errBadJSONDoc = fmt.Errorf("Invalid JSON Document")

type JsonPatchOperation struct {
	Operation string      `json:"op"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value,omitempty"`
}

func (j *JsonPatchOperation) Json() string {
	b, _ := json.Marshal(j)
	return string(b)
}

func (j *JsonPatchOperation) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("{")
	b.WriteString(fmt.Sprintf(`"op":"%s"`, j.Operation))
	b.WriteString(fmt.Sprintf(`,"path":"%s"`, j.Path))
	// Consider omitting Value for non-nullable operations.
	if j.Value != nil || j.Operation == "replace" || j.Operation == "add" {
		v, err := json.Marshal(j.Value)
		if err != nil {
			return nil, err
		}
		b.WriteString(`,"value":`)
		b.Write(v)
	}
	b.WriteString("}")
	return b.Bytes(), nil
}

type ByPath []JsonPatchOperation

func (a ByPath) Len() int           { return len(a) }
func (a ByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

func NewPatch(operation, path string, value interface{}) JsonPatchOperation {
	return JsonPatchOperation{Operation: operation, Path: path, Value: value}
}

// CreatePatch creates a patch as specified in http://jsonpatch.com/
//
// 'a' is original, 'b' is the modified document. Both are to be given as json encoded content.
// The function will return an array of JsonPatchOperations
//
// An error will be returned if any of the two documents are invalid.
func CreatePatch(a, b []byte) ([]JsonPatchOperation, error) {
	var aI interface{}
	var bI interface{}

	err := json.Unmarshal(a, &aI)
	if err != nil {
		return nil, errBadJSONDoc
	}
	err = json.Unmarshal(b, &bI)
	if err != nil {
		return nil, errBadJSONDoc
	}

	return diff(aI, bI, "", []JsonPatchOperation{})
}

// From http://tools.ietf.org/html/rfc6901#section-4 :
//
// Evaluation of each reference token begins by decoding any escaped
// character sequence.  This is performed by first transforming any
// occurrence of the sequence '~1' to '/', and then transforming any
// occurrence of the sequence '~0' to '~'.
//   TODO decode support:
//   var rfc6901Decoder = strings.NewReplacer("~1", "/", "~0", "~")

var rfc6901Encoder = strings.NewReplacer("~", "~0", "/", "~1")

func makePath(path string, newPart interface{}) string {
	key := rfc6901Encoder.Replace(fmt.Sprintf("%v", newPart))
	if path == "" {
		return "/" + key
	}
	if strings.HasSuffix(path, "/") {
		return path + key
	}
	return path + "/" + key
}

func diff(a, b interface{}, p string, patch []JsonPatchOperation) ([]JsonPatchOperation, error) {
	// If values are not of the same type simply replace
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		patch = append(patch, NewPatch("replace", p, b))
		return patch, nil
	}

	var err error
	switch at := a.(type) {
	case map[string]interface{}:
		bt := b.(map[string]interface{})
		patch, err = diffObjects(at, bt, p, patch)
		if err != nil {
			return nil, err
		}
	case string, float64, bool:
		if !reflect.DeepEqual(a, b) {
			patch = append(patch, NewPatch("replace", p, b))
		}
	case []interface{}:
		bt, ok := b.([]interface{})
		if !ok {
			// array replaced by non-array
			patch = append(patch, NewPatch("replace", p, b))
		} else if len(at) != len(bt) {
			// arrays are not the same length
			patch = append(patch, diffArrays(at, bt, p)...)
		} else {
			for i := range bt {
				patch, err = diff(at[i], bt[i], makePath(p, i), patch)
				if err != nil {
					return nil, err
				}
			}
		}
	case nil:
		switch b.(type) {
		case nil:
			// Both nil, fine.
		default:
			patch = append(patch, NewPatch("add", p, b))
		}
	default:
		panic(fmt.Sprintf("Unknown type:%T ", a))
	}
	return patch, nil
}

// diff returns the (recursive) difference between a and b as an array of JsonPatchOperations.
func diffObjects(a, b map[string]interface{}, path string, patch []JsonPatchOperation) ([]JsonPatchOperation, error) {
	for key, bv := range b {
		p := makePath(path, key)
		av, ok := a[key]
		// Key doesn't exist in original document, value was added
		if !ok {
			patch = append(patch, NewPatch("add", p, bv))
			continue
		}
		// If types have changed, replace completely
		if reflect.TypeOf(av) != reflect.TypeOf(bv) {
			patch = append(patch, NewPatch("replace", p, bv))
			continue
		}
		// Types are the same, compare values
		var err error
		patch, err = diff(av, bv, p, patch)
		if err != nil {
			return nil, err
		}
	}
	// Now add all deleted values as nil
	for key := range a {
		_, ok := b[key]
		if !ok {
			p := makePath(path, key)
			patch = append(patch, NewPatch("remove", p, nil))
		}
	}
	return patch, nil
}

func diffArrays(av, bv []interface{}, p string) []JsonPatchOperation {
	retval := []JsonPatchOperation{}
	// Find elements that need to be removed
	processArray(av, bv, func(i int, value interface{}) {
		retval = append(retval, NewPatch("remove", makePath(p, i), nil))
	})
	// Find elements that need to be added.
	// NOTE we pass in `bv` then `av` so that processArray can find the missing elements.
	processArray(bv, av, func(i int, value interface{}) {
		retval = append(retval, NewPatch("add", makePath(p, i), value))
	})
	return retval
}

// processArray processes `av` and `bv` calling `applyOp` whenever a value is absent.
// It keeps track of which indexes have already had `applyOp` called for and automatically skips them so you can process duplicate objects correctly.
func processArray(av, bv []interface{}, applyOp func(i int, value interface{})) {
	foundIndexes := make(map[int]bool, len(av))
	reverseFoundIndexes := make(map[int]bool, len(av))
	for i, v := range av {
		for i2, v2 := range bv {
			if _, ok := reverseFoundIndexes[i2]; ok {
				// We already found this index.
				continue
			}
			if reflect.DeepEqual(v, v2) {
				// Mark this index as found since it matches exactly.
				foundIndexes[i] = true
				reverseFoundIndexes[i2] = true
				break
			}
		}
		if _, ok := foundIndexes[i]; !ok {
			applyOp(i, v)
		}
	}
}
