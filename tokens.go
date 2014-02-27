package gitexport

type Token int

// Tokens
const (
	InvalidTok Token = iota
	EOFTok
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
	InlineTok
	DeleteAllTok
	TaggerTok
	DelimTok
	NTokens
)

var tokenStrings = []string{
	InvalidTok:    "<INVALID>",
	EOFTok:        "<EOF>",
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
	InlineTok:     "inline",
	DeleteAllTok:  "deleteall",
	TaggerTok:     "tagger",
	DelimTok:      "<<",
}

var tokenMap map[string]Token

func initTokenMap() {
	tokenMap = make(map[string]Token)
	for i, s := range tokenStrings {
		tokenMap[s] = Token(i)
	}
}

func (t Token) String() string {
	if t < 0 || t >= NTokens {
		return tokenStrings[InvalidTok]
	}
	return tokenStrings[t]
}
