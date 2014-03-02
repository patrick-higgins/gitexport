package gitexport

import (
	"errors"
	"fmt"
	"github.com/patrick-higgins/gitexport/lex"
	"io"
	"strings"
	"time"
)

type Parser struct {
	l *lex.Lexer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{l: lex.New(r)}
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
		c.From = p.commitish()
		p.l.Consume()
		tok = p.l.Token()
	}

	if tok == lex.MergeTok {
		c.Merge = p.commitish()
		p.l.Consume()
		tok = p.l.Token()
	}

	return &c
}

func (p *Parser) commitish() *Commitish {
	return nil
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
	if val[0] == "now" {
		return time.Now()
	}

	// raw format
	var sec int64
	n, err := fmt.Sscanf(val[0], "%d", &sec)
	if err != nil {
		panic(err)
	}
	if n == 1 {
		return time.Unix(sec, 0)
	}

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
