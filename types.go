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
type Commitish struct{}
type FileTweaker interface{}

type Commit struct {
	Ref          string
	Mark         *int64
	Author       *Person
	Committer    Person
	Message      string
	From         *Commitish
	Merge        *Commitish
	FileCommands []FileTweaker
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
