package structquery

import (
	"reflect"
	"strings"
)

type FieldMeta struct {
	QueryType   QueryType
	Options     map[string]string
	Name        string
	IsAnonymous bool
}

type Field struct {
	FieldMeta
	Value      interface{}
	FieldValue reflect.Value
	ColumnName string // TODO: add read from options
}

func ParseStruct(value interface{}) ([]*Field, error) {
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

func parseStruct(value reflect.Value) []*Field {
	valueType := value.Type()
	var fields = make([]*Field, 0, valueType.NumField())

	for i := 0; i < valueType.NumField(); i++ {
		f := valueType.Field(i)
		fv := value.Field(i)

		if valueType.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}

		filedMeta := toFieldInfo(f)
		if filedMeta.QueryType == "" && filedMeta.IsAnonymous {
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

		field := &Field{
			FieldMeta:  filedMeta,
			Value:      fv.Interface(),
			FieldValue: fv,
		}
		fields = append(fields, field)
	}
	return fields
}

func toFieldInfo(field reflect.StructField) FieldMeta {
	tag := field.Tag.Get(defaultTag)
	query, options := parseTag(tag)
	return FieldMeta{
		Name:        field.Name,
		IsAnonymous: field.Anonymous,
		QueryType:   QueryType(query),
		Options:     options,
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
