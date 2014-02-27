package gitexport

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Stream struct {
	r          *bufio.Reader
	lineNumber int
	cmd        []string
}

func NewStream(r io.Reader) *Stream {
	return &Stream{r: bufio.NewReader(r)}
}

func (s *Stream) error(msg interface{}) error {
	return fmt.Errorf("line %d: %v", s.lineNumber, msg)
}

func (s *Stream) nextLine() Token {
	var cmd string
	var err error
	for {
		s.lineNumber++
		cmd, err = s.r.ReadString('\n')
		if len(cmd) > 1 {
			break
		}
		if err == io.EOF {
			return EOFTok
		}
		if err != nil {
			panic(err)
		}
	}
	s.cmd = strings.Split(cmd, " ")
	if len(s.cmd) == 0 {
		panic(s.error("invalid command: " + cmd))
	}
	tok, ok := tokenMap[s.cmd[0]]
	if !ok {
		panic(s.error("unknown token: " + s.cmd[0]))
	}
	return tok
}

func (s *Stream) parseCommand() interface{} {
	tok := s.nextLine()

	switch tok {
	case CommitTok:
		return s.commit()
	case TagTok:
		return s.tag()
	case ResetTok:
		return s.reset()
	case BlobTok:
		return s.blob()
	case CheckPointTok:
		return s.checkpoint()
	case ProgressTok:
		return s.progress()
	case DoneTok:
		return s.done()
	case LsTok:
		return s.ls()
	case FeatureTok:
		return s.feature()
	case OptionTok:
		return s.option()
	default:
		panic(s.error(fmt.Sprintf("invalid top-level token: %s", tok)))
	}
}

var errUnsupported = errors.New("Unsupported operation")

func (s *Stream) commit() interface{} {
	var c Commit
	c.Ref = s.cmd[1]
	tok := s.nextLine()
	if tok == MarkTok {
		c.Mark = new(int64)
		*c.Mark = s.mark()
	}
	return &c
}

func (s *Stream) mark() int64 {
	var m int64
	n, err := fmt.Sscanf(s.cmd[1], ":%d", &m)
	if err != nil {
		panic(err)
	}
	if n == 1 {
		return m
	}
	panic(errors.New("could not parse mark: " + s.cmd[1]))
}

func (s *Stream) tag() interface{} {
	panic(errUnsupported)
}

func (s *Stream) reset() interface{} {
	panic(errUnsupported)
}

func (s *Stream) blob() interface{} {
	panic(errUnsupported)
}

func (s *Stream) checkpoint() interface{} {
	panic(errUnsupported)
}

func (s *Stream) progress() interface{} {
	panic(errUnsupported)
}

func (s *Stream) done() interface{} {
	panic(errUnsupported)
}

func (s *Stream) ls() interface{} {
	panic(errUnsupported)
}

func (s *Stream) feature() interface{} {
	panic(errUnsupported)
}

func (s *Stream) option() interface{} {
	panic(errUnsupported)
}
