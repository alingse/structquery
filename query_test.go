package structquery

import (
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
