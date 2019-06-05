package jsonpatch

import (
	"fmt"
	"reflect"
)

type arrayEl struct {
	val    interface{}
	isBase bool
}

func diffArrays(a, b []interface{}, p string) ([]JSONPatchOperation, error) {
	fullReplace := []JSONPatchOperation{NewPatch("replace", p, b)}
	patch := []JSONPatchOperation{}

	tmp := make([]arrayEl, len(a))

	for i, ae := range a {
		newEl := arrayEl{val: ae}
		for j := i; j < len(a); j++ {
			if len(b) <= j { //b is out of bounds
				break
			}
			be := b[j]
			if reflect.DeepEqual(ae, be) {
				newEl.isBase = true
			}
		}
		tmp[i] = newEl
	}
	// Now we have an array of elements in which we know the original, unmoved elements

	fmt.Println("a>>>", a)
	fmt.Println("TMP>>>", tmp)

	for i := 0; i < len(a); i++ {
		te := tmp[i]
		newPath := makePath(p, i)
		for j := i; j < len(a); j++ { //FIXME: what if b is longer than a?
			if len(b) <= j { //b is out of bounds
				break
			}
			be := b[j]
			fmt.Printf("Comparing i=%d j=%d ae=%v be=%v\n", i, j, te.val, be)
			if reflect.DeepEqual(te.val, be) {
				// element is already in b, move on
				break
			} else {
				if te.isBase {
					patch = append(patch, NewPatch("add", newPath, be))
				} else {
					patch = append(patch, NewPatch("remove", newPath, nil))
				}
			}
		}
	}

	fmt.Println("patch>>>", patch)

	return getSmallestPatch(fullReplace, patch), nil
}
