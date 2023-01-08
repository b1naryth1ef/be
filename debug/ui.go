package debug

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/inkyblackness/imgui-go/v4"
)

func isInt(t reflect.Kind) bool {
	return (t == reflect.Int || t == reflect.Int8 || t == reflect.Int16 || t == reflect.Int32 || t == reflect.Int64)
}

func isUnsignedInt(t reflect.Kind) bool {
	return (t == reflect.Uint || t == reflect.Uint8 || t == reflect.Uint16 || t == reflect.Uint32 || t == reflect.Uint64)
}

func RenderDebugValue(label string, rawptr interface{}) error {
	ptr := reflect.ValueOf(rawptr)
	if ptr.Kind() != reflect.Pointer {
		return fmt.Errorf("Cannot render debug value for non-pointer: %v (%v)", ptr, ptr.Kind())
	}

	return renderDebugValue(label, ptr, false)
}

func renderDebugValue(label string, ptr reflect.Value, editable bool) error {
	if ptr.IsNil() {
		imgui.LabelText(label, "nil")
		return nil
	}

	valueType := ptr.Type().Elem()
	if valueType.Kind() == reflect.Struct {
		if imgui.CollapsingHeaderV(label, imgui.TreeNodeFlagsFramed) {
			return renderDebugStructValue(ptr)
		}
	} else if valueType.Kind() == reflect.Bool {
		return renderDebugBool(label, ptr, editable)
	} else if valueType.Kind() == reflect.String {
		return renderDebugString(label, ptr, editable)
	} else if valueType.Kind() == reflect.Float32 || valueType.Kind() == reflect.Float64 {
		return renderDebugFloat(label, ptr, editable)
	} else if isInt(valueType.Kind()) {
		return renderDebugInt(label, ptr, editable)
	} else if isUnsignedInt(valueType.Kind()) {
		return renderDebugUnsignedInt(label, ptr, editable)
	} else if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
		return renderDebugArray(label, ptr, editable)
	} else {
		// log.Printf("value type: %v", valueType.Kind())
	}

	return nil
}

func renderDebugArray(label string, ptr reflect.Value, editable bool) error {
	rawValue, err := json.Marshal(ptr.Elem().Interface())
	if err != nil {
		return err
	}

	value := string(rawValue)

	flags := imgui.InputTextFlagsNone
	if !editable {
		flags |= imgui.InputTextFlagsReadOnly
	}

	if imgui.InputTextV(label, &value, flags, nil) {
		json.Unmarshal([]byte(value), ptr.Interface())
	}

	return nil
}

func renderDebugUnsignedInt(label string, ptr reflect.Value, editable bool) error {
	value := fmt.Sprintf("%v", ptr.Elem().Uint())
	flags := imgui.InputTextFlagsCharsDecimal
	if !editable {
		flags |= imgui.InputTextFlagsReadOnly
	}

	if imgui.InputTextV(label, &value, flags, nil) {
		v, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			ptr.Elem().SetUint(v)
		}
	}

	return nil
}

func renderDebugInt(label string, ptr reflect.Value, editable bool) error {
	value := fmt.Sprintf("%v", ptr.Elem().Int())
	flags := imgui.InputTextFlagsCharsDecimal
	if !editable {
		flags |= imgui.InputTextFlagsReadOnly
	}

	if imgui.InputTextV(label, &value, flags, nil) {
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			ptr.Elem().SetInt(v)
		}
	}

	return nil
}

func renderDebugFloat(label string, ptr reflect.Value, editable bool) error {
	value := fmt.Sprintf("%f", ptr.Elem().Float())
	flags := imgui.InputTextFlagsCharsDecimal
	if !editable {
		flags |= imgui.InputTextFlagsReadOnly
	}

	if imgui.InputTextV(label, &value, flags, nil) {
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			ptr.Elem().SetFloat(v)
		}
	}

	return nil
}

func renderDebugString(label string, ptr reflect.Value, editable bool) error {
	value := ptr.Elem().String()
	flags := imgui.InputTextFlagsNone
	if !editable {
		flags |= imgui.InputTextFlagsReadOnly
	}

	if imgui.InputTextV(label, &value, flags, nil) {
		ptr.Elem().SetString(value)
	}

	return nil
}

func renderDebugBool(label string, ptr reflect.Value, editable bool) error {
	value := ptr.Elem().Bool()

	if imgui.Checkbox(label, &value) && editable {
		ptr.Elem().SetBool(value)
	}

	return nil
}

func getUnexportedFieldPtr(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))
}

func renderDebugStructValue(ptr reflect.Value) error {
	structValue := ptr.Elem()
	structType := structValue.Type()

	for fieldIdx := 0; fieldIdx < structType.NumField(); fieldIdx++ {
		fieldType := structType.Field(fieldIdx)
		fieldValue := structValue.Field(fieldIdx)
		fieldTag := fieldType.Tag.Get("debug")
		editable := false

		if fieldTag == "-" || (fieldTag == "" && !fieldType.IsExported()) {
			continue
		} else if fieldTag == "editable" {
			editable = true
		}

		var valuePtr reflect.Value
		if fieldType.Type.Kind() == reflect.Pointer {
			if !fieldType.IsExported() {
				valuePtr = getUnexportedField(fieldValue)
			} else {
				valuePtr = fieldValue
			}
		} else {
			if !fieldType.IsExported() {
				valuePtr = getUnexportedFieldPtr(fieldValue)
			} else {
				valuePtr = fieldValue.Addr()
			}
		}

		err := renderDebugValue(fieldType.Name, valuePtr, editable)
		if err != nil {
			return fmt.Errorf("Failed to render field %v: %v", fieldType.Name, err)
		}
	}

	return nil
}
