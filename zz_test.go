package structquery

import (
	"reflect"
	"testing"
)

func assertTrue(t *testing.T, b bool, messages ...interface{}) {
	t.Helper()
	if !b {
		t.Errorf("assertTrue faild with messages %#v", messages)
	}
}

func assertEqual(t *testing.T, a, b interface{}, messages ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("assertEqual faild: %#v != %#v with messages: %#v", a, b, messages)
	}
}

func assertNotEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Errorf("assertNotEqual faild: %#v equal to %#v", a, b)
	}
}
