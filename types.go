package gitexport

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

type Person struct {
	Name  string    // empty means name is not present
	Email string    // non-empty
	When  time.Time // non-empty
}

func (p Person) Marshal() []byte {
	b := new(bytes.Buffer)
	if p.Name != "" {
		io.WriteString(b, p.Name)
		b.Write([]byte{' '})
	}
	io.WriteString(b, fmt.Sprintf("%s %d %s\n", p.Email, p.When.Unix(), p.When.Format("-0700")))
	return b.Bytes()
}

type Data struct{}

// FileCommand is one of the following in unparsed string form, including the final LF:
// filemodify | filedelete | filecopy | filerename | filedeleteall | notemodify
type FileCommand string

func (f FileCommand) String() string {
	return string(f)
}

type Commit struct {
	Ref       string        // non-empty
	Mark      int           // zero if mark is not present
	Author    *Person       // nil if author not present
	Committer Person        // non-empty
	Message   string        // non-empty
	From      string        // empty if from is not present
	Merge     []string      // maybe empty
	Commands  []FileCommand // maybe empty
}

func (c *Commit) Write(w io.Writer) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			if te, ok := e.(error); ok {
				err = te
			}
		}
	}()

	write := func(s string) {
		wn, e := io.WriteString(w, s)
		n += wn
		if e != nil {
			panic(e)
		}
	}

	write("commit " + c.Ref + "\n")
	if c.Mark != 0 {
		write(fmt.Sprintf("mark :%d\n", c.Mark))
	}
	if c.Author != nil {
		write("author ")
		write(string(c.Author.Marshal()))
	}
	write("committer ")
	write(string(c.Committer.Marshal()))

	msg := []byte(c.Message)
	write(fmt.Sprintf("data %d\n", len(msg)))

	var wn int
	wn, err = w.Write(msg)
	n += wn
	if err != nil {
		return
	}

	if c.From != "" {
		write("from " + c.From + "\n")
	}

	for _, m := range c.Merge {
		write("merge " + m + "\n")
	}

	for _, cmd := range c.Commands {
		write(string(cmd))
	}

	write("\n")

	return
}

type Tag struct{}
type Reset struct{}
type Blob struct{}
type Checkpoint struct{}
type Progress struct{}
type Done struct{}
type CatBlob struct{}
type Ls struct{}
type Feature struct{}
type Option struct{}
