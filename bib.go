// Package bib contains functions for (English-language) bibliographic sorting of strings.
//
// A bibliographic sort is one that ignores a leading article ("the," "a," "an")
// and treats leading numbers as if they're spelled out.
// Characters other than letters and digits are ignored,
// except that "&" is converted to the spelled-out word "and,"
// and hyphens are converted to spaces.
package bib

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/bobg/go-generics/v4/slices"
)

// Less tells whether a comes before b in a bibliograhic sort.
func Less(a, b string) bool {
	return Key(a) < Key(b)
}

// Sort sorts the input slice bibliographically.
func Sort(strs []string) {
	// We could just write:
	//
	//   sort.Slice(strs, func(i, j int) bool { return Less(strs[i], strs[j]) })
	//
	// but that would call Key on each string in strs more than once, on average,
	// which is inefficient.
	// So instead we compute keys for all the strings exactly once into a new slice,
	// then use slices.KeyedSort.

	keys := slices.Map(strs, Key)
	slices.KeyedSort(strs, sort.StringSlice(keys))
}

// Key converts an input string to a bibliographic sort key.
//
// TODO: Handle Unicode character categories.
// E.g. don't convert only "-" to a space,
// convert everything in the "dash punctuation" category
// (https://www.compart.com/en/unicode/category/Pd).
func Key(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "&", " and ")
	s = strings.ReplaceAll(s, "-", " ")

	// Keep only letters, digits, and whitespace.
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)

	f := strings.Fields(s)
	switch f[0] {
	case "a", "the", "an":
		if len(f) == 1 {
			// Unlikely case.
			return f[0]
		}
		f = f[1:]
	}

	m := numRegex.FindStringSubmatch(f[0])
	if len(m) > 0 {
		n, _ := strconv.ParseInt(m[1], 10, 64)
		f = slices.ReplaceN(f, 0, 1, intToWords(n, len(m[2]) > 0)...)
	}

	return strings.Join(f, " ")
}

var numRegex = regexp.MustCompile(`^(\d+)(st|nd|rd|th)?$`)

func intToWords(n int64, ordinal bool) []string {
	if ordinal && n < 10 {
		var x string

		switch n {
		case 0:
			x = "zeroth"
		case 1:
			x = "first"
		case 2:
			x = "second"
		case 3:
			x = "third"
		case 5:
			x = "fifth"
		case 8:
			x = "eighth"
		case 9:
			x = "ninth"
		default:
			w := intToWords(n, false)
			x = w[0] + "th"
		}
		return []string{x}
	}

	if n < 20 {
		var x string

		switch n {
		case 0:
			x = "zero"
		case 1:
			x = "one"
		case 2:
			x = "two"
		case 3:
			x = "three"
		case 4:
			x = "four"
		case 5:
			x = "five"
		case 6:
			x = "six"
		case 7:
			x = "seven"
		case 8:
			x = "eight"
		case 9:
			x = "nine"
		case 10:
			x = "ten"
		case 11:
			x = "eleven"
		case 12:
			x = "twelve"
		case 13:
			x = "thirteen"
		case 14:
			x = "fourteen"
		case 15:
			x = "fifteen"
		case 16:
			x = "sixteen"
		case 17:
			x = "seventeen"
		case 18:
			x = "eighteen"
		case 19:
			x = "nineteen"
		}
		if ordinal {
			// Won't be true for 0 through 9, which are handled above.
			x += "th"
		}

		return []string{x}
	}

	if ordinal {
		var x string

		switch n {
		case 20:
			x = "twentieth"
		case 30:
			x = "thirtieth"
		case 40:
			x = "fortieth"
		case 50:
			x = "fiftieth"
		case 60:
			x = "sixtieth"
		case 70:
			x = "seventieth"
		case 80:
			x = "eightieth"
		case 90:
			x = "ninetieth"
		}

		if x != "" {
			return []string{x}
		}
	}

	if n < 100 {
		var s string
		switch {
		case n < 30:
			s = "twenty"
		case n < 40:
			s = "thirty"
		case n < 50:
			s = "forty"
		case n < 60:
			s = "fifty"
		case n < 70:
			s = "sixty"
		case n < 80:
			s = "seventy"
		case n < 90:
			s = "eighty"
		default:
			s = "ninety"
		}
		if r := n % 10; r > 0 {
			w := intToWords(r, ordinal)
			s += "-" + w[0]
		}
		return []string{s}
	}

	if n < 1000 {
		q, r := n/100, n%100 // quotient, remainder
		w := intToWords(q, false)
		w = append(w, "hundred")
		if r > 0 {
			ww := intToWords(r, ordinal)
			w = append(w, ww...)
		}
		if ordinal && r == 0 {
			w[len(w)-1] += "th"
		}
		return w
	}

	// Years.
	if !ordinal && n >= 1100 && n < 3000 && !((n >= 2000) && (n < 2010)) {
		q, r := n/100, n%100
		w := intToWords(q, false)
		if r < 10 {
			w = append(w, "hundred")
		}
		if r > 0 {
			ww := intToWords(r, false)
			w = append(w, ww...)
		}
		return w
	}

	if n < 1000000 {
		q, r := n/1000, n%1000
		w := intToWords(q, false)
		w = append(w, "thousand")
		if r > 0 {
			ww := intToWords(r, ordinal)
			w = append(w, ww...)
		}
		if ordinal && n%100 == 0 {
			w[len(w)-1] += "th"
		}
		return w
	}

	if n < 1000000000 {
		q, r := n/1000000, n%1000000
		w := intToWords(q, false)
		w = append(w, "million")
		if r > 0 {
			ww := intToWords(r, ordinal)
			w = append(w, ww...)
		}
		if ordinal && n%100 == 0 {
			w[len(w)-1] += "th"
		}
		return w
	}

	q, r := n/1000000000, n%1000000000
	w := intToWords(q, false)
	w = append(w, "billion")
	if r > 0 {
		ww := intToWords(r, ordinal)
		w = append(w, ww...)
	}
	if ordinal && n%100 == 0 {
		w[len(w)-1] += "th"
	}
	return w
}
