package main

import (
	"testing"
)

func TestParseSingleRange(t *testing.T) {
	tests := []struct {
		input   string
		want    Range
		wantErr bool
	}{
		// Single line
		{"5", Range{5, 5, true}, false},
		{"1", Range{1, 1, true}, false},
		{"999999", Range{999999, 999999, true}, false},

		// Inclusive range
		{"3-10", Range{3, 10, true}, false},
		{"1-1", Range{1, 1, true}, false},
		{"1-100", Range{1, 100, true}, false},

		// From start
		{"...50", Range{1, 50, true}, false},
		{"...1", Range{1, 1, true}, false},

		// To EOF
		{"100...", Range{100, 0, false}, false},
		{"1...", Range{1, 0, false}, false},

		// Errors
		{"0", Range{}, true},
		{"-5", Range{}, true},
		{"abc", Range{}, true},
		{"5-3", Range{}, true},
		{"", Range{}, true},
		{"...", Range{}, true},
		{"5-", Range{}, true},
		{"-5-10", Range{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseSingleRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSingleRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseSingleRange(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRanges(t *testing.T) {
	tests := []struct {
		input   string
		want    []Range
		wantErr bool
	}{
		{"5,8,10-12", []Range{{5, 5, true}, {8, 8, true}, {10, 12, true}}, false},
		{"1-3,7,20...", []Range{{1, 3, true}, {7, 7, true}, {20, 0, false}}, false},
		{"...5,10-15", []Range{{1, 5, true}, {10, 15, true}}, false},
		{"42", []Range{{42, 42, true}}, false},

		// Errors
		{",5", nil, true},
		{"5,", nil, true},
		{"5,,8", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseRanges(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRanges(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("parseRanges(%q) returned %d ranges, want %d", tt.input, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseRanges(%q)[%d] = %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMatchesAny(t *testing.T) {
	tests := []struct {
		name    string
		lineNum int
		ranges  []Range
		want    bool
	}{
		{"in single", 5, []Range{{5, 5, true}}, true},
		{"not in single", 6, []Range{{5, 5, true}}, false},
		{"in range", 7, []Range{{5, 10, true}}, true},
		{"before range", 4, []Range{{5, 10, true}}, false},
		{"after range", 11, []Range{{5, 10, true}}, false},
		{"in open end", 100, []Range{{5, 0, false}}, true},
		{"before open end", 4, []Range{{5, 0, false}}, false},
		{"multi range hit", 8, []Range{{1, 3, true}, {8, 8, true}}, true},
		{"multi range miss", 5, []Range{{1, 3, true}, {8, 8, true}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesAny(tt.lineNum, tt.ranges); got != tt.want {
				t.Errorf("matchesAny(%d, %+v) = %v, want %v", tt.lineNum, tt.ranges, got, tt.want)
			}
		})
	}
}

func TestMaxEnd(t *testing.T) {
	tests := []struct {
		name     string
		ranges   []Range
		wantMax  int
		wantBool bool
	}{
		{"all bounded", []Range{{1, 5, true}, {10, 15, true}}, 15, true},
		{"has open end", []Range{{1, 5, true}, {10, 0, false}}, 0, false},
		{"single", []Range{{3, 3, true}}, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMax, gotBool := maxEnd(tt.ranges)
			if gotMax != tt.wantMax || gotBool != tt.wantBool {
				t.Errorf("maxEnd(%+v) = (%d, %v), want (%d, %v)", tt.ranges, gotMax, gotBool, tt.wantMax, tt.wantBool)
			}
		})
	}
}
