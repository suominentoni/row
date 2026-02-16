package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const usage = `Usage: row [flags] <range>

Filter lines from stdin by line number ranges.

Ranges:
  N        single line (e.g., 5)
  N-M      inclusive range (e.g., 3-10)
  ...M     from first line to M (e.g., ...50)
  N...     from N to end of input (e.g., 100...)
  A,B,C    combine ranges with commas (e.g., 1-3,7,20...)

Flags:
  -h, --hide        invert filter: hide matching lines
  -s, --separator   print --- between non-contiguous output segments
  --help             show this help`

func expandArgs(args []string) []string {
	var expanded []string
	for _, a := range args {
		if strings.HasPrefix(a, "-") && !strings.HasPrefix(a, "--") && len(a) > 2 {
			allFlags := true
			for _, c := range a[1:] {
				if c != 'h' && c != 's' {
					allFlags = false
					break
				}
			}
			if allFlags {
				for _, c := range a[1:] {
					switch c {
					case 'h':
						expanded = append(expanded, "-h")
					case 's':
						expanded = append(expanded, "-s")
					}
				}
				continue
			}
		}
		expanded = append(expanded, a)
	}
	return expanded
}

func run(args []string, input *bufio.Scanner) int {
	args = expandArgs(args)
	hide := false
	separator := false
	var rangeArg string

	for _, a := range args {
		switch a {
		case "--help", "-?", "help":
			fmt.Fprintln(os.Stderr, usage)
			return 0
		case "-h", "--hide":
			hide = true
		case "-s", "--separator":
			separator = true
		default:
			if rangeArg != "" {
				fmt.Fprintf(os.Stderr, "error: unexpected argument %q\n", a)
				return 1
			}
			rangeArg = a
		}
	}

	if rangeArg == "" {
		fmt.Fprintln(os.Stderr, usage)
		return 1
	}

	ranges, err := parseRanges(rangeArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	bounded, allBounded := maxEnd(ranges)
	canEarlyExit := allBounded && !hide

	lineNum := 0
	lastPrinted := 0
	firstOutput := true
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for input.Scan() {
		lineNum++
		matches := matchesAny(lineNum, ranges)
		shouldPrint := (matches && !hide) || (!matches && hide)

		if shouldPrint {
			if separator && !firstOutput && lastPrinted != lineNum-1 {
				_, err := fmt.Fprintln(writer, "---")
				if err != nil {
					return 0
				}
			}
			_, err := fmt.Fprintln(writer, input.Text())
			if err != nil {
				return 0
			}
			lastPrinted = lineNum
			firstOutput = false
		}

		if canEarlyExit && lineNum >= bounded {
			break
		}
	}

	return 0
}

func main() {
	signal.Ignore(syscall.SIGPIPE)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	os.Exit(run(os.Args[1:], scanner))
}
