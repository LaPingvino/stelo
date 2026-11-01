//go:build js

package store

import "syscall/js"

func storage() js.Value {
	defer func() { recover() }() // storage can be refused (private mode, blocked site data)
	return js.Global().Get("localStorage")
}

func read() ([]byte, bool) {
	s := storage()
	if s.IsUndefined() || s.IsNull() {
		return nil, false
	}
	v := s.Call("getItem", key)
	if v.IsNull() || v.IsUndefined() {
		return nil, false
	}
	return []byte(v.String()), true
}

func write(data []byte) {
	if s := storage(); !s.IsUndefined() && !s.IsNull() {
		func() {
			defer func() { recover() }()
			s.Call("setItem", key, string(data))
		}()
	}
}
