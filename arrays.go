package jsonpatch

import (
	"fmt"
	"reflect"
)

type arrayEl struct {
	val    interface{}
	isBase bool
}

func diffArrays(a, b []interface{}, p string, forceFullPatch bool) ([]JSONPatchOperation, error) {
	fullReplace := []JSONPatchOperation{NewPatch("replace", p, b)}
	patch := []JSONPatchOperation{}

	tmp := make([]arrayEl, len(a))
	for i, ae := range a {
		newEl := arrayEl{val: ae}
		for j := i; j < len(b); j++ {
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

	bPos := 0
	addedDelta := 0
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	for i := 0; i < maxLen; i++ {
		index := i + addedDelta
		if index >= maxLen {
			break
		}
		newPath := makePath(p, index)
		if len(a) <= i {
			patch = append(patch, NewPatch("add", newPath, b[index]))
			addedDelta++
			continue
		}
		te := tmp[i]
		for j := bPos; j < maxLen; j++ {
			be := b[j]
			fmt.Printf("Comparing i=%d j=%d ae=%v be=%v\n", i, j, te.val, be)
			if reflect.DeepEqual(te.val, be) {
				// element is already in b, move on
				bPos++
				break
			} else {
				if te.isBase {
					fmt.Println("add", newPath, be)
					patch = append(patch, NewPatch("add", newPath, be))
					addedDelta++
					bPos++
				} else {
					fmt.Println("remove", newPath, be)
					patch = append(patch, NewPatch("remove", newPath, nil))
					addedDelta--
					break
				}
			}
		}
	}

	fmt.Println("patch>>>", patch)

	if forceFullPatch {
		return patch, nil
	}
	return getSmallestPatch(fullReplace, patch), nil
}
