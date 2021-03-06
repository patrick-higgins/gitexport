package lex

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	for i, tt := range []struct {
		line string
		err  error
		want Token
	}{
		// Lines must have an LF
		{"data", nil, InvalidTok},
		{"data ", nil, InvalidTok},
		{"data arg1", nil, InvalidTok},
		{"data arg1 arg2", nil, InvalidTok},

		// Arguments are allowed
		{"data\n", nil, DataTok},
		{"data arg1\n", nil, DataTok},
		{"data arg1 arg2\n", nil, DataTok},

		// Unknown commands are invalid
		{"foo\n", nil, InvalidTok},

		// Comments
		{"#", nil, CommentTok},
		{"#a", nil, CommentTok},
		{"#a\n", nil, CommentTok},
		{"# a", nil, CommentTok},
		{"# a\n", nil, CommentTok},

		// Errors are ignored if the line isn't empty
		{"\n", io.EOF, EmptyTok},
		{"data\n", io.EOF, DataTok},
		{"data\n", io.ErrClosedPipe, DataTok},

		// Errors are returned when line is empty
		{"", io.EOF, EOFTok},
		{"", io.ErrClosedPipe, ErrTok},
	} {
		got := classify(tt.line, tt.err)
		if got != tt.want {
			t.Errorf("[%d] classify(%#v, %#v)=%v, want=%v", i, tt.line, tt.err, got, tt.want)
		}
	}
}

func TestLex(t *testing.T) {
	r := strings.NewReader(strings.Replace(`reset refs/heads/a
commit refs/heads/a
mark :1
author Some Guy <someguy@gmail.com.uk> 1393367434 -0700
committer Some Guy <someguy@gmail.com.uk> 1393367434 -0700
M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt

reset refs/heads/master
from :2

# a comment
tag c.Merge-to-a-1
#another comment
from :2
tagger Some Guy <someguy@gmail.com.uk> 1393367459 -0700

`, "\r\n", "\n", -1))

	l := New(r)

	for i, tt := range []struct {
		tok    Token
		fields []string
	}{
		{ResetTok, []string{"reset", "refs/heads/a"}},
		{CommitTok, []string{"commit", "refs/heads/a"}},
		{MarkTok, []string{"mark", ":1"}},
		{AuthorTok, []string{"author", "Some", "Guy", "<someguy@gmail.com.uk>", "1393367434", "-0700"}},
		{CommitterTok, []string{"committer", "Some", "Guy", "<someguy@gmail.com.uk>", "1393367434", "-0700"}},
		{MTok, []string{"M", "100644", "e79c5e8f964493290a409888d5413a737e8e5dd5", "test.txt"}},
		{EmptyTok, []string{""}},
		{ResetTok, []string{"reset", "refs/heads/master"}},
		{FromTok, []string{"from", ":2"}},
		{EmptyTok, []string{""}},
		{CommentTok, []string{"#", "a", "comment"}},
		{TagTok, []string{"tag", "c.Merge-to-a-1"}},
		{CommentTok, []string{"#another", "comment"}},
		{FromTok, []string{"from", ":2"}},
		{TaggerTok, []string{"tagger", "Some", "Guy", "<someguy@gmail.com.uk>", "1393367459", "-0700"}},
		{EmptyTok, []string{""}},
		{EOFTok, nil},
	} {
		tok := l.Token()
		fields := l.Fields()
		if tok != tt.tok {
			t.Errorf("[%d]: token got=%v, want=%v", i, tok, tt.tok)
		}
		if !reflect.DeepEqual(fields, tt.fields) {
			t.Errorf("[%d]: fields got=%#v, want=%#v", i, fields, tt.fields)
		}
		l.Consume()
	}
}

func TestData(t *testing.T) {
	for i, tt := range []struct {
		data       string
		want       []byte
		lineNumber int
		err        bool
	}{
		// Simple case
		{"data 3\nfoo", []byte("foo"), 2, false},

		// Error on short read
		{"data 4\nfoo", nil, 1, true},

		// Increment line number for each \n
		{"data 12\nfoo\nbar\nbaz\n", []byte("foo\nbar\nbaz\n"), 5, false},

		// Delimited format
		{"data <<EOF\nfoo\nEOF\n", []byte("foo\n"), 4, false},
	} {
		r := strings.NewReader(tt.data)
		l := New(r)

		got, err := l.ConsumeData()
		if tt.err && err == nil {
			t.Errorf("[%d] expected error but didn't get one", i)
		}
		if !tt.err && err != nil {
			t.Errorf("[%d] error: %v", i, err)
		}

		if !bytes.Equal(got, tt.want) {
			t.Errorf("[%d] got=%v, want=%v", i, got, tt.want)
		}

		if l.LineNumber() != tt.lineNumber {
			t.Errorf("[%d] lineNumber=%v, want=%v", i, l.LineNumber(), tt.lineNumber)
		}
	}
}
