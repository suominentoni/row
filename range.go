package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	Start  int
	End    int
	HasEnd bool
}

func parseSingleRange(s string) (Range, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Range{}, fmt.Errorf("empty range")
	}
	if s == "..." {
		return Range{}, fmt.Errorf("invalid range: %q", s)
	}

	// ...M — from start to M
	if strings.HasPrefix(s, "...") {
		val, err := strconv.Atoi(s[3:])
		if err != nil || val < 1 {
			return Range{}, fmt.Errorf("invalid range: %q", s)
		}
		return Range{Start: 1, End: val, HasEnd: true}, nil
	}

	// N... — from N to EOF
	if strings.HasSuffix(s, "...") {
		val, err := strconv.Atoi(s[:len(s)-3])
		if err != nil || val < 1 {
			return Range{}, fmt.Errorf("invalid range: %q", s)
		}
		return Range{Start: val, End: 0, HasEnd: false}, nil
	}

	// N-M — inclusive range
	if idx := strings.Index(s, "-"); idx >= 0 {
		left := s[:idx]
		right := s[idx+1:]
		if left == "" || right == "" {
			return Range{}, fmt.Errorf("invalid range: %q", s)
		}
		start, err := strconv.Atoi(left)
		if err != nil || start < 1 {
			return Range{}, fmt.Errorf("invalid range: %q", s)
		}
		end, err := strconv.Atoi(right)
		if err != nil || end < 1 {
			return Range{}, fmt.Errorf("invalid range: %q", s)
		}
		if start > end {
			return Range{}, fmt.Errorf("invalid range: start %d > end %d", start, end)
		}
		return Range{Start: start, End: end, HasEnd: true}, nil
	}

	// N — single line
	val, err := strconv.Atoi(s)
	if err != nil || val < 1 {
		return Range{}, fmt.Errorf("invalid range: %q", s)
	}
	return Range{Start: val, End: val, HasEnd: true}, nil
}

func parseRanges(arg string) ([]Range, error) {
	parts := strings.Split(arg, ",")
	var ranges []Range
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("empty range segment in %q", arg)
		}
		r, err := parseSingleRange(p)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, r)
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})
	return ranges, nil
}

func matchesAny(lineNum int, ranges []Range) bool {
	for _, r := range ranges {
		if r.HasEnd {
			if lineNum >= r.Start && lineNum <= r.End {
				return true
			}
		} else {
			if lineNum >= r.Start {
				return true
			}
		}
	}
	return false
}

func maxEnd(ranges []Range) (int, bool) {
	max := 0
	for _, r := range ranges {
		if !r.HasEnd {
			return 0, false
		}
		if r.End > max {
			max = r.End
		}
	}
	return max, true
}
