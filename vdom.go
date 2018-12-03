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
	maximumPatches := len(newElement.Children)
	if maximumPatches < len(oldElement.Children) {
		maximumPatches = len(oldElement.Children)
	}
	for i := 0; i < maximumPatches; i++ {
		var newElementChild *Element
		if i < len(newElement.Children) {
			newElementChild = &newElement.Children[i]
		}
		var oldeElementChild *Element
		if i < len(oldElement.Children) {
			oldeElementChild = &oldElement.Children[i]
		}
		patches = append(patches, diff(newElementChild, oldeElementChild))
	}
	return
}
