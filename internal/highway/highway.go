package highway

import (
	"context"
	"encoding/json"
	"sync"
)

// JobUnmarshaler reconstructs a Job from persisted JSON bytes.
type JobUnmarshaler func([]byte) (Job, error)

// Job is a unit of work submitted to a Highway.
type Job interface {
	ID() string
	Type() string
	Run(ctx context.Context, progress chan<- Progress) error
	Marshal() ([]byte, error)
}

// Highway manages bounded-concurrency execution of jobs and exposes a progress channel.
type Highway struct {
	workers      int
	jobs         []Job
	progress     chan Progress
	mu           sync.Mutex
	completed    map[string]bool
	statePath    string
	unmarshalers map[string]JobUnmarshaler
}

// New creates a Highway with the given concurrency limit.
func New(workers int) *Highway {
	return &Highway{
		workers:      workers,
		progress:     make(chan Progress, 100),
		completed:    make(map[string]bool),
		unmarshalers: make(map[string]JobUnmarshaler),
	}
}

// RegisterType registers an unmarshaler for a job type, used by LoadState for resume.
func (h *Highway) RegisterType(jobType string, unmarshal JobUnmarshaler) {
	h.unmarshalers[jobType] = unmarshal
}

// Submit adds jobs to the highway.
func (h *Highway) Submit(jobs ...Job) {
	h.jobs = append(h.jobs, jobs...)
}

// Run executes all submitted jobs with bounded concurrency.
// It closes the progress channel when all jobs are done.
// It blocks until completion.
func (h *Highway) Run(ctx context.Context) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, h.workers)

	for _, j := range h.jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				h.progress <- Progress{
					JobID:   j.ID(),
					Message: "cancelled",
					Done:    true,
					Error:   ctx.Err(),
				}
				return
			case sem <- struct{}{}:
			}
			defer func() { <-sem }()

			// Jobs are responsible for sending their own Done progress.
			j.Run(ctx, h.progress)

			h.mu.Lock()
			h.completed[j.ID()] = true
			h.mu.Unlock()
		}()
	}

	wg.Wait()
	h.deleteState()
	close(h.progress)
}

// Progress returns the read-only progress channel.
func (h *Highway) Progress() <-chan Progress {
	return h.progress
}

// JobCount returns the number of submitted jobs.
func (h *Highway) JobCount() int {
	return len(h.jobs)
}

// SetStatePath configures the path for resume state persistence.
func (h *Highway) SetStatePath(path string) {
	h.statePath = path
}

// LoadState reads persisted state and reconstructs pending jobs for resume.
// Not implemented — this project does not require persistence.
func (h *Highway) LoadState() error {
	return nil
}

// resumeState is the JSON schema for persisted highway state.
type resumeState struct {
	Completed []string          `json:"completed"`
	Pending   []json.RawMessage `json:"pending"`
}

// saveState persists completed and pending job state to disk.
// Not implemented — this project does not require persistence.
func (h *Highway) saveState() error {
	return nil
}

// deleteState removes the persisted state file after successful completion.
// Not implemented — this project does not require persistence.
func (h *Highway) deleteState() {
}
