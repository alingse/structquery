package structquery

import (
	"fmt"

	"gorm.io/gorm/clause"
)

const (
	EmptyQueryType             QueryType = ""
	EqQueryType                QueryType = "eq"
	NeqQueryType               QueryType = "neq"
	LikeQueryType              QueryType = "like"
	LLikeQueryType             QueryType = "llike"
	RLikeQueryType             QueryType = "rlike"
	InQueryType                QueryType = "in"
	NotInQueryType             QueryType = "not_in"
	GtQueryType                QueryType = "gt"
	GteQueryType               QueryType = "gte"
	LtQueryType                QueryType = "lt"
	LteQueryType               QueryType = "lte"
	JSONExtractEqQueryType     QueryType = "json_extract_eq"
	JSONExtractLikeQueryType   QueryType = "json_extract_like"
	MySQLJSONContainsQueryType QueryType = "my_json_contains"
	RawSQLQueryType            QueryType = "unsaferaw" // dangerous
)

func RegisterBuiltin(q *Queryer) {
	q.Register(EmptyQueryType, EqQueryer)
	q.Register(EqQueryType, EqQueryer)
	q.Register(NeqQueryType, func(f Field) clause.Expression {
		return clause.Neq{
			Column: f.FieldMeta.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(LikeQueryType, func(f Field) clause.Expression {
		return clause.Like{
			Column: f.FieldMeta.ColumnName,
			Value:  `%` + fmt.Sprintf("%v", f.Value) + `%`,
		}
	})
	q.Register(LLikeQueryType, func(f Field) clause.Expression {
		return clause.Like{
			Column: f.FieldMeta.ColumnName,
			Value:  `%` + fmt.Sprintf("%v", f.Value),
		}
	})
	q.Register(RLikeQueryType, func(f Field) clause.Expression {
		return clause.Like{
			Column: f.FieldMeta.ColumnName,
			Value:  fmt.Sprintf("%v", f.Value) + `%`,
		}
	})
	q.Register(InQueryType, EqQueryer)
	q.Register(NotInQueryType, NeqQueryer)
	q.Register(GtQueryType, func(f Field) clause.Expression {
		return clause.Gt{
			Column: f.FieldMeta.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(GteQueryType, func(f Field) clause.Expression {
		return clause.Gte{
			Column: f.FieldMeta.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(LtQueryType, func(f Field) clause.Expression {
		return clause.Lt{
			Column: f.FieldMeta.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(LteQueryType, func(f Field) clause.Expression {
		return clause.Lte{
			Column: f.FieldMeta.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(JSONExtractEqQueryType, JSONExtractEqQueryer)
	q.Register(JSONExtractLikeQueryType, JSONExtractLikeQueryer)
	// support for mysql
	q.Register(MySQLJSONContainsQueryType, MySQLJSONContainsQueryer)
	// unsafe raw sql
	q.Register(RawSQLQueryType, UnsafeRawSQLQueryer)
}

func EqQueryer(field Field) clause.Expression {
	return clause.Eq{
		Column: field.FieldMeta.ColumnName,
		Value:  field.Value,
	}
}

func NeqQueryer(field Field) clause.Expression {
	return clause.Neq{
		Column: field.FieldMeta.ColumnName,
		Value:  field.Value,
	}
}

func JSONExtractEqQueryer(field Field) clause.Expression {
	jsonPath := field.FieldMeta.Options["path"]
	if jsonPath == "" {
		return nil
	}
	sql := fmt.Sprintf("JSON_EXTRACT(%s, '%s') = ?", field.FieldMeta.ColumnName, jsonPath)
	var values = []interface{}{field.Value}
	return clause.NamedExpr{
		SQL:  sql,
		Vars: values,
	}
}

func JSONExtractLikeQueryer(field Field) clause.Expression {
	jsonPath := field.FieldMeta.Options["path"]
	if jsonPath == "" {
		return nil
	}
	sql := fmt.Sprintf("JSON_EXTRACT(%s, '%s') LIKE ?", field.FieldMeta.ColumnName, jsonPath)
	var values = []interface{}{`%` + fmt.Sprintf("%v", field.Value) + `%`}
	return clause.NamedExpr{
		SQL:  sql,
		Vars: values,
	}
}

func MySQLJSONContainsQueryer(field Field) clause.Expression {
	jsonPath := field.FieldMeta.Options["path"]

	var sql string
	if jsonPath == "" {
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?)", field.FieldMeta.ColumnName)
	} else {
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?, '%s')", field.FieldMeta.ColumnName, jsonPath)
	}
	var values = []interface{}{field.Value}
	return clause.NamedExpr{
		SQL:  sql,
		Vars: values,
	}
}

func UnsafeRawSQLQueryer(field Field) clause.Expression {
	return clause.NamedExpr{
		SQL:  field.FieldMeta.Options["sql"],
		Vars: []interface{}{field.Value},
	}
}
