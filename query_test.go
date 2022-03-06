package structquery

import (
	"errors"
	"reflect"
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

func TestCache(t *testing.T) {
	cache := newCache(defaultTag)
	var q QueryA
	q.name = "hello"
	qv := reflect.ValueOf(q)
	st := cache.parse(qv.Type())

	assertNotEqual(t, st, nil)
	assertEqual(t, len(st.fields), 5)
	assertEqual(t, st.fields[0].name, "QueryB")
	assertEqual(t, st.fields[1].name, "A")
	assertEqual(t, st.fields[2].name, "QueryC")
	assertEqual(t, st.fields[3].name, "name")
	assertEqual(t, st.fields[4].name, "name")

	assertEqual(t, st.fields[0].canonicalName, "QueryB")
	assertEqual(t, st.fields[1].canonicalName, "A")
	assertEqual(t, st.fields[2].canonicalName, "QueryB.QueryC")
	assertEqual(t, st.fields[3].canonicalName, "QueryB.name")
	assertEqual(t, st.fields[4].canonicalName, "QueryB.QueryC.name")
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

func TestNewCacheGet(t *testing.T) {
	cache := newCache(defaultTag)
	var q QueryA
	s := cache.get(reflect.TypeOf(q))
	assertNotEqual(t, s, nil)
	s2 := cache.get(reflect.TypeOf(q))
	assertEqual(t, s, s2)
}

func TestIndirectType(t *testing.T) {
	var q *QueryA
	t1 := indirectType(reflect.TypeOf(q))
	var q2 QueryA
	t2 := indirectType(reflect.TypeOf(q2))
	assertEqual(t, t1, t2)
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
		ID:             0,
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
	assertEqual(t, err, ErrBadQueryValue)

	_, err = queryer.And(q)
	assertEqual(t, err, ErrBadQueryValue)

	_, err = queryer.Or(q)
	assertEqual(t, err, ErrBadQueryValue)

	var q2 struct {
		Name string `sq:"not_exist"`
	}
	_, err = queryer.toExprs(&q2)
	assertTrue(t, err != nil)
	assertTrue(t, errors.Is(err, ErrBadQueryType))
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
