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
	props           map[string]interface{}
	childrenPatches []patch
}

func changed(newElement, oldElement *Element) bool {
	return newElement.DomType != oldElement.DomType ||
		newElement.DomType == "content" && newElement.DomContent != oldElement.DomContent
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
			mode: UPDATE,
			// props:    diffProps(newElement, oldElement),
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
			newEle = CreateElementContent(p.element.DomContent)[0]
		} else {
			newEle = CreateElement(p.element.DomType, p.element.Props, *p.element.Children)
		}
		parent.Dom.Call("appendChild", newEle.Dom)
		*parent.Children = append(*parent.Children, newEle)
	case REMOVE:
		parent.Dom.Call("removeChild", currentEle.Dom)
		return index
	case REPLACE:
		newEle := Element{}
		if p.element.DomType == "content" {
			newEle = CreateElementContent(p.element.DomContent)[0]
		} else {
			newEle = CreateElement(p.element.DomType, p.element.Props, *p.element.Children)
		}
		parent.Dom.Call("replaceChild", newEle.Dom, currentEle.Dom)
		(*parent.Children)[index] = newEle
	case UPDATE:
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
