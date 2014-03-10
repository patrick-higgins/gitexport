package gitexport

import (
	"bytes"
	"testing"
	"time"
)

func TestCommitWrite(t *testing.T) {
	refTime, err := time.Parse(time.RFC1123Z, time.RFC1123Z)
	if err != nil {
		t.Fatalf("could not parse reference time: %v", err)
	}

	for i, tt := range []struct {
		c    Commit
		want string
	}{
		{
			Commit{
				Ref: "ref",
				Committer: Person{
					Email: "<some.guy@domain.us.uk>",
					When:  refTime,
				},
				Message: "message\n",
			},
			`commit ref
committer <some.guy@domain.us.uk> 1136239445 -0700
data 8
message

`,
		},
		{
			Commit{
				Ref:  "ref",
				Mark: 42,
				Author: &Person{
					Name:  "Some Guy",
					Email: "<some.guy@domain.us.uk>",
					When:  refTime,
				},
				Committer: Person{
					Name:  "Some Other Guy",
					Email: "<some.other.guy@domain.us.uk>",
					When:  refTime,
				},
				Message:  "message\n",
				From:     ":43",
				Merge:    []string{":44", ":45"},
				Commands: []FileCommand{"D foo/bar\n", "M 100755 :46 foo/bar\n"},
			},
			`commit ref
mark :42
author Some Guy <some.guy@domain.us.uk> 1136239445 -0700
committer Some Other Guy <some.other.guy@domain.us.uk> 1136239445 -0700
data 8
message
from :43
merge :44
merge :45
D foo/bar
M 100755 :46 foo/bar

`,
		},
	} {
		b := new(bytes.Buffer)
		tt.c.Write(b)
		got := b.String()
		want := lfOnly(tt.want)
		if got != want {
			t.Errorf("[%d] commit write:\ngot=%#v\n%v\nwant=%#v\n%v", i, got, got, want, want)
		}
	}

}
