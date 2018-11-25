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
	Children   []Element
	Dom        js.Value
}

// CreateElementContent is for creating Content element (text node)
func CreateElementContent(content string) (elements []Element) {
	e := Element{
		DomType:    "content",
		DomContent: content,
		Dom:        js.Global().Get("document").Call("createTextNode", content),
	}
	elements = append(elements, e)
	return
}

// CreateElement is for creating element
func CreateElement(domType string, props map[string]interface{}, children []Element) (e Element) {
	e.DomType = domType
	e.Dom = js.Global().Get("document").Call("createElement", domType)

	for _, c := range children {
		e.Dom.Call("appendChild", c.Dom)
	}
	for k, v := range props {
		if v, ok := v.(string); ok {
			e.Dom.Call("setAttribute", k, v)
		}
		if v, ok := v.(map[string]string); ok {
			str := ""
			for k, a := range v {
				str = str + fmt.Sprintf(" %s:%s;", k, a)
			}
			e.Dom.Call("setAttribute", k, str)
		}
	}
	return e
}

// Render is ..
func Render(element Element, rootElement js.Value) {
	rootElement.Call("appendChild", element.Dom)
}
