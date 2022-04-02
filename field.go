package structquery

import (
	"reflect"
	"strings"
)

type fieldInfo struct {
	typ         reflect.Type
	name        string
	isAnonymous bool
	tag         string
	query       string
	options     map[string]string
}

type fieldWithValue struct {
	*fieldInfo
	value      interface{}
	fieldValue reflect.Value
}

func parseValue(value interface{}) ([]*fieldWithValue, error) {
	v := reflect.ValueOf(value)
	v = indirectValue(v)
	if v.Kind() != reflect.Struct {
		return nil, ErrBadQueryValue
	}
	fields := parseStruct(v)
	return fields, nil
}

func indirectValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return indirectValue(value.Elem())
	}
	return value
}

func parseStruct(value reflect.Value) []*fieldWithValue {
	valueType := value.Type()
	var fields = make([]*fieldWithValue, 0, valueType.NumField())
	for i := 0; i < valueType.NumField(); i++ {
		f := valueType.Field(i)
		fv := value.Field(i)
		if valueType.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}

		filedMeta := toFieldInfo(f)
		if filedMeta.query == "" && filedMeta.isAnonymous {
			fv = indirectValue(fv)
			if fv.Type().Kind() == reflect.Struct {
				anonymousFields := parseStruct(fv)
				fields = append(fields, anonymousFields...)
				continue
			}
		}

		if fv.IsZero() {
			continue
		}

		fields = append(fields, &fieldWithValue{
			fieldInfo:  filedMeta,
			value:      fv.Interface(),
			fieldValue: fv,
		})
	}
	return fields
}

func toFieldInfo(field reflect.StructField) *fieldInfo {
	tag := field.Tag.Get(defaultTag)
	query, options := parseTag(tag)
	return &fieldInfo{
		typ:         field.Type,
		name:        field.Name,
		isAnonymous: field.Anonymous,
		tag:         tag,
		query:       query,
		options:     options,
	}
}

type tagOptions map[string]string

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, `;`)
	m := make(tagOptions, len(s))
	for _, option := range s[1:] {
		option = strings.TrimSpace(option)
		if option == "" {
			continue
		}
		kv := strings.SplitN(option, `:`, 2)
		if len(kv) > 1 {
			m[kv[0]] = kv[1]
		} else {
			m[kv[0]] = ""
		}
	}
	return s[0], m
}
