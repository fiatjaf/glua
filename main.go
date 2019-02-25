package main

import (
	lua "github.com/J-J-J/goluajit"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	run := func(code string) {
		L := lua.NewState()
		defer L.Close()
		if err := L.DoString(code); err != nil {
			panic(err)
		}
	}
	withGlobals := func(globals map[string]interface{}, code string) {
		L := lua.NewState()
		defer L.Close()

		for name, value := range globals {
			L.SetGlobal(name, lvalueFromInterface(L, value))
		}

		if err := L.DoString(code); err != nil {
			panic(err)
		}
	}
	withModules := func(modules map[string]string, globals map[string]interface{}, code string) {
		L := lua.NewState()
		defer L.Close()

		preload := L.GetField(L.GetField(L.Get(lua.EnvironIndex), "package"), "preload")
		for moduleName, code := range modules {
			mod, err := L.LoadString(code)
			if err != nil {
				panic(err)
			}
			L.SetField(preload, moduleName, mod)
		}

		for name, value := range globals {
			L.SetGlobal(name, lvalueFromInterface(L, value))
		}

		if err := L.DoString(code); err != nil {
			panic(err)
		}
	}

	if js.Module != js.Undefined {
		js.Module.Get("exports").Set("run", run)
		js.Module.Get("exports").Set("runWithGlobals", withGlobals)
		js.Module.Get("exports").Set("runWithModules", withModules)
	} else {
		js.Global.Set("glua", map[string]interface{}{
			"run":            run,
			"runWithGlobals": withGlobals,
			"runWithModules": withModules,
		})
	}
}

func lvalueFromInterface(L *lua.LState, value interface{}) lua.LValue {
	switch val := value.(type) {
	case string:
		return lua.LString(val)
	case float64:
		return lua.LNumber(val)
	case bool:
		return lua.LBool(val)
	case map[string]interface{}:
		table := L.NewTable()
		for k, iv := range val {
			table.RawSetString(k, lvalueFromInterface(L, iv))
		}
		return table
	case []interface{}:
		table := L.NewTable()
		for i, iv := range val {
			table.RawSetInt(i+1, lvalueFromInterface(L, iv))
		}
		return table
	case func(...interface{}) *js.Object:
		fn := val
		return L.NewFunction(func(L *lua.LState) int {
			var args []interface{}

			for a := 1; ; a++ {
				arg := L.Get(a)
				if arg == lua.LNil {
					break
				}
				args = append(args, lvalueToInterface(arg))
			}

			jsreturn := fn(args...)

			if jsreturn == js.Undefined {
				return 0
			}

			L.Push(lvalueFromInterface(L, jsreturn.Interface()))
			return 1
		})
	default:
		return lua.LNil
	}
}

func lvalueToInterface(lvalue lua.LValue) interface{} {
	switch value := lvalue.(type) {
	case *lua.LTable:
		size := value.Len()

		// it will be either an object or an array
		if size == 0 {
			return make(map[string]interface{}, 0)
		}

		// otherwise we'll have to figure out
		object := make(map[string]interface{}, size)
		array := make([]interface{}, size)

		isArray := true
		value.ForEach(func(k lua.LValue, lv lua.LValue) {
			v := lvalueToInterface(lv)

			if isArray {
				ln, ok := k.(lua.LNumber)
				if !ok || float64(int(ln)) != float64(ln) {
					// has a non-int key, so not an array
					isArray = false
				} else if ln == 0 || int(ln) > size /* because lua arrays are 1-based */ {
					// int out of the allowed range, so not an array
					isArray = false
				} else {
					// keep storing everything in the array
					array[int(ln)-1 /* because lua arrays are 1-based */] = v
				}
			}

			// if in the last key we discover this isn't an array we already have all values
			object[lua.LVAsString(k)] = lvalueToInterface(lv)
		})

		if isArray {
			return array
		} else {
			return object
		}
	case lua.LNumber:
		return float64(value)
	case lua.LString:
		return string(value)
	default:
		switch lvalue {
		case lua.LTrue:
			return true
		case lua.LFalse:
			return false
		case lua.LNil:
			return nil
		}
	}
	return nil
}
