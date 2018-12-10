package ago

import (
	"syscall/js"
)

const HISTORY_MODE = "HISTORY"
const HASH_MODE = "HASH"
const NAV_TO_JS_FUNC = "agoRouterNavigateTo"

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

// JsNavigate ...
func (rt *Router) JsNavigate(i []js.Value) {
	path := i[0].String()
	rt.Navigate(path)
}

// JsGoBack ...
func (rt *Router) JsGoBack(i []js.Value) {
	path := js.Global().Get("window").Get("location").Get("pathname").String()
	for _, route := range rt.Routes {
		if route.Path == path {
			RemoveAllChild(rt.RootDom)
			route.Component.Render(rt.RootDom)
			return
		}
	}
	rt.RootDom.Set("innerHTML", "not found")
}

// Register ...
func (rt *Router) Register() {
	js.Global().Set(NAV_TO_JS_FUNC, js.NewCallback(rt.JsNavigate))
	js.Global().Get("window").Set("onpopstate", js.NewCallback(rt.JsGoBack))
}

// Navigate ...
func (rt *Router) Navigate(path string) {
	for _, route := range rt.Routes {
		if route.Path == path {
			js.Global().Get("window").Get("history").Call("pushState", nil, "", path)
			RemoveAllChild(rt.RootDom)
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
