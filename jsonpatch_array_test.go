package jsonpatch

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	arrayBase = `{
  "persons": [{"name":"Ed"},{}]
}`

	arrayUpdated = `{
  "persons": [{"name":"Ed"},{},{}]
}`
)

func TestArrayAddMultipleEmptyObjects(t *testing.T) {
	patch, e := CreatePatch([]byte(arrayBase), []byte(arrayUpdated))
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch), "they should be equal")
	sort.Sort(ByPath(patch))

	change := patch[0]
	assert.Equal(t, "add", change.Operation, "they should be equal")
	assert.Equal(t, "/persons/2", change.Path, "they should be equal")
	assert.Equal(t, map[string]interface{}{}, change.Value, "they should be equal")
}

func TestArrayRemoveMultipleEmptyObjects(t *testing.T) {
	patch, e := CreatePatch([]byte(arrayUpdated), []byte(arrayBase))
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch), "they should be equal")
	sort.Sort(ByPath(patch))

	change := patch[0]
	assert.Equal(t, "remove", change.Operation, "they should be equal")
	assert.Equal(t, "/persons/2", change.Path, "they should be equal")
	assert.Equal(t, nil, change.Value, "they should be equal")
}

var (
	arrayWithSpacesBase = `{
	"persons": [{"name":"Ed"},{},{},{"name":"Sally"},{}]
}`

	arrayWithSpacesUpdated = `{
  "persons": [{"name":"Ed"},{},{"name":"Sally"},{}]
}`
)

// TestArrayRemoveSpaceInbetween tests removing one blank item from a group blanks which is in between non blank items which also end with a blank item. This tests that the correct index is removed
func TestArrayRemoveSpaceInbetween(t *testing.T) {
	t.Skip("This test fails. TODO change compareArray algorithm to match by index instead of by object equality")
	patch, e := CreatePatch([]byte(arrayWithSpacesBase), []byte(arrayWithSpacesUpdated))
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch), "they should be equal")
	sort.Sort(ByPath(patch))

	change := patch[0]
	assert.Equal(t, "remove", change.Operation, "they should be equal")
	assert.Equal(t, "/persons/2", change.Path, "they should be equal")
	assert.Equal(t, nil, change.Value, "they should be equal")
}
