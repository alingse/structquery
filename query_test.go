package structquery

import (
	"errors"
	"testing"
)

type QueryA struct {
	QueryB
	A int `json:"a" sq:"eq;column:a;"`
}

type QueryB struct {
	QueryC
	name string `sq:""`
}

type QueryC struct {
	// *QueryA
	name string `sq:"eq;json_path:xxx;"`
}

func TestParseTag(t *testing.T) {
	qt, options := parseTag("")
	assertEqual(t, qt, "")
	assertEqual(t, len(options), 0)

	qt, options = parseTag(`eq;a:b;type:varchar(100);uniq;`)
	assertEqual(t, qt, "eq")
	assertEqual(t, len(options), 3)
	assertEqual(t, options["a"], "b")
	assertEqual(t, options["type"], "varchar(100)")
	assertEqual(t, options["uniq"], "")
}

func TestNewQueryer(t *testing.T) {
	queryer := NewQueryer()
	assertNotEqual(t, queryer, nil)
}

type UserQuery struct {
	Name           string `sq:"like"`
	Email          string `sq:"eq"`
	ID             int    `sq:""`
	CreatedAtStart int64  `sq:"gte;column:created_at"`
	CreatedAtEnd   int64  `sq:"lt;column:created_at"`
	Tag            string `sq:"json_extract_like;column:tags;path:$[*].name"`
}

func TestNewQueryerWithAnd(t *testing.T) {
	queryer := NewQueryer()
	assertNotEqual(t, queryer, nil)
	var q = UserQuery{
		Name:           "hello",
		Email:          "",
		ID:             1,
		CreatedAtStart: 1000000,
		CreatedAtEnd:   2000000,
		Tag:            "gorm",
	}
	exprs, err := queryer.toExprs(&q)
	assertEqual(t, err, nil)
	assertNotEqual(t, exprs, nil)

	expr, err := queryer.And(&q)
	assertEqual(t, err, nil)
	assertNotEqual(t, expr, nil)

	expr, err = queryer.Or(&q)
	assertEqual(t, err, nil)
	assertNotEqual(t, expr, nil)
}

func TestNewQueryerWithError(t *testing.T) {
	queryer := NewQueryer()

	var q UserQuery
	_, err := queryer.toExprs(q)
	assertEqual(t, err, nil)

	_, err = queryer.And(q)
	assertEqual(t, err, nil)

	_, err = queryer.Or(q)
	assertEqual(t, err, nil)
}

func TestMoreBuiltin(t *testing.T) {
	queryer := NewQueryer()
	var q = struct {
		Tag            string `sq:"my_json_contains;column:tags;path:$.name"`
		Location       string `sq:"json_extract_eq;column:locations;path:$.name"`
		FilterID       int    `sq:"neq;column:id"`
		NameStart      string `sq:"llike;column:name"`
		NameEnd        string `sq:"rlike;column:name"`
		LastVisitStart int64  `sq:"gt;column:last_visit"`
		LastVisitEnd   int64  `sq:"lte;column:last_visit"`
	}{
		Tag:            "gorm",
		Location:       "Tokyo",
		FilterID:       100,
		NameStart:      "h",
		NameEnd:        "d",
		LastVisitStart: 1000000,
		LastVisitEnd:   2000000,
	}
	expr, err := queryer.And(&q)
	assertEqual(t, err, nil)
	assertNotEqual(t, expr, nil)
}

func TestQueryWithBadQueryType(t *testing.T) {
	queryer := NewQueryer()

	var q2 = struct {
		Name string `sq:"not_exist"`
	}{
		Name: "hello",
	}
	var err error

	_, err = queryer.toExprs(&q2)
	assertEqual(t, errors.Is(err, ErrBadQueryType), true, err)
	_, err = queryer.And(q2)
	assertEqual(t, errors.Is(err, ErrBadQueryType), true, err)
	_, err = queryer.Or(q2)
	assertEqual(t, errors.Is(err, ErrBadQueryType), true, err)
}

func TestQueryWithBadQueryValue(t *testing.T) {
	queryer := NewQueryer()
	var q2 = "hello"
	var err error

	_, err = queryer.toExprs(q2)
	assertEqual(t, errors.Is(err, ErrBadQueryValue), true, err)
	_, err = queryer.And(q2)
	assertEqual(t, errors.Is(err, ErrBadQueryValue), true, err)
	_, err = queryer.Or(q2)
	assertEqual(t, errors.Is(err, ErrBadQueryValue), true, err)
	_, err = queryer.Where(nil, q2)
	assertEqual(t, errors.Is(err, ErrBadQueryValue), true, err)
}
