package jsonpatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	patch, e := CreatePatch([]byte(`{"foo": ["a", "b", "c", "d"]}`), []byte(`{"foo": [1, "a", "c", "d", 2]}`))
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	// assert.Equal(t, 1, len(patch), "they should be equal")
	// sort.Sort(ByPath(patch))

	// change := patch[0]
	// assert.Equal(t, "add", change.Operation, "they should be equal")
	// assert.Equal(t, "/persons/2", change.Path, "they should be equal")
	// assert.Equal(t, map[string]interface{}{}, change.Value, "they should be equal")
}
