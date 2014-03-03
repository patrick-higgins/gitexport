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
author Patrick Higgins <patrick.allen.higgins@gmail.com> 1393367434 -0700
committer Patrick Higgins <patrick.allen.higgins@gmail.com> 1393367434 -0700
data 8
initial
M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt

reset refs/heads/master
from :2

tag c.Merge-to-a-1
from :2
tagger Patrick Higgins <patrick.allen.higgins@gmail.com> 1393367459 -0700
data 15
c.Merge-to-a-1

`))

	NewParser(r)
}

func TestCommit(t *testing.T) {
	r := strings.NewReader(lfOnly(`commit refs/heads/a
mark :1
author Patrick Higgins <patrick.allen.higgins@gmail.com> 1393367434 -0700
committer Patrick Higgins <patrick.allen.higgins@gmail.com> 1393367434 -0700
data 8
initial
M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt
`))

	p := NewParser(r)
	c, err := p.Commit()
	if err != nil {
		t.Fatal(err)
	}

	if c.Ref != "refs/heads/a" {
		t.Errorf("c.Ref=%#v, want=%#v", c.Ref, "refs/heads/a")
	}

	if *c.Mark != 1 {
		t.Errorf("c.Mark=%v, want=%v", c.Mark, 1)
	}

	personName := "Patrick Higgins"
	personTime := time.Unix(0, 0).Add(1393367434 * time.Second)
	wantPerson := &Person{
		Name:  &personName,
		Email: "<patrick.allen.higgins@gmail.com>",
		When:  personTime,
	}

	if !reflect.DeepEqual(c.Author, wantPerson) {
		t.Errorf("c.Author=%#v, want=%#v", c.Author, wantPerson)
	}

	if !reflect.DeepEqual(&c.Committer, wantPerson) {
		t.Errorf("c.Committer=%#v, want=%#v", c.Committer, wantPerson)
	}

	if c.Message != "initial\n" {
		t.Errorf("c.Message=%#v, want=%#v", c.Message, "initial\n")
	}

	if c.From != nil {
		t.Errorf("c.From=%#v, want nil", *c.From)
	}

	if c.Merge != nil {
		t.Errorf("c.Merge=%#v, want nil", c.Merge)
	}

	wantCommands := []FileCommand{"M 100644 e79c5e8f964493290a409888d5413a737e8e5dd5 test.txt\n"}
	if !reflect.DeepEqual(c.Commands, wantCommands) {
		t.Errorf("c.Commands=%#v, want=%#v", c.Commands, wantCommands)
	}
}
