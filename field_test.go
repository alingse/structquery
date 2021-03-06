package structquery

import (
	"fmt"
	"reflect"
	"testing"
)

type simpleQuery struct {
	Email string `sq:"eq"`
	Name  string `sq:"like"`
}

func TestParseWithValue(t *testing.T) {
	type testCase struct {
		query  interface{}
		length int
		err    error
	}
	cases := []testCase{
		{
			query:  struct{}{},
			length: 0,
			err:    nil,
		},
		{
			query:  struct{ Email string }{},
			length: 0,
			err:    nil,
		},
		{
			query: &simpleQuery{
				Email: "hello",
			},
			length: 1,
			err:    nil,
		},
		{
			query: simpleQuery{
				Email: "hello",
			},
			length: 1,
			err:    nil,
		},
	}

	for _, c := range cases {
		fields, err := ParseStruct(c.query)
		assertEqual(t, err, c.err)
		assertEqual(t, len(fields), c.length, c)
	}
}
func TestParseBase(t *testing.T) {
	var q = struct {
		Email string `sq:"eq"`
	}{
		Email: "hello",
	}

	fields, err := ParseStruct(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 1)
	assertEqual(t, fields[0].Name, "Email")
	assertEqual(t, fields[0].Value, "hello")
}

func TestParseBaseV1(t *testing.T) {
	var q = struct {
		Email string `sq:"eq"`
		Name  string `sq:"like"`
	}{
		Email: "hello",
	}

	fields, err := ParseStruct(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 1)
	assertEqual(t, fields[0].Name, "Email")
	assertEqual(t, fields[0].Value, "hello")
}

type complexQuery struct {
	simpleQuery
	ItemID int64 `sq:"eq"`
}

func TestComplexQuery1(t *testing.T) {
	var q = complexQuery{
		simpleQuery: simpleQuery{
			Email: "hello",
		},
		ItemID: 1,
	}

	fields, err := ParseStruct(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 2)
	assertEqual(t, fields[0].Name, "Email")
	assertEqual(t, fields[0].Value, "hello")
	assertEqual(t, fields[1].Name, "ItemID")
	assertEqual(t, fields[1].Value, int64(1))
}

type complexQuery2 struct {
	*complexQuery
	Query *complexQuery `sq:"eq"`
}

func TestComplexQuery2(t *testing.T) {
	var q = complexQuery2{
		complexQuery: &complexQuery{
			simpleQuery: simpleQuery{
				Email: "hello",
				Name:  "world",
			},
		},
	}

	fields, err := ParseStruct(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 2)
	assertEqual(t, fields[0].Name, "Email")
	assertEqual(t, fields[0].Value, "hello")
	assertEqual(t, fields[1].Name, "Name")
	assertEqual(t, fields[1].Value, "world")
}

type ItemID int64
type complexQuery3 struct {
	ItemID
}

func TestComplexQuery3(t *testing.T) {
	var q = complexQuery3{
		ItemID: 0,
	}

	fields, err := ParseStruct(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 0, fields)
}

type complexQuery4 struct {
	ItemID *ItemID
}

func TestComplexQuery3Case1(t *testing.T) {
	var itemID ItemID = 0
	var q = complexQuery4{
		ItemID: &itemID,
	}

	fields, err := ParseStruct(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 1, fields)
}

func TestParseStructErr(t *testing.T) {
	var itemID int = 1
	fields, err := ParseStruct(itemID)
	assertEqual(t, err, ErrBadQueryValue)
	assertEqual(t, len(fields), 0)
}

func toMessage(messages []interface{}) string {
	if len(messages) == 0 {
		return ""
	}
	var msg string
	for _, m := range messages {
		msg += ", " + fmt.Sprintf("%#v", m)
	}
	return msg
}

func assertEqual(t *testing.T, a, b interface{}, messages ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("assertEqual faild: %#v != %#v with messages %s", a, b,
			toMessage(messages))
	}
}

func assertNotEqual(t *testing.T, a, b interface{}, messages ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Errorf("assertNotEqual faild: %#v equal to %#v with messages %s", a, b,
			toMessage(messages))
	}
}
