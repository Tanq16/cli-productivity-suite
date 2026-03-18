package highway

type ProgressType int

const (
	ProgressTypeSubStatus ProgressType = iota
	ProgressTypeProgress
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
}
