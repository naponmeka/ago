package ago

import (
	"fmt"
	"syscall/js"
)

// Element is ...
type Element struct {
	DomType    string
	Props      map[string]interface{}
	DomContent string
	Children   *[]Element
	Dom        js.Value
}

// CreateElementContent is for creating Content element (text node)
func CreateElementContent(content string, createRealDom bool) (elements []Element) {
	e := Element{
		DomType:    "content",
		DomContent: content,
	}
	if createRealDom == true {
		e.Dom = js.Global().Get("document").Call("createTextNode", content)
	}
	elements = append(elements, e)
	return
}

// CreateElement is for creating element
func CreateElement(domType string, props map[string]interface{}, children []Element, createRealDom bool) (e Element) {
	e.DomType = domType
	e.Children = &children
	e.Props = props
	if createRealDom == true {
		e.Dom = js.Global().Get("document").Call("createElement", domType)
		for _, c := range children {
			e.Dom.Call("appendChild", c.Dom)
		}
		for k, v := range props {
			setProp(&e, k, v)
		}

	}
	return e
}

func setProp(e *Element, key string, prop interface{}) {
	e.Props[key] = prop
	if prop, ok := prop.(string); ok {
		e.Dom.Call("setAttribute", key, prop)
		if key == "route" {
			e.Dom.Call("setAttribute",
				"onclick",
				fmt.Sprintf("javascript:%s(\"%s\");", NAV_TO_JS_FUNC, prop),
			)
		}
	}
	if prop, ok := prop.(map[string]string); ok {
		str := ""
		for k, a := range prop {
			str = str + fmt.Sprintf(" %s:%s;", k, a)
		}
		e.Dom.Call("setAttribute", key, str)
	}
}

func removeProp(e *Element, key string) {
	delete(e.Props, key)
	e.Dom.Call("removeAttribute", key)
}

// CreateElementRecursive is ..
func CreateElementRecursive(domType string, domContent string, props map[string]interface{}, children []Element, createRealDom bool) Element {
	childElements := []Element{}
	for _, child := range children {
		childrenOfChild := []Element{}
		if child.Children != nil {
			childrenOfChild = *child.Children
		}
		childElements = append(childElements, CreateElementRecursive(child.DomType, child.DomContent, child.Props, childrenOfChild, createRealDom))
	}
	if len(children) == 0 && string(domContent) != "" {
		childElements = CreateElementContent(string(domContent), createRealDom)
	}
	x := CreateElement(domType, props, childElements, createRealDom)
	return x
}

// Render is ..
func Render(rootElement js.Value, element Element) {
	rootElement.Call("appendChild", element.Dom)
}

// RemoveAllChild ...
func RemoveAllChild(dom js.Value) {
	for dom.Call("hasChildNodes").Bool() == true {
		dom.Call("removeChild", dom.Get("lastChild"))
	}
}
