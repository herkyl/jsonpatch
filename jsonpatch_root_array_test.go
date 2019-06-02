package jsonpatch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPatchCreate(t *testing.T) {
	cases := map[string]struct {
		a    string
		b    string
		diff string
	}{
		// "object": {
		// 	`{"asdf":"qwerty"}`,
		// 	`{"asdf":"zzz"}`,
		// 	`[{"op":"replace","path":"/asdf","value":"zzz"}]`,
		// },
		// "object with array": {
		// 	`{"items":[{"asdf":"qwerty"}]}`,
		// 	`{"items":[{"asdf":"bla"},{"asdf":"zzz"}]}`,
		// 	`[{"op":"remove","path":"/items/0"},{"op":"add","path":"/items/0","value":{"asdf":"bla"}},{"op":"add","path":"/items/1","value":{"asdf":"zzz"}}]`,
		// },
		"array": {
			`[{"asdf":"qwerty"}]`,
			`[{"asdf":"bla"},{"asdf":"zzz"}]`,
			`[{"op":"replace","path":"/0/asdf","value":"bla"},{"op":"add","path":"/1","value":{"asdf":"zzz"}}]`,
		},
		// "from empty array": {
		// 	`[]`,
		// 	`[{"asdf":"bla"},{"asdf":"zzz"}]`,
		// 	`[{"op":"add","path":"/0","value":{"asdf":"bla"}},{"op":"add","path":"/1","value":{"asdf":"zzz"}}]`,
		// },
		// "to empty array": {
		// 	`[{"asdf":"bla"},{"asdf":"zzz"}]`,
		// 	`[]`,
		// 	`[{"op":"remove","path":"/0"},{"op":"remove","path":"/1"}]`,
		// },
		// "from object to array": {
		// 	`{"foo":"bar"}`,
		// 	`[{"foo":"bar"}]`,
		// 	`[{"op":"replace","path":"","value":[{"foo":"bar"}]}]`,
		// },
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Logf(`Running test: "%s"`, name)
			patch, err := CreatePatch([]byte(tc.a), []byte(tc.b))
			assert.NoError(t, err)

			patchBytes, err := json.Marshal(patch)
			assert.NoError(t, err)

			assert.Equal(t, tc.diff, string(patchBytes))
		})
	}
}
