package spago

import (
	"fmt"
	"syscall/js"
)

// T text Node
type text string

// Render ...
func (t text) html(bind bool) js.Value {
	return document.Call("createTextNode", string(t))
}

func (t text) apply(n *Node) {
	n.children = append(n.children, t)
}

// T is
func T(s ...interface{}) Markup {
	return text(fmt.Sprint(s...))
}

// ClassMap ...
type ClassMap map[string]bool

func (c ClassMap) apply(n *Node) {
	for k, v := range c {
		n.classmap[k] = v
	}
}

type binded struct {
	name string
	fn   js.Func
}

// Node ...
type Node struct {
	namespace  string
	tag        string
	attributes []attribute
	classmap   ClassMap
	children   []ComponentOrHTML
	listeners  []listener
}

func (n *Node) apply(nn *Node) {
	nn.children = append(nn.children, n)
}

func appendChild(parent, children js.Value) {
	for _, v := range expandNodes(children) {
		parent.Call("appendChild", v)
	}
}

// Render ...
func (n *Node) html(bind bool) js.Value {
	var jsv js.Value
	if len(n.namespace) > 0 {
		jsv = document.Call("createElementNS", n.namespace, n.tag)
	} else {
		jsv = document.Call("createElement", n.tag)
	}
	for _, a := range n.attributes {
		jsv.Call("setAttribute", a.Key, a.Value)
	}
	clist := jsv.Get("classList")
	for k, v := range n.classmap {
		if v {
			clist.Call("add", k)
		} else {
			clist.Call("remove", k)
		}
	}
	for _, c := range n.children {
		switch v := c.(type) {
		case HTML:
			appendChild(jsv, v.html(bind))
		case Component:
			if um, ok := v.(Unmounter); ok {
				if !v.get().target.IsUndefined() {
					um.Unmount()
				}
			}
			if m, ok := v.(Mounter); ok {
				mounts = append(mounts, m)
			}
			appendChild(jsv, v.Render().html(bind))
		}
	}
	binds := []binded{}
	if bind {
		for _, l := range n.listeners {
			fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				l.Func(args[0])
				return nil
			})
			binds = append(binds, binded{l.Name, fn})
			jsv.Call("addEventListener", l.Name, fn)
		}
	}
	var cb js.Func
	// jsv, binds, cb を closure に渡す
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		childNodes := jsv.Get("childNodes")
		for i := 0; i < childNodes.Length(); i++ {
			child := childNodes.Index(i)
			if v := child.Get("release"); !v.IsUndefined() {
				v.Invoke()
			}
		}
		for _, b := range binds {
			jsv.Call("removeEventListener", b.name, b.fn)
			b.fn.Release()
		}
		return nil
	})
	jsv.Set("release", cb)
	return jsv
}
