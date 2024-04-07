package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/cog/generated/cog"
)

type SomeStruct struct {
	Id string `json:"id"`

	Title         *string
	LeaveMeNilPlz *string

	Options any

	private string
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

func main() {
	spew.Dump(Dump("foo"))
	spew.Dump(Dump(4))
	spew.Dump(Dump(-4))
	spew.Dump(Dump(4.3))
	spew.Dump(Dump(-4.3))
	spew.Dump(Dump(true))
	spew.Dump(Dump(false))

	spew.Dump(Dump([]string{"foo", "bar", "baz"}))
	spew.Dump(Dump([]int{1, 2, 3}))
	spew.Dump(Dump([][]string{
		{"foo", "bar"},
		{"baz"},
	}))

	spew.Dump(Dump(map[string]string{
		"foo": "bar",
		"baz": "biz",
	}))
	spew.Dump(Dump(map[int]string{
		3: "333",
		4: "444",
	}))

	spew.Dump(Dump(SomeStruct{
		Id:      "some-id",
		Title:   cog.ToPtr("some-title"),
		Options: "foo",
		private: "private stuff is private",
	}))

	spew.Dump(Dump(struct {
		Id string `json:"id"`

		Title         *string
		LeaveMeNilPlz *string

		Options any

		private string
	}{
		Id:      "some-id",
		Title:   cog.ToPtr("some-title"),
		Options: "foo",
		private: "private stuff is private",
	}))
}
