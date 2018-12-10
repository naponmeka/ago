package ago

import (
	"syscall/js"
)

const HISTORY_MODE = "HISTORY"
const HASH_MODE = "HASH"

// Route ...
type Route struct {
	Path      string
	Component Component
}

// Router ...
type Router struct {
	Routes  []Route
	Mode    string
	Root    string
	RootDom js.Value
}

// Navigate ...
func (rt Router) Navigate(path string) {
	for _, route := range rt.Routes {
		if route.Path == path {
			js.Global().Get("window").Get("history").Call("pushState", nil, "", path)
			for rt.RootDom.Call("hasChildNodes").Bool() == true {
				rt.RootDom.Call("removeChild", rt.RootDom.Get("lastChild"))
			}
			route.Component.Render(rt.RootDom)
			return
		}
	}
	js.Global().Get("window").Get("history").Call("pushState", nil, "", path)
	rt.RootDom.Set("innerHTML", "not found")
}

// AddRoute ...
func (rt *Router) AddRoute(r Route) {
	rt.Routes = append(rt.Routes, r)
}
