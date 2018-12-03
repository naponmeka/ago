package ago

import (
	"fmt"
	"syscall/js"
)

// Component ...
type Component struct {
	gox   string
	state interface{}
	VDom  Element
	root  Element
}

// CreateComponent ...
func CreateComponent(gox string, state interface{}) Component {
	return Component{
		gox:   gox,
		state: state,
		VDom:  Transform(gox, state),
	}
}

// ChangeState ...
func (c *Component) ChangeState(value interface{}) {
	newElement := Transform(c.gox, value)
	patch := diff(&newElement, &c.VDom)
	fmt.Println(patch)
	patchDiff(&c.root, patch, 0)
	fmt.Println("DONE change state")
}

func (c *Component) Render(parentID string) {
	rootDom := js.Global().Get("document").Call("getElementById", parentID)
	rootDom.Call("appendChild", c.VDom.Dom)
	rootElement := Element{
		Dom:      rootDom,
		Children: &[]Element{c.VDom},
	}
	c.root = rootElement
}
