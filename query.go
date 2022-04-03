package structquery

import (
	"errors"
	"fmt"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type QueryType string
type QueryerFunc func(f Field) clause.Expression

type Queryer struct {
	Namer    schema.Namer
	queryFns map[QueryType]QueryerFunc
}

const defaultTag = `sq`

func NewQueryer() *Queryer {
	q := &Queryer{
		queryFns: make(map[QueryType]QueryerFunc),
		Namer:    schema.NamingStrategy{}, // default gorm namer
	}
	registerBuiltin(q)
	return q
}

func (q *Queryer) Register(t QueryType, fn QueryerFunc) {
	q.queryFns[t] = fn
}

var (
	ErrBadQueryValue = errors.New("structquery: query must be a pointer to struct")
	ErrBadQueryType  = errors.New("structquery: query type not registered")
)

func (q *Queryer) And(value interface{}) (clause.Expression, error) {
	exprs, err := q.toExprs(value)
	if err != nil {
		return nil, err
	}
	return clause.And(exprs...), nil
}

func (q *Queryer) Or(value interface{}) (clause.Expression, error) {
	exprs, err := q.toExprs(value)
	if err != nil {
		return nil, err
	}
	return clause.Or(exprs...), nil
}

func (q *Queryer) toExprs(value interface{}) ([]clause.Expression, error) {
	fields, err := ParseStruct(value)
	if err != nil {
		return nil, err
	}
	return q.translate(fields)
}

func (q *Queryer) translate(fields []*Field) ([]clause.Expression, error) {
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

		f := *field
		f.ColumnName = q.Namer.ColumnName("", field.Name)

		expr := fn(f)
		if expr != nil {
			exprs = append(exprs, expr)
		}
	}
	return exprs, nil
}
