package cli

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// FlagFormat formats a flag
func FlagFormat(flag Flag) string {
	buffer := &bytes.Buffer{}

	formatName(buffer, flag)
	formatUsage(buffer, flag)
	formatValue(buffer, flag)
	formatEnv(buffer, flag)
	formatFile(buffer, flag)

	return buffer.String()
}

func formatName(buffer *bytes.Buffer, flag Flag) {
	def := flag.Definition()
	hide := isBool(flag.Get())

	for index, name := range split(def.Name) {
		if index > 0 {
			buffer.WriteString(", ")
		}

		if len(name) == 1 {
			buffer.WriteString("-")
		} else {
			buffer.WriteString("--")
		}

		buffer.WriteString(name)

		if !hide {
			buffer.WriteString(" value")
		}
	}

}

func formatUsage(buffer *bytes.Buffer, flag Flag) {
	def := flag.Definition()

	if def.Usage == "" {
		return
	}

	if buffer.Len() > 0 {
		buffer.WriteString("\t")
	}

	buffer.WriteString(def.Usage)
}

func formatValue(buffer *bytes.Buffer, flag Flag) {
	value := toString(flag.Get())

	if value == "" {
		return
	}

	if buffer.Len() > 0 {
		buffer.WriteString(" ")
	}

	fmt.Fprintf(buffer, "(default: %v)", value)
}

func formatEnv(buffer *bytes.Buffer, flag Flag) {
	envs := flag.Definition().EnvVar

	if envs = strings.TrimSpace(envs); envs == "" {
		return
	}

	if buffer.Len() > 0 {
		buffer.WriteString(" ")
	}

	buffer.WriteString("[")

	for index, envar := range split(envs) {
		if index > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString("$")
		buffer.WriteString(envar)
	}

	buffer.WriteString("]")
}

func formatFile(buffer *bytes.Buffer, flag Flag) {
	path := flag.Definition().FilePath

	if path = strings.TrimSpace(path); path == "" {
		return
	}

	if buffer.Len() > 0 {
		buffer.WriteString(" ")
	}

	fmt.Fprintf(buffer, "[%s]", path)
}

func split(text string) []string {
	items := strings.Split(text, ",")

	for index, item := range items {
		items[index] = strings.TrimSpace(item)
	}

	return items
}

func toString(value interface{}) string {
	v := reflect.Indirect(reflect.ValueOf(value))

	if !v.IsValid() {
		return ""
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		items := make([]string, v.Len())

		for i := 0; i < v.Len(); i++ {
			items[i] = fmt.Sprintf("%v", v.Index(i).Interface())
		}

		return strings.Join(items, ", ")
	case reflect.Map:
		items := make([]string, v.Len())

		for i, key := range v.MapKeys() {
			value := v.MapIndex(key)
			items[i] = fmt.Sprintf("%v => %v", key.Interface(), value.Interface())
		}

		return strings.Join(items, ", ")
	case reflect.Bool:
		return ""
	}

	zero := reflect.Zero(v.Type())

	if v.Interface() == zero.Interface() {
		return ""
	}

	return fmt.Sprintf("%v", v.Interface())
}

func isBool(value interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(value))
	return v.Kind() == reflect.Bool
}
