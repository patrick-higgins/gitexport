package gitexport

import (
	"errors"
	"fmt"
	"github.com/patrick-higgins/gitexport/lex"
	"io"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

type Parser struct {
	l *lex.Lexer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{l: lex.New(r)}
}

func NewLexerParser(l *lex.Lexer) *Parser {
	return &Parser{l: l}
}

func (p *Parser) error(msg interface{}) error {
	return fmt.Errorf("line %d: %v", p.l.LineNumber(), msg)
}

func (p *Parser) parseCommand() interface{} {
	tok := p.l.Token()

	switch tok {
	case lex.CommitTok:
		return p.commit()
	case lex.TagTok:
		return p.tag()
	case lex.ResetTok:
		return p.reset()
	case lex.BlobTok:
		return p.blob()
	case lex.CheckPointTok:
		return p.checkpoint()
	case lex.ProgressTok:
		return p.progress()
	case lex.DoneTok:
		return p.done()
	case lex.LsTok:
		return p.ls()
	case lex.FeatureTok:
		return p.feature()
	case lex.OptionTok:
		return p.option()
	default:
		panic(p.error(fmt.Sprintf("invalid top-level token: %s", tok)))
	}
}

var errUnsupported = errors.New("Unsupported operation")
var errWrongType = errors.New("Internal error: unexpected type")
var errWrongTok = errors.New("Not looking at requested token type")

func (p *Parser) Commit() (c *Commit, err error) {
	var ok bool
	defer func() {
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	if p.l.Token() != lex.CommitTok {
		return nil, errWrongTok
	}

	c, ok = p.commit().(*Commit)
	if !ok {
		return nil, errWrongType
	}
	return
}

func (p *Parser) commit() interface{} {
	var c Commit
	c.Ref = p.l.Field(1)

	p.l.Consume()
	tok := p.l.Token()
	if tok == lex.MarkTok {
		c.Mark = new(int64)
		*c.Mark = p.mark()
		p.l.Consume()
		tok = p.l.Token()
	}
	if tok == lex.AuthorTok {
		c.Author = p.person()
		p.l.Consume()
		tok = p.l.Token()
	}

	if tok != lex.CommitterTok {
		panic(errors.New("missing committer"))
	}

	c.Committer = *p.person()
	p.l.Consume()
	tok = p.l.Token()

	if tok != lex.DataTok {
		panic(errors.New("missing commit message"))
	}

	data, err := p.l.ConsumeData()
	if err != nil {
		panic(err)
	}

	// Messages are UTF-8, thankfully
	c.Message = string(data)
	tok = p.l.Token()

	if tok == lex.FromTok {
		from := p.l.Field(1)
		c.From = &from
		p.l.Consume()
		tok = p.l.Token()
	}

	for tok == lex.MergeTok {
		merge := p.l.Field(1)
		c.Merge = append(c.Merge, merge)
		p.l.Consume()
		tok = p.l.Token()
	}

FileCommands:
	for {
		switch p.l.Token() {
		case lex.MTok, lex.DTok, lex.CTok, lex.RTok, lex.DeleteAllTok, lex.NTok:
			c.Commands = append(c.Commands, FileCommand(p.l.Line()))
			p.l.Consume()
		default:
			break FileCommands
		}
	}

	return &c
}

func (p *Parser) mark() int64 {
	var m int64
	n, err := fmt.Sscanf(p.l.Field(1), ":%d", &m)
	if err != nil {
		panic(err)
	}
	if n == 1 {
		return m
	}
	panic(errors.New("could not parse mark: " + p.l.Field(1)))
}

func (p *Parser) person() *Person {
	var person Person
	emailIdx := 1
	fields := p.l.Fields()
	for fields[emailIdx][0] != '<' {
		emailIdx++
	}
	if emailIdx > 1 {
		s := strings.Join(fields[1:emailIdx], " ")
		person.Name = &s
	}
	person.Email = fields[emailIdx]
	person.When = p.when(fields[emailIdx+1:])
	return &person
}

func (p *Parser) when(val []string) time.Time {
	// now format
	if val[0] == "now" {
		return time.Now()
	}

	// raw format
	var sec int64
	_, err := fmt.Sscanf(val[0], "%d", &sec)
	if err == nil {
		return time.Unix(sec, 0)
	}

	// rfc2822 format
	const rfc2822 = "Mon Jan 2 15:04:05 2006 -0700"
	t, err := time.Parse(rfc2822, strings.Join(val, " "))
	if err != nil {
		panic(err)
	}
	return t
}

func (p *Parser) tag() interface{} {
	panic(errUnsupported)
}

func (p *Parser) reset() interface{} {
	panic(errUnsupported)
}

func (p *Parser) blob() interface{} {
	panic(errUnsupported)
}

func (p *Parser) checkpoint() interface{} {
	panic(errUnsupported)
}

func (p *Parser) progress() interface{} {
	panic(errUnsupported)
}

func (p *Parser) done() interface{} {
	panic(errUnsupported)
}

func (p *Parser) ls() interface{} {
	panic(errUnsupported)
}

func (p *Parser) feature() interface{} {
	panic(errUnsupported)
}

func (p *Parser) option() interface{} {
	panic(errUnsupported)
}
