// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     GoRuntime

package cog

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/grafana/cog/testdata/generated/cog/variants"
)

var runtimeInstance *Runtime

type Runtime struct {
	panelcfgVariants  map[string]variants.PanelcfgConfig
	dataqueryVariants map[string]variants.DataqueryConfig
}

func NewRuntime() *Runtime {
	if runtimeInstance != nil {
		return runtimeInstance
	}

	runtimeInstance = &Runtime{
		panelcfgVariants:  make(map[string]variants.PanelcfgConfig),
		dataqueryVariants: make(map[string]variants.DataqueryConfig),
	}

	return runtimeInstance
}

func (runtime *Runtime) RegisterPanelcfgVariant(config variants.PanelcfgConfig) {
	runtime.panelcfgVariants[config.Identifier] = config
}

func (runtime *Runtime) ConfigForPanelcfgVariant(identifier string) (variants.PanelcfgConfig, bool) {
	config, found := runtime.panelcfgVariants[identifier]

	return config, found
}

func (runtime *Runtime) RegisterDataqueryVariant(config variants.DataqueryConfig) {
	runtime.dataqueryVariants[config.Identifier] = config
}

func (runtime *Runtime) UnmarshalDataqueryArray(raw []byte, dataqueryTypeHint string) ([]variants.Dataquery, error) {
	rawDataqueries := []json.RawMessage{}
	if err := json.Unmarshal(raw, &rawDataqueries); err != nil {
		return nil, err
	}

	dataqueries := make([]variants.Dataquery, 0, len(rawDataqueries))
	for _, rawDataquery := range rawDataqueries {
		dataquery, err := runtime.UnmarshalDataquery(rawDataquery, dataqueryTypeHint)
		if err != nil {
			return nil, err
		}

		dataqueries = append(dataqueries, dataquery)
	}

	return dataqueries, nil
}

func (runtime *Runtime) UnmarshalDataquery(raw []byte, dataqueryTypeHint string) (variants.Dataquery, error) {
	// A hint tells us the dataquery type: let's use it.
	if dataqueryTypeHint != "" {
		config, found := runtime.dataqueryVariants[dataqueryTypeHint]
		if found {
			dataquery, err := config.DataqueryUnmarshaler(raw)
			if err != nil {
				return nil, err
			}

			return dataquery.(variants.Dataquery), nil
		}
	}

	// Dataqueries might reference the datasource to use, and its type. Let's use that.
	partialDataquery := struct {
		Datasource struct {
			Type string `json:"type"`
		} `json:"datasource"`
	}{}
	if err := json.Unmarshal(raw, &partialDataquery); err != nil {
		return nil, err
	}
	if partialDataquery.Datasource.Type != "" {
		config, found := runtime.dataqueryVariants[partialDataquery.Datasource.Type]
		if found {
			dataquery, err := config.DataqueryUnmarshaler(raw)
			if err != nil {
				return nil, err
			}

			return dataquery.(variants.Dataquery), nil
		}
	}

	// We have no idea what type the dataquery is: use our `UnknownDataquery` bag to not lose data.
	dataquery := variants.UnknownDataquery{}
	if err := json.Unmarshal(raw, &dataquery); err != nil {
		return nil, err
	}

	return dataquery, nil
}

func UnmarshalDataqueryArray(raw []byte, dataqueryTypeHint string) ([]variants.Dataquery, error) {
	return NewRuntime().UnmarshalDataqueryArray(raw, dataqueryTypeHint)
}

func UnmarshalDataquery(raw []byte, dataqueryTypeHint string) (variants.Dataquery, error) {
	return NewRuntime().UnmarshalDataquery(raw, dataqueryTypeHint)
}

func ConfigForPanelcfgVariant(identifier string) (variants.PanelcfgConfig, bool) {
	return NewRuntime().ConfigForPanelcfgVariant(identifier)
}

func (runtime *Runtime) ConvertPanelToGo(inputPanel any, panelType string) string {
	config, found := runtime.panelcfgVariants[panelType]
	if found && config.GoConverter != nil {
		return config.GoConverter(inputPanel)
	}

	return "/* could not convert panel to go */"
}

func (runtime *Runtime) ConvertDataqueryToGo(inputPanel any, dataqueryTypeHints ...string) string {
	for _, dataqueryTypeHint := range dataqueryTypeHints {
		config, found := runtime.dataqueryVariants[dataqueryTypeHint]
		if found && config.GoConverter != nil {
			return config.GoConverter(inputPanel)
		}
	}

	return "/* could not convert dataquery to go */"
}

func ConvertPanelToCode(inputPanel any, panelType string) string {
	return NewRuntime().ConvertPanelToGo(inputPanel, panelType)
}

func ConvertDataqueryToCode(inputPanel any, dataqueryTypeHints ...string) string {
	return NewRuntime().ConvertDataqueryToGo(inputPanel, dataqueryTypeHints...)
}

func Dump(root any) string {
	return dumpValue(reflect.ValueOf(root))
}

func dumpValue(value reflect.Value) string {
	if !value.IsValid() {
		return "<invalid>"
	}

	switch value.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%#v", value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", value.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%#v", value.Float())
	case reflect.String:
		return fmt.Sprintf("%#v", value.String())
	case reflect.Array, reflect.Slice:
		return dumpArray(value)
	case reflect.Map:
		return dumpMap(value)
	case reflect.Struct:
		return dumpStruct(value)
	case reflect.Interface:
		if !value.CanInterface() {
			return "<interface: can't interface>"
		}

		return Dump(value.Interface())
	case reflect.Pointer:
		if value.IsNil() {
			return "nil"
		}

		pointed := value.Elem()

		return fmt.Sprintf("cog.ToPtr[%s](%s)", pointed.Type().String(), dumpValue(pointed))
	default:
		return fmt.Sprintf("<unknown: type=%s, kind=%s>", value.Type(), value.Kind().String())
	}
}

func dumpArray(value reflect.Value) string {
	if value.IsNil() {
		return "nil"
	}

	parts := make([]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		parts = append(parts, dumpValue(value.Index(i)))
	}

	return fmt.Sprintf("%s{%s}", value.Type().String(), strings.Join(parts, ", "))
}

func dumpMap(value reflect.Value) string {
	if value.IsNil() {
		return "nil"
	}

	parts := make([]string, 0, value.Len())
	iter := value.MapRange()
	for iter.Next() {
		line := fmt.Sprintf("%s: %s", dumpValue(iter.Key()), dumpValue(iter.Value()))
		parts = append(parts, line)
	}

	return fmt.Sprintf("%s{%s}", value.Type().String(), strings.Join(parts, ", "))
}

func dumpStruct(value reflect.Value) string {
	parts := make([]string, 0, value.NumField())
	structType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := structType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldValue := value.Field(i)
		fieldValueKind := fieldValue.Kind()
		if (fieldValueKind == reflect.Pointer || fieldValueKind == reflect.Interface || fieldValueKind == reflect.Array || fieldValueKind == reflect.Slice || fieldValueKind == reflect.Map) && fieldValue.IsNil() {
			continue
		}

		line := fmt.Sprintf("%s: %s", field.Name, dumpValue(fieldValue))
		parts = append(parts, line)
	}

	return fmt.Sprintf("%s{%s}", value.Type().String(), strings.Join(parts, ", "))
}
