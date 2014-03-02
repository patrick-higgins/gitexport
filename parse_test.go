package gitexport

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(strings.Replace(`reset refs/heads/a
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

`, "\r\n", "\n", -1))

	NewParser(r)
}
