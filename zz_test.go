package structquery

import (
	"reflect"
	"testing"
)

func assertTrue(t *testing.T, b bool, args ...interface{}) {
	t.Helper()
	if !b {
		t.Errorf("assertTrue faild args %+v", args)
	}
}

func assertEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("assertEqual faild: %#v != %#v", a, b)
	}
}

func assertNotEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Errorf("assertNotEqual faild: %#v equal to %#v", a, b)
	}
}
