package structquery

import (
	"fmt"
	"reflect"
	"testing"
)

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
