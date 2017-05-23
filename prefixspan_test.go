package prefixspan

import (
	"reflect"
	"sort"
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
	{{e, f}, {a, b}, {d, f}, {c}, {b}},
	{{e}, {g}, {a, f}, {c}, {b}, {c}},
}

var testDBProjectedA = []Sequence{
	{{a, b, c}, {a, c}, {d}, {c, f}},
	{{placeholder, d}, {c}, {b, c}, {a, e}},
	{{placeholder, b}, {d, f}, {c}, {b}},
	{{placeholder, f}, {c}, {b}, {c}},
}

var testDBProjectedF = []Sequence{
	{},
	{{a, b}, {d, f}, {c}, {b}},
	{{c}, {b}, {c}},
}

var testDBProjectedFB = []Sequence{
	{{d, f}, {c}, {b}},
	{{c}},
}

func lessItemSet(a, b ItemSet) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return len(a) < len(b)
}

func sortDB(db []Sequence) {
	sort.Slice(db, func(i, j int) bool {
		a, b := db[i], db[j]

		for k := 0; k < len(a) && k < len(b); k++ {
			if !reflect.DeepEqual(a[k], b[k]) {
				return lessItemSet(a[k], b[k])
			}
		}
		return len(a) < len(b)
	})
}

func TestPrefixSpan(t *testing.T) {
	result := PrefixSpan(testDB, testMinSupport)

	want := []Sequence{
		{{a}},
		{{a}, {a}},
		{{a}, {b}},
		{{a}, {b, c}},
		{{a}, {b, c}, {a}},
		{{a}, {b}, {a}},
		{{a}, {b}, {c}},
		{{a, b}},
		{{a, b}, {c}},
		{{a, b}, {d}},
		{{a, b}, {f}},
		{{a, b}, {d}, {c}},
		{{a}, {c}},
		{{a}, {c}, {a}},
		{{a}, {c}, {b}},
		{{a}, {c}, {c}},
		{{a}, {d}},
		{{a}, {d}, {c}},
		{{a}, {f}},

		{{b}},
		{{b}, {a}},
		{{b}, {c}},
		{{b, c}},
		{{b, c}, {a}},
		{{b}, {d}},
		{{b}, {d}, {c}},
		{{b}, {f}},

		{{c}},
		{{c}, {a}},
		{{c}, {b}},
		{{c}, {c}},

		{{d}},
		{{d}, {b}},
		{{d}, {c}},
		{{d}, {c}, {b}},

		{{e}},
		{{e}, {a}},
		{{e}, {a}, {b}},
		{{e}, {a}, {c}},
		{{e}, {a}, {c}, {b}},
		{{e}, {b}},
		{{e}, {b}, {c}},
		{{e}, {c}},
		{{e}, {c}, {b}},
		{{e}, {f}},
		{{e}, {f}, {b}},
		{{e}, {f}, {c}},
		{{e}, {f}, {c}, {b}},

		{{f}},
		{{f}, {b}},
		{{f}, {b}, {c}},
		{{f}, {c}},
		{{f}, {c}, {b}},
	}

	sortDB(want)
	sortDB(result)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("PrefixSpan() = \n%v\n, want \n%v", result, want)
	}
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
		{{placeholder, e}}, // Typo in paper?
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
		{},
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
	result := appendToSequence(testDB, testMinSupport, c)

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

func TestAppendToSequence_f(t *testing.T) {
	result := appendToSequence(testDB, testMinSupport, f)

	want := testDBProjectedF

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(%v) = \n%v\n, want \n%v", f, result, want)
	}
}

func TestAppendToSequence_fb(t *testing.T) {
	result := appendToSequence(testDBProjectedF, testMinSupport, b)

	want := testDBProjectedFB

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(appendToSequence(%v), %v) = \n%v\n, want \n%v", f, b, result, want)
	}
}

func TestAppendToSequence_fbc(t *testing.T) {
	result := appendToSequence(testDBProjectedFB, testMinSupport, c)

	want := []Sequence{
		{{b}},
		{},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("appendToSequence(appendToSequence(appendToSequence(%v), %v), %v) = \n%v\n, want \n%v", f, b, c, result, want)
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

func TestSequence_String(t *testing.T) {
	seq := Sequence{{placeholder, c}, {a, c}, {d}, {c, f}}
	s := seq.String()
	want := "<(_c)(ac)d(cf)>"

	if s != want {
		t.Errorf("Sequence.String() = %v, want %v", s, want)
	}
}
