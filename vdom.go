package ago

const CREATE = "CREATE"
const REMOVE = "REMOVE"
const REPLACE = "REPLACE"
const UPDATE = "UPDATE"
const SET_PROP = "SET_PROP"
const REMOVE_PROP = "REMOVE PROP"

type patch struct {
	mode            string
	element         *Element
	propsPatches    []propPatch
	childrenPatches []patch
}

type propPatch struct {
	mode  string
	name  string
	value interface{}
}

func changed(newElement, oldElement *Element) bool {
	return newElement.DomType != oldElement.DomType ||
		newElement.DomType == "content" && newElement.DomContent != oldElement.DomContent
}

func diffProps(newElement, oldElement *Element) []propPatch {
	propPatches := []propPatch{}
	for k, val := range oldElement.Props {
		if _, exits := newElement.Props[k]; !exits {
			propPatches = append(propPatches, propPatch{mode: REMOVE_PROP, name: k, value: val})
		}
	}
	for k, val := range newElement.Props {
		oldVal, oldExits := oldElement.Props[k]
		if !oldExits || oldVal != val {
			propPatches = append(propPatches, propPatch{mode: SET_PROP, name: k, value: val})

		}
	}
	return propPatches
}

func diff(newElement, oldElement *Element) patch {
	if oldElement == nil {
		return patch{mode: CREATE, element: newElement}
	} else if newElement == nil {
		return patch{mode: REMOVE}
	} else if changed(newElement, oldElement) {
		return patch{mode: REPLACE, element: newElement}
	} else if newElement.DomType != "content" {
		return patch{
			mode:            UPDATE,
			propsPatches:    diffProps(newElement, oldElement),
			childrenPatches: diffChildren(newElement, oldElement),
		}
	}
	return patch{}
}

func diffChildren(newElement, oldElement *Element) (patches []patch) {
	maximumPatches := len(*newElement.Children)
	if maximumPatches < len(*oldElement.Children) {
		maximumPatches = len(*oldElement.Children)
	}
	for i := 0; i < maximumPatches; i++ {
		var newElementChild *Element
		if i < len(*newElement.Children) {
			newElementChild = &(*newElement.Children)[i]
		} else {
			newElementChild = nil
		}
		var oldElementChild *Element
		if i < len(*oldElement.Children) {
			oldElementChild = &(*oldElement.Children)[i]
		} else {
			oldElementChild = nil
		}
		patches = append(patches, diff(newElementChild, oldElementChild))
	}
	return
}

func patchDiff(parent *Element, p patch, index int) int {
	currentEle := &Element{}
	if index < len(*parent.Children) {
		currentEle = &(*parent.Children)[index]
	}

	switch p.mode {
	case CREATE:
		newEle := Element{}
		if p.element.DomType == "content" {
			newEle = CreateElementContent(p.element.DomContent, true)[0]
		} else {
			newEle = CreateElementRecursive(p.element.DomType, p.element.DomContent, p.element.Props, *p.element.Children, true)
		}
		parent.Dom.Call("appendChild", newEle.Dom) // !!!!
		*parent.Children = append(*parent.Children, newEle)
	case REMOVE:
		parent.Dom.Call("removeChild", currentEle.Dom)
		return index
	case REPLACE:
		newEle := Element{}
		if p.element.DomType == "content" {
			newEle = CreateElementContent(p.element.DomContent, true)[0]
		} else {
			newEle = CreateElementRecursive(p.element.DomType, p.element.DomContent, p.element.Props, *p.element.Children, true)
		}
		parent.Dom.Call("replaceChild", newEle.Dom, currentEle.Dom)
		(*parent.Children)[index] = newEle
	case UPDATE:
		patchProp(currentEle, p.propsPatches)
		newChildren := []Element{}
		indexToBeDeleted := make(map[int]bool)
		for i := 0; i < len(p.childrenPatches); i++ {
			childTobeRemoved := patchDiff(currentEle, p.childrenPatches[i], i)
			if childTobeRemoved != -1 && currentEle.Children != nil {
				indexToBeDeleted[childTobeRemoved] = true
			}
		}
		for j := 0; j < len(*currentEle.Children); j++ {
			if ok, _ := indexToBeDeleted[j]; !ok {
				newChildren = append(newChildren, (*currentEle.Children)[j])
			}
		}
		if len(indexToBeDeleted) != 0 {
			currentEle.Children = &newChildren
		}

	}
	return -1
}

func patchProp(element *Element, propPatches []propPatch) {
	for _, pp := range propPatches {
		if pp.mode == SET_PROP {
			setProp(element, pp.name, pp.value)
		} else if pp.mode == REMOVE_PROP {
			removeProp(element, pp.name)
		}
	}
}
