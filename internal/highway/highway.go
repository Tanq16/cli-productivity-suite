package highway

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Job interface {
	ID() string
	Run(ctx context.Context, progress chan<- Progress) error
}

type Highway struct {
	workers int

	mu        sync.Mutex
	pending   []Job
	completed map[string]bool
	progress  chan Progress
}

func New(workers int) *Highway {
	if workers < 1 {
		workers = 1
	}
	return &Highway{
		workers:   workers,
		completed: make(map[string]bool),
		progress:  make(chan Progress, 100),
	}
}

func (h *Highway) Submit(jobs ...Job) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.pending = append(h.pending, jobs...)
}

func (h *Highway) Progress() <-chan Progress {
	return h.progress
}

func (h *Highway) Run(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(h.workers)

	h.mu.Lock()
	jobs := h.pending
	h.mu.Unlock()

	for _, j := range jobs {
		if h.isCompleted(j.ID()) {
			continue
		}
		g.Go(func() error {
			h.executeJob(gCtx, j)
			return nil
		})
	}

	err := g.Wait()
	close(h.progress)
	return err
}

func (h *Highway) executeJob(ctx context.Context, job Job) {
	if ctx.Err() != nil {
		h.progress <- Progress{
			JobID:   job.ID(),
			Message: "cancelled",
			Done:    true,
			Error:   ctx.Err(),
		}
		return
	}

	err := job.Run(ctx, h.progress)
	if err != nil {
		h.progress <- Progress{
			JobID:  job.ID(),
			Done:   true,
			Error:  err,
			ErrMsg: err.Error(),
		}
		h.markFailed(job.ID())
	} else {
		h.markCompleted(job.ID())
	}
}

func (h *Highway) isCompleted(id string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.completed[id]
}

func (h *Highway) markCompleted(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.completed[id] = true
}

func (h *Highway) markFailed(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.completed[id] = true
}

func (h *Highway) JobCount() int {
	return len(h.pending)
}
