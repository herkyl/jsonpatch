package jsonpatch

import (
	"fmt"
	"reflect"
)

type tmpEl struct {
	val     interface{}
	isFixed bool
}

func diffArrays(a, b []interface{}, p string, forceFullPatch bool) ([]JSONPatchOperation, error) {
	fullReplace := []JSONPatchOperation{NewPatch("replace", p, b)}
	patch := []JSONPatchOperation{}

	tmp := make([]tmpEl, len(a))
	for i, ae := range a {
		newEl := tmpEl{val: ae}
		for j := i; j < len(b); j++ {
			if len(b) <= j { //b is out of bounds
				break
			}
			be := b[j]
			if reflect.DeepEqual(ae, be) {
				newEl.isFixed = true // this element should remain in place
			}
		}
		tmp[i] = newEl
	}
	// Now we have an array of elements in which we know the original, unmoved elements

	fmt.Println("a>>>", a)
	fmt.Println("TMP>>>", tmp)

	aIndex := 0
	bIndex := 0
	addedDelta := 0
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	for aIndex+addedDelta < maxLen {
		tmpIndex := aIndex + addedDelta
		newPath := makePath(p, tmpIndex)
		if aIndex >= len(a) && bIndex >= len(b) {
			break
		}
		if aIndex >= len(a) { // a is out of bounds, all new items in b must be adds
			patch = append(patch, NewPatch("add", newPath, b[tmpIndex]))
			addedDelta++
			continue
		}
		if bIndex >= len(b) { // b is out of bounds, all new items in a must be removed
			patch = append(patch, NewPatch("remove", newPath, nil))
			addedDelta--
			aIndex++
			continue
		}
		// can compare elements, so let's compare them
		te := tmp[aIndex]
		for j := bIndex; j < maxLen; j++ {
			be := b[j]
			fmt.Printf("Comparing i=%d j=%d ae=%v be=%v\n", aIndex, j, te.val, be)
			if reflect.DeepEqual(te.val, be) {
				// element is already in b, move on
				bIndex++
				aIndex++
				break
			} else {
				if te.isFixed {
					fmt.Println("add", newPath, be)
					patch = append(patch, NewPatch("add", newPath, be))
					addedDelta++
					bIndex++
					break
				} else {
					fmt.Println("remove", newPath, be)
					patch = append(patch, NewPatch("remove", newPath, nil))
					addedDelta--
					aIndex++
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
