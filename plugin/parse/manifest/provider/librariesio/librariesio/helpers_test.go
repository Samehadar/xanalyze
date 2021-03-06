package librariesio

import (
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	b := true
	p := new(bool)

	if p = Bool(b); *p != b {
		t.Errorf("Bool did not return a *bool, got %v", p)
	}
}

func TestString(t *testing.T) {
	s := "HelloWorld"
	p := new(string)

	if p = String(s); *p != s {
		t.Errorf("String did not return a *string, got %v", p)
	}
}

func TestInt(t *testing.T) {
	i := 1234
	p := new(int)

	if p = Int(i); *p != i {
		t.Errorf("Int did not return a *int, got %v", p)
	}
}

func TestTime(t *testing.T) {
	i := time.Time{}
	p := new(time.Time)

	if p = Time(i); *p != i {
		t.Errorf("Time did not return a *time.Time, got %v", p)
	}
}
