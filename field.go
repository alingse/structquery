package structquery

import (
	"reflect"
	"strings"
	"sync"
)

type FieldMeta struct {
	Type               reflect.Type
	ColumnName         string
	QueryType          QueryType
	Options            map[string]string
	FieldName          string
	FieldCanonicalName string // like A.B.C
}

type fieldInfo struct {
	typ           reflect.Type
	name          string
	canonicalName string
	isAnonymous   bool
	tag           string
	query         string
	options       map[string]string
}

type structInfo struct {
	fields []*fieldInfo
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
		if fv.IsZero() {
			continue
		}
		filedMeta := structFieldTofiledInfo(f)
		if filedMeta.query == "" && filedMeta.isAnonymous {
			fv := indirectValue(fv)
			anonymousFields := parseStruct(fv.Type(), fv)
			fields = append(fields, anonymousFields...)
		} else {
			fields = append(fields, &fieldWithValue{
				fieldInfo: filedMeta,
				value:     fv.Interface(),
			})
		}
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

// ----------------------------------------------------------------------------

// some code was copied and modified from
// https://github.com/gorilla/schema/blob/master/cache.go

// Copyright (c) 2012 Rodrigo Moraes. All rights reserved.

// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ----------------------------------------------------------------------------

// cache caches meta-data about a struct.
type cache struct {
	l   sync.RWMutex
	m   map[reflect.Type]*structInfo
	tag string
}

// newCache returns a new cache.
func newCache(tag string) *cache {
	c := cache{
		m:   make(map[reflect.Type]*structInfo),
		tag: tag,
	}
	return &c
}

// get returns a cached structInfo, creating it if necessary.
func (c *cache) get(t reflect.Type) *structInfo {
	c.l.RLock()
	info := c.m[t]
	c.l.RUnlock()
	if info == nil {
		info = c.parse(t)
		c.l.Lock()
		c.m[t] = info
		c.l.Unlock()
	}
	return info
}

func (c *cache) parse(t reflect.Type) *structInfo {
	fields := c.create(t, "")
	return &structInfo{fields: fields}
}

// create creates a structInfo with meta-data about a struct.
func (c *cache) create(t reflect.Type, parentName string) []*fieldInfo {
	var fields []*fieldInfo

	var anonymousFields []*fieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := c.createField(t.Field(i), parentName)
		fields = append(fields, f)
		if f.isAnonymous && f.tag == "" {
			ft := indirectType(f.typ)
			if ft.Kind() == reflect.Struct {
				anonymousFields = append(anonymousFields, c.create(ft, f.canonicalName)...)
			}
		}
	}

	fields = append(fields, anonymousFields...)
	return fields
}

// createField creates a fieldInfo for the given field.
func (c *cache) createField(field reflect.StructField, parentName string) *fieldInfo {
	tag := field.Tag.Get(c.tag)
	query, options := parseTag(tag)
	canonicalName := field.Name
	if parentName != "" {
		canonicalName = parentName + "." + field.Name
	}

	return &fieldInfo{
		typ:           field.Type,
		name:          field.Name,
		canonicalName: canonicalName,
		isAnonymous:   field.Anonymous,
		tag:           tag,
		query:         query,
		options:       options,
	}
}

type tagOptions map[string]string

// parseTag splits a struct field's tag into its name and map[string]string
// options.
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

func indirectType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem()
	}
	return typ
}
