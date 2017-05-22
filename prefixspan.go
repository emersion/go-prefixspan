// Package prefixspan implements the PrefixSpan algorithm, as defined in
// http://hanj.cs.illinois.edu/pdf/span01.pdf.
package prefixspan

// "_" in the paper.
const placeholder = -1

func itoa(item int) string {
	if item == placeholder {
		return "_"
	}
	r := rune(int('a') + item)
	return string(r)
}

// An ItemSet is a subset of items.
type ItemSet []int

// String implements Stringer.
func (set ItemSet) String() string {
	if len(set) == 1 {
		return itoa(set[0])
	}

	s := "("
	for _, it := range set {
		s += itoa(it)
	}
	s += ")"
	return s
}

// contains checks if an item is in an itemset.
func (set ItemSet) contains(item int) bool {
	for _, it := range set {
		if it == item {
			return true
		}
	}
	return false
}

// A sequence is an ordered list of itemsets.
type Sequence []ItemSet

// String implements Stringer.
func (seq Sequence) String() string {
	s := "<"
	for _, set := range seq {
		s += set.String()
	}
	s += ">"
	return s
}

func (seq Sequence) copy() Sequence {
	copied := make(Sequence, len(seq))
	copy(copied, seq)
	return copied
}

func (seq Sequence) suffix(i, j int) Sequence {
	suffix := seq[i:]
	if len(suffix) == 0 {
		return nil
	}
	itemSetTrail := suffix[0][j+1:]
	if len(itemSetTrail) == 0 {
		suffix = suffix[1:]
	} else {
		suffix = suffix.copy()
		suffix[0] = append(ItemSet{placeholder}, itemSetTrail...)
	}
	return suffix
}

func (seq Sequence) sequenceSuffix(item int) Sequence {
	for i, set := range seq {
		for j, it := range set {
			if it == item {
				if i == 0 && j > 0 && set[0] == placeholder {
					continue
				}
				return seq.suffix(i, j)
			}
		}
	}
	return nil
}

func (seq Sequence) itemSetSuffix(itemSet ItemSet) Sequence {
	lastItem := itemSet[len(itemSet)-1]
	for i, set := range seq {
		if set[0] == placeholder && set[1] == lastItem {
			return seq.suffix(i, 1)
		}

		delta := 0
		for j, it := range set {
			if it == itemSet[delta] {
				delta++
				if delta == len(itemSet) {
					return seq.suffix(i, j)
				}
			} else {
				delta = 0
			}
		}
	}
	return nil
}

func frequentItems(db []Sequence, minSupport int) ItemSet {
	var list ItemSet
	m := make(map[int]int)
	for _, seq := range db {
		for _, itemSet := range seq {
			for _, item := range itemSet {
				if item == placeholder {
					continue
				}

				m[item]++

				if m[item] == minSupport {
					list = append(list, item)
				}
			}
		}
	}
	return list
}

func appendToSequence(db []Sequence, minSupport int, item int) []Sequence {
	var projected []Sequence
	for _, seq := range db {
		suffix := seq.sequenceSuffix(item)
		if len(suffix) == 0 {
			continue
		}

		projected = append(projected, suffix)
	}
	return projected
}

func appendToItemSet(db []Sequence, minSupport int, itemSet ItemSet) []Sequence {
	var projected []Sequence
	for _, seq := range db {
		suffix := seq.itemSetSuffix(itemSet)
		if len(suffix) == 0 {
			continue
		}

		projected = append(projected, suffix)
	}
	return projected
}

func prefixSpan(db []Sequence, minSupport int, pattern Sequence) []Sequence {
	// Scan db once, find the set of frequent items
	items := frequentItems(db, minSupport)

	var patterns []Sequence
	if len(pattern) > 0 {
		patterns = append(patterns, pattern)
	}

	for _, item := range items {
		// Append item to the last element of pattern to form a sequential pattern
		if len(pattern) > 0 {
			lastItemSet := pattern[len(pattern)-1]
			if !lastItemSet.contains(item) {
				lastItemSet = append(lastItemSet, item)
				projected := appendToItemSet(db, minSupport, lastItemSet)
				if len(projected) >= minSupport {
					superPattern := pattern.copy()
					superPattern[len(superPattern)-1] = lastItemSet
					patterns = append(patterns, prefixSpan(projected, minSupport, superPattern)...)
				}
			}
		}

		// Append item to pattern to form a sequential pattern
		projected := appendToSequence(db, minSupport, item)
		if len(projected) >= minSupport {
			superPattern := append(pattern, ItemSet{item})
			patterns = append(patterns, prefixSpan(projected, minSupport, superPattern)...)
		}
	}

	return patterns
}

// PrefixSpan searches frequent sequences in the database db. Sequences must
// only contain positive values.
func PrefixSpan(db []Sequence, minSupport int) []Sequence {
	return prefixSpan(db, minSupport, nil)
}
