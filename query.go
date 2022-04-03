package structquery

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type FieldMeta struct {
	Type        reflect.Type
	QueryType   QueryType
	Tag         string
	Options     map[string]string
	FieldName   string
	IsAnonymous bool
}

type Field struct {
	FieldMeta
	ColumnName string
	Value      interface{}
	FieldValue reflect.Value
}

type QueryType string
type FieldQueryer func(Field) clause.Expression

type Queryer struct {
	Namer    schema.Namer
	queryFns map[QueryType]FieldQueryer
}

const defaultTag = `sq`

func NewQueryer() *Queryer {
	q := &Queryer{
		queryFns: make(map[QueryType]FieldQueryer),
		Namer:    schema.NamingStrategy{}, // default gorm namer
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
	fields, err := parseValue(query)
	if err != nil {
		return nil, err
	}
	return q.translate(fields)
}

func (q *Queryer) translate(fields []*fieldWithValue) ([]clause.Expression, error) {
	var exprs []clause.Expression
	for _, field := range fields {
		if field.QueryType == "" {
			continue
		}
		queryType := field.QueryType
		fn, ok := q.queryFns[queryType]
		if !ok {
			return nil, fmt.Errorf("%w:%s ", ErrBadQueryType, queryType)
		}
		meta := *field.FieldMeta

		f := Field{
			FieldMeta:  meta,
			ColumnName: q.Namer.ColumnName("", field.FieldName),
			Value:      field.value,
			FieldValue: field.fieldValue,
		}
		expr := fn(f)
		if expr != nil {
			exprs = append(exprs, expr)
		}
	}
	return exprs, nil
}
