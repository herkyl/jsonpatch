package jsonpatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "b", "c", "d"},
		[]interface{}{1, "a", "c", "d", 2},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 3, len(patch))

	change := patch[0]
	assert.Equal(t, "add", change.Operation)
	assert.Equal(t, "/0", change.Path)
	assert.Equal(t, 1, change.Value)

	change = patch[1]
	assert.Equal(t, "remove", change.Operation)
	assert.Equal(t, "/2", change.Path)
	assert.Equal(t, nil, change.Value)

	change = patch[2]
	assert.Equal(t, "add", change.Operation)
	assert.Equal(t, "/4", change.Path)
	assert.Equal(t, 2, change.Value)
}

func TestRemovingCenterElement(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "b", "c"},
		[]interface{}{"a", "c"},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch))

	change := patch[0]
	assert.Equal(t, "remove", change.Operation)
	assert.Equal(t, "/1", change.Path)
	assert.Equal(t, nil, change.Value)
}

func TestAddingCenterElement(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "c"},
		[]interface{}{"a", "b", "c"},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch))

	change := patch[0]
	assert.Equal(t, "add", change.Operation)
	assert.Equal(t, "/1", change.Path)
	assert.Equal(t, "b", change.Value)
}

func TestAddingElementToEnd(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "b"},
		[]interface{}{"a", "b", "c"},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch))

	change := patch[0]
	assert.Equal(t, "add", change.Operation)
	assert.Equal(t, "/2", change.Path)
	assert.Equal(t, "c", change.Value)
}

func TestAddingElementToStart(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "b"},
		[]interface{}{"0", "a", "b"},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch))

	change := patch[0]
	assert.Equal(t, "add", change.Operation)
	assert.Equal(t, "/0", change.Path)
	assert.Equal(t, "0", change.Value)
}

func TestRemovingLastElement(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "b", "c"},
		[]interface{}{"a", "b"},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch))

	change := patch[0]
	assert.Equal(t, "remove", change.Operation)
	assert.Equal(t, "/2", change.Path)
	assert.Equal(t, nil, change.Value)
}

func TestRemovingFirstElement(t *testing.T) {
	patch, e := diffArrays(
		[]interface{}{"a", "b", "c"},
		[]interface{}{"b", "c"},
		"",
		true,
	)
	assert.NoError(t, e)
	t.Log("Patch:", patch)
	assert.Equal(t, 1, len(patch))

	change := patch[0]
	assert.Equal(t, "remove", change.Operation)
	assert.Equal(t, "/0", change.Path)
	assert.Equal(t, nil, change.Value)
}
