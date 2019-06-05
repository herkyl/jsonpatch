package jsonpatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var errBadJSONDoc = fmt.Errorf("Invalid JSON Document")

type JSONPatchOperation struct {
	Operation string      `json:"op"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value,omitempty"`
}

func (j *JSONPatchOperation) JSON() string {
	b, _ := json.Marshal(j)
	return string(b)
}

func (j *JSONPatchOperation) MarshalJSON() ([]byte, error) {
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

type ByPath []JSONPatchOperation

func (a ByPath) Len() int           { return len(a) }
func (a ByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

func NewPatch(operation, path string, value interface{}) JSONPatchOperation {
	return JSONPatchOperation{Operation: operation, Path: path, Value: value}
}

// CreatePatch creates a patch as specified in http://jsonpatch.com/
//
// 'a' is original, 'b' is the modified document. Both are to be given as json encoded content.
// The function will return an array of JSONPatchOperations
//
// An error will be returned if any of the two documents are invalid.
func CreatePatch(a, b []byte) ([]JSONPatchOperation, error) {
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

	return diff(aI, bI, "", []JSONPatchOperation{})
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

func diff(a, b interface{}, p string, patch []JSONPatchOperation) ([]JSONPatchOperation, error) {
	// If values are not of the same type simply replace
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		patch = append(patch, NewPatch("replace", p, b))
		return patch, nil
	}

	var err error
	var patch2 []JSONPatchOperation
	switch at := a.(type) {
	case map[string]interface{}:
		bt := b.(map[string]interface{})
		patch2, err = diffObjects(at, bt, p)
		if err != nil {
			return nil, err
		}
		patch = append(patch, patch2...)
	case string, float64, bool:
		if !reflect.DeepEqual(a, b) {
			patch = append(patch, NewPatch("replace", p, b))
		}
	case []interface{}:
		bt, ok := b.([]interface{})
		if !ok {
			// array replaced by non-array
			patch = append(patch, NewPatch("replace", p, b))
		} else {
			// arrays are not the same length
			patch2, err = diffArrays(at, bt, p)
			if err != nil {
				return nil, err
			}
			patch = append(patch, patch2...)
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
func diffObjects(a, b map[string]interface{}, path string) ([]JSONPatchOperation, error) {
	fullReplace := []JSONPatchOperation{NewPatch("replace", path, b)}
	patch := []JSONPatchOperation{}
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
	return getSmallestPatch(fullReplace, patch), nil
}

func getSmallestPatch(patches ...[]JSONPatchOperation) []JSONPatchOperation {
	smallestPatch := patches[0]
	b, _ := json.Marshal(patches[0])
	smallestSize := len(b)
	for i := 1; i < len(patches); i++ {
		p := patches[i]
		b, _ := json.Marshal(p)
		size := len(b)
		if size < smallestSize {
			smallestPatch = p
			smallestSize = size
		}
	}
	return smallestPatch
}
