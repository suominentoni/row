package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureRun(args []string, input string) (string, int) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	scanner := bufio.NewScanner(strings.NewReader(input))
	code := run(args, scanner)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	return buf.String(), code
}

func lines(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "%d\n", i)
	}
	return b.String()
}

func TestRunSingleLine(t *testing.T) {
	out, code := captureRun([]string{"3"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if out != "3\n" {
		t.Errorf("output = %q, want %q", out, "3\n")
	}
}

func TestRunRange(t *testing.T) {
	out, code := captureRun([]string{"5-10"}, lines(100))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "5\n6\n7\n8\n9\n10\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunOpenStart(t *testing.T) {
	out, code := captureRun([]string{"...3"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "1\n2\n3\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunOpenEnd(t *testing.T) {
	out, code := captureRun([]string{"95..."}, lines(100))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "95\n96\n97\n98\n99\n100\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunEmptyInput(t *testing.T) {
	out, code := captureRun([]string{"1-5"}, "")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if out != "" {
		t.Errorf("output = %q, want empty", out)
	}
}

func TestRunSingleLineInput(t *testing.T) {
	out, code := captureRun([]string{"1"}, "only\n")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if out != "only\n" {
		t.Errorf("output = %q, want %q", out, "only\n")
	}
}

func TestRunHide(t *testing.T) {
	out, code := captureRun([]string{"-h", "3-5"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "1\n2\n6\n7\n8\n9\n10\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunHideMultiRange(t *testing.T) {
	out, code := captureRun([]string{"--hide", "1,10"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "2\n3\n4\n5\n6\n7\n8\n9\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunSeparator(t *testing.T) {
	out, code := captureRun([]string{"-s", "3-5,10-12"}, lines(20))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "3\n4\n5\n---\n10\n11\n12\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunSeparatorContiguous(t *testing.T) {
	out, code := captureRun([]string{"-s", "3-5,6-8"}, lines(20))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "3\n4\n5\n6\n7\n8\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunCombinedFlags(t *testing.T) {
	out, code := captureRun([]string{"-hs", "3-5"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "1\n2\n---\n6\n7\n8\n9\n10\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunNumber(t *testing.T) {
	out, code := captureRun([]string{"-n", "3-5"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "3: 3\n4: 4\n5: 5\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunNumberMultiRange(t *testing.T) {
	out, code := captureRun([]string{"-n", "2-3,8-9"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "2: 2\n3: 3\n8: 8\n9: 9\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunNumberWithHide(t *testing.T) {
	out, code := captureRun([]string{"-nh", "3-5"}, lines(6))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "1: 1\n2: 2\n6: 6\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunNumberWithHideMultiRange(t *testing.T) {
	out, code := captureRun([]string{"-nh", "2-3,8-9"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "1: 1\n4: 4\n5: 5\n6: 6\n7: 7\n10: 10\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestRunNoArgs(t *testing.T) {
	_, code := captureRun([]string{}, "")
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func TestRunHelp(t *testing.T) {
	_, code := captureRun([]string{"--help"}, "")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
}

func TestRunInvalidRange(t *testing.T) {
	_, code := captureRun([]string{"abc"}, lines(10))
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func TestRunMultiRange(t *testing.T) {
	out, code := captureRun([]string{"1,5,9-10"}, lines(10))
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "1\n5\n9\n10\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}
