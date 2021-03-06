package errors

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"gocloud.dev/gcerrors"
)

func TestNewf(t *testing.T) {
	e := Newf(gcerrors.Internal, nil, "a %d b", 3)
	got := e.Error()
	want := "a 3 b (code=Internal)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatting(t *testing.T) {
	for i, test := range []struct {
		err  *Error
		verb string
		want []string // regexps, one per line
	}{
		{
			New(gcerrors.NotFound, nil, 1, "message"),
			"%v",
			[]string{`^message \(code=NotFound\)$`},
		},
		{
			New(gcerrors.NotFound, nil, 1, "message"),
			"%+v",
			[]string{
				`^message \(code=NotFound\):$`,
				`\s+github.com/felipeweb/devctl/internal/errors.TestFormatting$`,
				`\s+.*/internal/errors/errors_test.go:\d+$`,
			},
		},
		{
			New(gcerrors.AlreadyExists, errors.New("wrapped"), 1, "message"),
			"%v",
			[]string{`^message \(code=AlreadyExists\): wrapped$`},
		},
		{
			New(gcerrors.AlreadyExists, errors.New("wrapped"), 1, "message"),
			"%+v",
			[]string{
				`^message \(code=AlreadyExists\):`,
				`^\s+github.com/felipeweb/devctl/internal/errors.TestFormatting$`,
				`^\s+.*/internal/errors/errors_test.go:\d+$`,
				`^\s+- wrapped$`,
			},
		},
		{
			New(gcerrors.AlreadyExists, errors.New("wrapped"), 1, ""),
			"%v",
			[]string{`^code=AlreadyExists: wrapped`},
		},
		{
			New(gcerrors.AlreadyExists, errors.New("wrapped"), 1, ""),
			"%+v",
			[]string{
				`^code=AlreadyExists:`,
				`^\s+github.com/felipeweb/devctl/internal/errors.TestFormatting$`,
				`^\s+.*/internal/errors/errors_test.go:\d+$`,
				`^\s+- wrapped$`,
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotString := fmt.Sprintf(test.verb, test.err)
			gotLines := strings.Split(gotString, "\n")
			if got, want := len(gotLines), len(test.want); got != want {
				t.Fatalf("got %d lines, want %d. got:\n%s", got, want, gotString)
			}
			for j, gl := range gotLines {
				matched, err := regexp.MatchString(test.want[j], gl)
				if err != nil {
					t.Fatal(err)
				}
				if !matched {
					t.Fatalf("line #%d: got %q, which doesn't match %q", j, gl, test.want[j])
				}
			}
		})
	}
}

func TestError(t *testing.T) {
	// Check that err.Error() == fmt.Sprintf("%s", err)
	for _, err := range []*Error{
		New(gcerrors.NotFound, nil, 1, "message"),
		New(gcerrors.AlreadyExists, errors.New("wrapped"), 1, "message"),
		New(gcerrors.AlreadyExists, errors.New("wrapped"), 1, ""),
	} {
		got := err.Error()
		want := fmt.Sprint(err)
		if got != want {
			t.Errorf("%v: got %q, want %q", err, got, want)
		}
	}
}
