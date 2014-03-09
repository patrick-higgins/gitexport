package gitexport

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// Replaces an \r\n sequences with \n in case Windows raw strings contain them.
func lfOnly(s string) string {
	return strings.Replace(s, "\r\n", "\n", -1)
}

func TestParse(t *testing.T) {
	r := strings.NewReader(lfOnly(`reset refs/heads/a
commit refs/heads/a
mark :1
author Some Guy <someguy@domain.us.uk> 1393367434 -0700
committer Some Guy <someguy@domain.us.uk> 1393367434 -0700
data 8
initial
M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt

reset refs/heads/master
from :2

tag c.Merge-to-a-1
from :2
tagger Some Guy <someguy@domain.us.uk> 1393367459 -0700
data 15
c.Merge-to-a-1

`))

	NewParser(r)
}

func TestCommit(t *testing.T) {
	r := strings.NewReader(lfOnly(`commit refs/heads/a
mark :4
author Some Guy <someguy@domain.us.uk> 1393367434 -0700
committer Some Guy <someguy@domain.us.uk> 1393367434 -0700
data 8
initial
from :1
merge :2
merge :3
M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt
`))

	p := NewParser(r)
	got, err := p.Commit()
	if err != nil {
		t.Fatal(err)
	}

	personTime := time.Unix(0, 0).Add(1393367434 * time.Second)
	wantPerson := Person{
		Name:  "Some Guy",
		Email: "<someguy@domain.us.uk>",
		When:  personTime,
	}

	want := &Commit{
		Ref:       "refs/heads/a",
		Mark:      4,
		Author:    &wantPerson,
		Committer: wantPerson,
		Message:   "initial\n",
		From:      ":1",
		Merge:     []string{":2", ":3"},
		Commands:  []FileCommand{"M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt\n"},
	}

	if !reflect.DeepEqual(got.Mark, want.Mark) {
		t.Errorf("got:\n\t%#v\nwant:\n\t%#v", got, want)
	}

}
