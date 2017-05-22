package prefixspan

import (
	"reflect"
	"testing"
)

const (
	a = iota
	b
	c
	d
	e
	f
	g
)

var testMinSupport = 2

var testDB = []Sequence{
	{{a}, {a, b, c}, {a, c}, {d}, {c, f}},
	{{a, d}, {c}, {b, c}, {a, e}},
	{{c, f}, {a, b}, {d, f}, {c}, {b}},
	{{c}, {g}, {a, f}, {c}, {b}, {c}},
}

var testDBProjectedA = []Sequence{
	{{a, b, c}, {a, c}, {d}, {c, f}},
	{{placeholder, d}, {c}, {b, c}, {a, e}},
	{{placeholder, b}, {d, f}, {c}, {b}},
	{{placeholder, f}, {c}, {b}, {c}},
}

func TestPrefixSpan(t *testing.T) {
	result := PrefixSpan(testDB, testMinSupport)
	t.Log(result) // TODO
}

func TestAppendToSequence_a(t *testing.T) {
	result := appendToSequence(testDB, testMinSupport, a)

	want := testDBProjectedA

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(%v) = \n%v\n, want \n%v", a, result, want)
	}
}

func TestAppendToSequence_aa(t *testing.T) {
	result := appendToSequence(testDBProjectedA, testMinSupport, a)

	want := []Sequence{
		{{placeholder, b, c}, {a, c}, {d}, {c, f}},
		{{placeholder, e}},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(appendToSequence(%v), %v) = \n%v\n, want \n%v", a, a, result, want)
	}
}

func TestAppendToSequence_ab(t *testing.T) {
	result := appendToSequence(testDBProjectedA, testMinSupport, b)

	want := []Sequence{
		{{placeholder, c}, {a, c}, {d}, {c, f}},
		{{placeholder, c}, {a, e}}, // Typo in paper?
		{{c}},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(appendToSequence(%v), %v) = \n%v\n, want \n%v", a, b, result, want)
	}
}

func TestAppendToSequence_b(t *testing.T) {
	result := appendToSequence(testDB, testMinSupport, b)

	want := []Sequence{
		{{placeholder, c}, {a, c}, {d}, {c, f}},
		{{placeholder, c}, {a, e}},
		{{d, f}, {c}, {b}},
		{{c}},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(%v) = \n%v\n, want \n%v", b, result, want)
	}
}

func TestAppendToSequence_c(t *testing.T) {
	result := appendToSequence(testDBProjectedA, testMinSupport, c)

	// <c> in the paper
	want := []Sequence{
		{{a, c}, {d}, {c, f}},
		{{b, c}, {a, e}},
		{{b}},
		{{b}, {c}},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(%v) = \n%v\n, want \n%v", c, result, want)
	}
}

func TestAppendToItemSet_ab(t *testing.T) {
	result := appendToItemSet(testDBProjectedA, testMinSupport, ItemSet{a, b})

	want := []Sequence{
		{{placeholder, c}, {a, c}, {d}, {c, f}},
		{{d, f}, {c}, {b}},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToItemSet(appendToItemSet(%v), %v) = \n%v\n, want \n%v", ItemSet{a}, ItemSet{a, b}, result, want)
	}
}
