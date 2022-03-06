package structquery

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Field struct {
	FieldMeta
	Value interface{}
}

type QueryType string
type FieldQueryer func(Field) clause.Expression

type Queryer struct {
	cache    *cache
	queryFns map[QueryType]FieldQueryer
	namer    schema.Namer
}

const defaultTag = `sq`

func NewQueryer() *Queryer {
	q := &Queryer{
		cache:    newCache(defaultTag),
		queryFns: make(map[QueryType]FieldQueryer),
		namer:    schema.NamingStrategy{},
	}
	RegisterBuiltin(q)
	return q
}

func (q *Queryer) Register(qt QueryType, fn FieldQueryer) {
	q.queryFns[qt] = fn
}

var (
	ErrBadQueryValue = errors.New("structquery: query must be a pointer to struct")
	ErrBadQueryType  = errors.New("structquery: query type not registered")
)

func (q *Queryer) And(queryValue interface{}) (clause.Expression, error) {
	exprs, err := q.toExprs(queryValue)
	if err != nil {
		return nil, err
	}
	return clause.And(exprs...), nil
}

func (q *Queryer) Or(queryValue interface{}) (clause.Expression, error) {
	exprs, err := q.toExprs(queryValue)
	if err != nil {
		return nil, err
	}
	return clause.Or(exprs...), nil
}

func (q *Queryer) toExprs(query interface{}) ([]clause.Expression, error) {
	v := reflect.ValueOf(query)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, ErrBadQueryValue
	}
	queryValue := v.Elem()
	stInfo := q.cache.get(queryValue.Type())
	return q.bindStructInfo(stInfo, queryValue)
}

func (q *Queryer) bindStructInfo(st *structInfo, queryValue reflect.Value) ([]clause.Expression, error) {
	var exprs []clause.Expression
	for _, field := range st.fields {
		queryType := QueryType(field.query)
		fn, ok := q.queryFns[queryType]
		if !ok {
			return nil, fmt.Errorf("%w:%s ", ErrBadQueryType, queryType)
		}
		valueField := queryValue.Field(0)
		if valueField.IsZero() {
			continue
		}
		value := valueField.Interface()
		meta := FieldMeta{
			Type:       field.typ,
			ColumnName: q.namer.ColumnName("", field.name),
			QueryType:  queryType,
			Options:    field.options,
		}
		expr := fn(Field{FieldMeta: meta, Value: value})
		if expr != nil {
			exprs = append(exprs, expr)
		}
	}
	return exprs, nil
}
