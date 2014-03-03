package gitexport

import (
	"time"
)

type Person struct {
	Name  *string
	Email string
	When  time.Time
}

type Data struct{}

// FileCommand is one of the following in unparsed string form, including the final LF:
// filemodify | filedelete | filecopy | filerename | filedeleteall | notemodify
type FileCommand string

func (f FileCommand) String() string {
	return string(f)
}

type Commit struct {
	Ref       string
	Mark      *int64
	Author    *Person
	Committer Person
	Message   string
	From      *string
	Merge     []string
	Commands  []FileCommand
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
