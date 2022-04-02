package structquery

import "testing"

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
		fields, err := parse(c.query)
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

	fields, err := parse(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 1)
	assertEqual(t, fields[0].name, "Email")
	assertEqual(t, fields[0].value, "hello")
}

func TestParseBaseV1(t *testing.T) {
	var q = struct {
		Email string `sq:"eq"`
		Name  string `sq:"like"`
	}{
		Email: "hello",
	}

	fields, err := parse(q)
	assertEqual(t, err, nil)
	assertEqual(t, len(fields), 1)
	assertEqual(t, fields[0].name, "Email")
	assertEqual(t, fields[0].value, "hello")
}
