package highway

type ProgressType int

const (
	ProgressTypeProgress  ProgressType = iota
	ProgressTypeSubStatus
)

type Progress struct {
	JobID     string
	Message   string
	SubStatus string
	Extra     string
	Type      ProgressType
	Current   int64
	Total     int64
	Done      bool
	Error     error
	ErrMsg    string
}
