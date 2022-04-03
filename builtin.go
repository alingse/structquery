package structquery

import (
	"fmt"

	"gorm.io/gorm/clause"
)

const (
	Empty             QueryType = ""
	Eq                QueryType = "eq"
	Neq               QueryType = "neq"
	Like              QueryType = "like"
	LLike             QueryType = "llike"
	RLike             QueryType = "rlike"
	In                QueryType = "in"
	NotIn             QueryType = "not_in"
	Gt                QueryType = "gt"
	Gte               QueryType = "gte"
	Lt                QueryType = "lt"
	Lte               QueryType = "lte"
	JSONExtractEq     QueryType = "json_extract_eq"
	JSONExtractLike   QueryType = "json_extract_like"
	MySQLJSONContains QueryType = "my_json_contains"
	UnsafeRawSQL      QueryType = "unsaferaw" // dangerous
)

func registerBuiltin(q *Queryer) {
	q.Register(Empty, QueryEq)
	q.Register(Eq, QueryEq)
	q.Register(Neq, func(f Field) clause.Expression {
		return clause.Neq{
			Column: f.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(Like, func(f Field) clause.Expression {
		return clause.Like{
			Column: f.ColumnName,
			Value:  `%` + fmt.Sprintf("%v", f.Value) + `%`,
		}
	})
	q.Register(LLike, func(f Field) clause.Expression {
		return clause.Like{
			Column: f.ColumnName,
			Value:  `%` + fmt.Sprintf("%v", f.Value),
		}
	})
	q.Register(RLike, func(f Field) clause.Expression {
		return clause.Like{
			Column: f.ColumnName,
			Value:  fmt.Sprintf("%v", f.Value) + `%`,
		}
	})
	q.Register(In, QueryEq)
	q.Register(NotIn, QueryNeq)
	q.Register(Gt, func(f Field) clause.Expression {
		return clause.Gt{
			Column: f.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(Gte, func(f Field) clause.Expression {
		return clause.Gte{
			Column: f.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(Lt, func(f Field) clause.Expression {
		return clause.Lt{
			Column: f.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(Lte, func(f Field) clause.Expression {
		return clause.Lte{
			Column: f.ColumnName,
			Value:  f.Value,
		}
	})
	q.Register(JSONExtractEq, QueryJSONExtractEq)
	q.Register(JSONExtractLike, QueryJSONExtractLike)
	// support for mysql
	q.Register(MySQLJSONContains, QueryMySQLJSONContains)
	// unsafe raw sql
	q.Register(UnsafeRawSQL, QueryUnsafeRawSQL)
}

func QueryEq(field Field) clause.Expression {
	return clause.Eq{
		Column: field.ColumnName,
		Value:  field.Value,
	}
}

func QueryNeq(field Field) clause.Expression {
	return clause.Neq{
		Column: field.ColumnName,
		Value:  field.Value,
	}
}

func QueryJSONExtractEq(field Field) clause.Expression {
	jsonPath := field.FieldMeta.Options["path"]
	if jsonPath == "" {
		return nil
	}
	sql := fmt.Sprintf("JSON_EXTRACT(%s, '%s') = ?", field.ColumnName, jsonPath)
	var values = []interface{}{field.Value}
	return clause.NamedExpr{
		SQL:  sql,
		Vars: values,
	}
}

func QueryJSONExtractLike(field Field) clause.Expression {
	jsonPath := field.FieldMeta.Options["path"]
	if jsonPath == "" {
		return nil
	}
	sql := fmt.Sprintf("JSON_EXTRACT(%s, '%s') LIKE ?", field.ColumnName, jsonPath)
	var values = []interface{}{`%` + fmt.Sprintf("%v", field.Value) + `%`}
	return clause.NamedExpr{
		SQL:  sql,
		Vars: values,
	}
}

func QueryMySQLJSONContains(field Field) clause.Expression {
	jsonPath := field.FieldMeta.Options["path"]

	var sql string
	if jsonPath == "" {
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?)", field.ColumnName)
	} else {
		sql = fmt.Sprintf("JSON_CONTAINS(%s, ?, '%s')", field.ColumnName, jsonPath)
	}
	var values = []interface{}{field.Value}
	return clause.NamedExpr{
		SQL:  sql,
		Vars: values,
	}
}

func QueryUnsafeRawSQL(field Field) clause.Expression {
	return clause.NamedExpr{
		SQL:  field.FieldMeta.Options["sql"],
		Vars: []interface{}{field.Value},
	}
}
