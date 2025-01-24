package hashmap

import (
	"testing"
)

func TestGetSet(t *testing.T) {
	m := New[string, string](100)
	m.Set("testing", "testing")
	if m.Get("testing") != "testing" {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	m := New[string, string](10)
	m.Set("testing", "testing")
	m.Remove("testing")
	if m.Get("testing") != "" {
		t.Fail()
	}
}

func TestAutogrow(t *testing.T) {
	m := New[string, string](2)
	m.Set("testing", "testing")
	prevCap := m.capacity
	m.Set("testing2", "testing2")

	if prevCap == m.capacity {
		t.Fatal("expected capacity to change bit it didn't")
	}
	if res := m.Get("testing"); res != "testing" {
		t.Fatalf("expected 'testing' but recieved %s", res)
	}
	if res := m.Get("testing2"); res != "testing2" {
		t.Fatalf("expected 'testing2' but recieved %s", res)
	}
}

func TestGrow(t *testing.T) {
	m := New[string, string](10)
	m.Set("testing", "testing")
	m.grow()
	val := m.Get("testing")
	if val != "testing" {
		t.Fatalf("expected 'testing' but recieved %v", val)
	}
}

func TestNewFromSlices(t *testing.T) {
	keys := []string{"1", "2", "3"}
	values := []string{"1", "2", "3"}
	m, err := NewFromSlices(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	for i := range keys {
		if res := m.Get(keys[i]); res != values[i] {
			t.Fatalf("expected %s but recieved %s", values[i], res)
		}
	}
}

func TestNewFromSlicesInt(t *testing.T) {
	keys := []int{1, 2, 3}
	values := []string{"1", "2", "3"}
	m, err := NewFromSlices(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	for i := range keys {
		if res := m.Get(keys[i]); res != values[i] {
			t.Fatalf("expected %s but recieved %s", values[i], res)
		}
	}
}
