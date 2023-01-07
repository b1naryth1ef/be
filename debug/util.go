package debug

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/inkyblackness/imgui-go/v4"
)

type CustomEditableTypeHandler = func(label string, value reflect.Value)

var customEditableTypeRegistry = map[reflect.Type]CustomEditableTypeHandler{}

func RegisterCustomEditableType(value interface{}, handler CustomEditableTypeHandler) {
	valueType := reflect.ValueOf(value).Type()
	customEditableTypeRegistry[valueType] = handler
}

func getUnexportedField(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}

func isSigned(k reflect.Kind) bool {
	return (k == reflect.Int8 ||
		k == reflect.Int16 ||
		k == reflect.Int32 ||
		k == reflect.Int64 ||
		k == reflect.Int)
}

func isNumber(k reflect.Kind) bool {
	return (k == reflect.Float32 ||
		k == reflect.Float64 ||
		k == reflect.Int8 ||
		k == reflect.Int16 ||
		k == reflect.Int32 ||
		k == reflect.Int64 ||
		k == reflect.Int ||
		k == reflect.Uint ||
		k == reflect.Uint8 ||
		k == reflect.Uint16 ||
		k == reflect.Uint32 ||
		k == reflect.Uint64)
}

func renderEditable(name string, value reflect.Value) {
	valueType := value.Type()

	if handler, ok := customEditableTypeRegistry[valueType]; ok {
		handler(name, value)
		return
	}

	if valueType.Kind() == reflect.String || isNumber(valueType.Kind()) {
		var contents = fmt.Sprintf("%v", value.Interface())
		imgui.InputTextV(name, &contents, imgui.ImGuiInputTextFlagsCallbackEdit, func(data imgui.InputTextCallbackData) int32 {
			if valueType.Kind() == reflect.String {
				value.SetString(string(data.Buffer()))
			} else if valueType.Kind() == reflect.Float64 || valueType.Kind() == reflect.Float32 {
				rawValue, err := strconv.ParseFloat(string(data.Buffer()), valueType.Bits())
				if err == nil {
					value.SetFloat(rawValue)
				}
			} else if isNumber(valueType.Kind()) {
				rawValue, err := strconv.ParseInt(string(data.Buffer()), 10, valueType.Bits())
				if err == nil {
					if isSigned(valueType.Kind()) {
						value.SetInt(rawValue)
					} else {
						value.SetUint(uint64(rawValue))
					}
				}
			} else {
				log.Printf("Unsupported edit type %v: %s", valueType, data.Buffer())
			}
			return 0
		})
	} else if valueType.Kind() == reflect.Bool {
		checked := value.Bool()
		if imgui.Checkbox(name, &checked) {
			value.SetBool(checked)
		}
	}
}

func renderEditableStruct(value reflect.Value) {
	imgui.Text("Would render editable struct here...")
}

func RenderStruct(s interface{}) {
	structPtrValue := reflect.ValueOf(s)
	if structPtrValue.IsNil() {
		imgui.Text("nil")
		return
	}
	structValue := structPtrValue.Elem()
	structType := structValue.Type()

	for fieldIdx := 0; fieldIdx < structType.NumField(); fieldIdx++ {
		field := structType.Field(fieldIdx)

		tag := field.Tag.Get("debug")
		if tag == "-" {
			continue
		} else if tag == "editable" {
			renderEditable(field.Name, getUnexportedField(structValue.FieldByIndex(field.Index)))
		} else if tag == "struct" {
			if imgui.CollapsingHeaderV(field.Name, imgui.TreeNodeFlagsDefaultOpen) {
				RenderStruct(getUnexportedField(structValue.FieldByIndex(field.Index)).Interface())
			}
		} else if tag == "since" {
			value := getUnexportedField(structValue.FieldByIndex(field.Index)).Interface().(time.Time)
			imgui.Text(fmt.Sprintf("%s: %dms", field.Name, time.Since(value).Milliseconds()))
		} else if tag != "" {
			imgui.Text(fmt.Sprintf("%s: "+tag, field.Name, getUnexportedField(structValue.FieldByIndex(field.Index)).Interface()))
		} else {
			if field.Type.Kind() == reflect.Pointer {
				imgui.Text(fmt.Sprintf("%s: %p", field.Name, getUnexportedField(structValue.FieldByIndex(field.Index)).Interface()))
			} else {
				imgui.Text(fmt.Sprintf("%s: %v", field.Name, getUnexportedField(structValue.FieldByIndex(field.Index)).Interface()))
			}
		}
	}
}
