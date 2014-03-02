package lex

func init() {
	tokenMap = make(map[string]Token)
	for i, s := range tokenStrings {
		tokenMap[s] = Token(i)
	}
}

type Token int

// Tokens
const (
	InvalidTok Token = iota
	ErrTok
	EOFTok
	EmptyTok
	CommentTok
	CommitTok
	TagTok
	ResetTok
	BlobTok
	CheckPointTok
	ProgressTok
	DoneTok
	CatBlobTok
	LsTok
	FeatureTok
	OptionTok
	MarkTok
	AuthorTok
	CommitterTok
	DataTok
	FromTok
	MergeTok
	CTok
	DTok
	MTok
	NTok
	RTok
	DeleteAllTok
	TaggerTok
	NTokens
)

var tokenStrings = []string{
	InvalidTok:    "<INVALID>",
	ErrTok:        "<ERROR>",
	EOFTok:        "<EOF>",
	EmptyTok:      "<EMPTY>",
	CommentTok:    "#comment",
	CommitTok:     "commit",
	TagTok:        "tag",
	ResetTok:      "reset",
	BlobTok:       "blob",
	CheckPointTok: "checkpoint",
	ProgressTok:   "progress",
	DoneTok:       "done",
	CatBlobTok:    "cat-blob",
	LsTok:         "ls",
	FeatureTok:    "feature",
	OptionTok:     "option",
	MarkTok:       "mark",
	AuthorTok:     "author",
	CommitterTok:  "committer",
	DataTok:       "data",
	FromTok:       "from",
	MergeTok:      "merge",
	CTok:          "C",
	DTok:          "D",
	MTok:          "M",
	NTok:          "N",
	RTok:          "R",
	DeleteAllTok:  "deleteall",
	TaggerTok:     "tagger",
}

var tokenMap map[string]Token

func (t Token) String() string {
	if t < 0 || t >= NTokens {
		return tokenStrings[InvalidTok]
	}
	return tokenStrings[t]
}
