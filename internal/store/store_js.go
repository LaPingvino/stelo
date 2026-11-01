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

const levelPrefix = "stelo-level-"

func readLevel(name string) (string, bool) {
	s := storage()
	if s.IsUndefined() || s.IsNull() {
		return "", false
	}
	v := s.Call("getItem", levelPrefix+name)
	if v.IsNull() || v.IsUndefined() {
		return "", false
	}
	return v.String(), true
}

func writeLevel(name, text string) {
	if s := storage(); !s.IsUndefined() && !s.IsNull() {
		func() {
			defer func() { recover() }()
			s.Call("setItem", levelPrefix+name, text)
		}()
	}
}

// Download offers the text to the browser as a file to save (<name>.txt).
func Download(name, text string) bool {
	defer func() { recover() }()
	doc := js.Global().Get("document")
	blob := js.Global().Get("Blob").New([]any{text}, map[string]any{"type": "text/plain"})
	url := js.Global().Get("URL").Call("createObjectURL", blob)
	a := doc.Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", name+".txt")
	doc.Get("body").Call("appendChild", a)
	a.Call("click")
	a.Call("remove")
	js.Global().Get("URL").Call("revokeObjectURL", url)
	return true
}

// OpenFile asks the browser for a file (a picker, from a tap or click just
// made) and hands its name and text to got when read. It reports whether a
// picker could be shown.
func OpenFile(got func(name, text string)) bool {
	defer func() { recover() }()
	doc := js.Global().Get("document")
	in := doc.Call("createElement", "input")
	in.Set("type", "file")
	in.Set("accept", ".txt,text/plain")
	var onChange js.Func
	onChange = js.FuncOf(func(js.Value, []js.Value) any {
		onChange.Release()
		files := in.Get("files")
		if files.Length() == 0 {
			return nil
		}
		f := files.Index(0)
		var then js.Func
		then = js.FuncOf(func(_ js.Value, args []js.Value) any {
			then.Release()
			got(f.Get("name").String(), args[0].String())
			return nil
		})
		f.Call("text").Call("then", then)
		return nil
	})
	in.Call("addEventListener", "change", onChange)
	in.Call("click")
	return true
}
