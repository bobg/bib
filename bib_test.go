package bib

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKey(t *testing.T) {
	cases := []struct {
		inp, want string
	}{{
		inp:  "The Gumball Rally",
		want: "gumball rally",
	}, {
		inp:  "1917",
		want: "nineteen seventeen",
	}, {
		inp:  "9 to 5",
		want: "nine to 5",
	}, {
		inp:  "It's Garry Shandling's Show",
		want: "its garry shandlings show",
	}, {
		inp:  "The 40-Year-Old Virgin",
		want: "forty year old virgin",
	}, {
		inp:  "42nd Street",
		want: "forty-second street",
	}, {
		inp:  "The 30th Floor",
		want: "thirtieth floor",
	}, {
		inp:  "The 501st Legion",
		want: "five hundred first legion",
	}, {
		inp:  "The 600th Floor",
		want: "six hundredth floor",
	}, {
		inp:  "350000000 Years of Solitude",
		want: "three hundred fifty million years of solitude",
	}}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			got := Key(tc.inp)
			if got != tc.want {
				t.Errorf(`input "%s", got "%s", want "%s"`, tc.inp, got, tc.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	x := []string{
		"The Gumball Rally",
		"The 501st Legion",
		"The 600th Floor",
		"1917",
		"42nd Street",
		"350000000 Years of Solitude",
		"It's Garry Shandling's Show",
		"The 40-Year-Old Virgin",
		"The 30th Floor",
		"9 to 5",
	}
	want := []string{
		"The 501st Legion",
		"The 40-Year-Old Virgin",
		"42nd Street",
		"The Gumball Rally",
		"It's Garry Shandling's Show",
		"9 to 5",
		"1917",
		"The 600th Floor",
		"The 30th Floor",
		"350000000 Years of Solitude",
	}
	Sort(x)
	if !reflect.DeepEqual(x, want) {
		t.Errorf("got %v, want %v", x, want)
	}
}
