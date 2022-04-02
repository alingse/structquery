package structquery

import (
	"reflect"
	"strings"
)

type fieldInfo struct {
	typ           reflect.Type
	name          string
	canonicalName string
	isAnonymous   bool
	tag           string
	query         string
	options       map[string]string
}

type fieldWithValue struct {
	*fieldInfo
	value interface{}
}

func parse(value interface{}) ([]*fieldWithValue, error) {
	v := reflect.ValueOf(value)
	v = indirectValue(v)
	if v.Kind() != reflect.Struct {
		return nil, ErrBadQueryValue
	}
	fields := parseStruct(v.Type(), v)
	return fields, nil
}

func indirectValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return indirectValue(value.Elem())
	}
	return value
}

func parseStruct(typ reflect.Type, value reflect.Value) []*fieldWithValue {
	var fields []*fieldWithValue
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fv := value.Field(i)
		if typ.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}

		filedMeta := structFieldTofiledInfo(f)
		if filedMeta.query == "" && filedMeta.isAnonymous {
			fv = indirectValue(fv)
			if fv.Type().Kind() == reflect.Struct {
				anonymousFields := parseStruct(fv.Type(), fv)
				fields = append(fields, anonymousFields...)
				continue
			}
		}

		if fv.IsZero() {
			continue
		}

		fields = append(fields, &fieldWithValue{
			fieldInfo: filedMeta,
			value:     fv.Interface(),
		})
	}
	return fields
}

func structFieldTofiledInfo(field reflect.StructField) *fieldInfo {
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
