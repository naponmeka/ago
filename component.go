package ago

import "fmt"

// Component ...
type Component struct {
	gox   string
	state interface{}
	VDom  Element
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
	fmt.Println(diff(&newElement, &c.VDom))
	fmt.Println("DONE change state")
}
