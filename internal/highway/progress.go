package highway

// ProgressType distinguishes between sub-status messages and numeric progress.
type ProgressType int

const (
	ProgressTypeSubStatus ProgressType = iota
	ProgressTypeProgress
)

// Progress represents a single status update from a running job.
type Progress struct {
	JobID     string
	Message   string
	SubStatus string
	Type      ProgressType
	Current   int
	Total     int
	Done      bool
	Error     error
}
