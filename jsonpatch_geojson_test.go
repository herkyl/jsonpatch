package jsonpatch

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var point = `{"type":"Point", "coordinates":[0.0, 1.0], "weight":"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book."}`
var lineString = `{"type":"LineString", "coordinates":[[0.0, 1.0], [2.0, 3.0]], "weight":"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book."}`

func TestPointLineStringReplace(t *testing.T) {
	patch, e := CreatePatch([]byte(point), []byte(lineString))
	assert.NoError(t, e)
	t.Log(patch)
	assert.Equal(t, 3, len(patch), "they should be equal")
	sort.Sort(ByPath(patch))
	change := patch[0]
	assert.Equal(t, change.Operation, "replace", "they should be equal")
	assert.Equal(t, change.Path, "/coordinates/0", "they should be equal")
	assert.Equal(t, change.Value, []interface{}{0.0, 1.0}, "they should be equal")
	change = patch[1]
	assert.Equal(t, change.Operation, "replace", "they should be equal")
	assert.Equal(t, change.Path, "/coordinates", "they should be equal")
	assert.Equal(t, change.Value, []interface{}{[]interface{}{0.0, 1.0}, []interface{}{2.0, 3.0}}, "they should be equal")
}

func TestLineStringPointReplace(t *testing.T) {
	patch, e := CreatePatch([]byte(lineString), []byte(point))
	assert.NoError(t, e)
	t.Log(patch)
	assert.Equal(t, len(patch), 3, "they should be equal")
	sort.Sort(ByPath(patch))
	change := patch[0]
	assert.Equal(t, change.Operation, "replace", "they should be equal")
	assert.Equal(t, change.Path, "/type", "they should be equal")
	assert.Equal(t, change.Value, "Point", "they should be equal")
	change = patch[1]
	assert.Equal(t, change.Operation, "replace", "they should be equal")
	assert.Equal(t, change.Path, "/coordinates", "they should be equal")
	assert.Equal(t, change.Value, []interface{}{0.0, 1.0}, "they should be equal")
}
