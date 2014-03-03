package lex

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Lexer struct {
	r      *bufio.Reader
	token  Token
	line   string
	lineno int
	fields []string
	err    error
}

func splitFields(line string) []string {
	if len(line) == 0 {
		return nil
	}
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	return strings.Split(line, " ")
}

func New(r io.Reader) *Lexer {
	bufr := bufio.NewReader(r)
	line, err := bufr.ReadString('\n')

	return &Lexer{
		r:      bufr,
		token:  classify(line, err),
		lineno: 1,
		line:   line,
		fields: splitFields(line),
		err:    err,
	}
}

func (l *Lexer) Token() Token {
	return l.token
}

func (l *Lexer) Line() string {
	return l.line
}

func (l *Lexer) LineNumber() int {
	return l.lineno
}

func (l *Lexer) Field(i int) string {
	return l.fields[i]
}

func (l *Lexer) Fields() []string {
	return l.fields
}

func (l *Lexer) Error() error {
	return l.err
}

func (l *Lexer) Consume() {
	line, err := l.r.ReadString('\n')
	l.token = classify(line, err)
	l.line = line
	l.lineno++
	l.fields = splitFields(line)
	l.err = err
}

// Consumes and returns a data element.
//
// Returns nil and an error if not currently looking at a data token or
// an I/O error occurs.
func (l *Lexer) ConsumeData() ([]byte, error) {
	if l.Token() != DataTok {
		return nil, errors.New("not on a data token")
	}

	var count int64
	_, err := fmt.Sscanf(l.line, "data %d\n", &count)
	if err != nil {
		var delim string
		_, err = fmt.Sscanf(l.line, "data <<%s\n", &delim)
		if err == nil {
			data, err := l.consumeDelimitedData(delim)
			l.Consume()
			return data, err
		}
		return nil, fmt.Errorf("invalid data byte count: %v", err)
	}

	data := make([]byte, count)

	_, err = io.ReadFull(l.r, data)
	if err != nil {
		return nil, err
		// return fmt.Errorf("could not read all of data: %v", err)
	}

	l.lineno += bytes.Count(data, []byte{'\n'})

	l.Consume()

	return data, nil
}

func (l *Lexer) consumeDelimitedData(delimStr string) ([]byte, error) {
	delim := make([]byte, len(delimStr)+1)
	copy(delim, delimStr)
	delim[len(delimStr)] = '\n'

	var data []byte
	for {
		line, err := l.r.ReadBytes('\n')
		if line[len(line)-1] == '\n' {
			l.lineno++
		}
		if bytes.Equal(line, delim) {
			return data, err
		}
		data = append(data, line...)
		if err != nil {
			return data, err
		}
	}
}

func classify(line string, err error) Token {
	if err != nil && len(line) == 0 {
		if err == io.EOF {
			return EOFTok
		}
		return ErrTok
	}

	if len(line) == 0 || len(line) == 1 && line[0] == '\n' {
		return EmptyTok
	}

	if line[0] == '#' {
		return CommentTok
	}

	i := strings.IndexAny(line, " \n")
	if i == -1 || line[len(line)-1] != '\n' {
		return InvalidTok
	}

	cmd := line[:i]
	tok, ok := tokenMap[cmd]
	if !ok {
		return InvalidTok
	}
	return tok
}
