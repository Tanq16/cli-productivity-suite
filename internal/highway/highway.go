package highway

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Job interface {
	ID() string
	Type() string
	Run(ctx context.Context, progress chan<- Progress) error
	Marshal() ([]byte, error)
}

type Highway struct {
	workers  int
	jobs     []Job
	progress chan Progress
}

func New(workers int) *Highway {
	return &Highway{
		workers:  workers,
		progress: make(chan Progress, 100),
	}
}

func (h *Highway) Submit(jobs ...Job) {
	h.jobs = append(h.jobs, jobs...)
}

func (h *Highway) Run(ctx context.Context) error {
	defer close(h.progress)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(h.workers)

	for _, j := range h.jobs {
		g.Go(func() error {
			if ctx.Err() != nil {
				h.progress <- Progress{
					JobID:   j.ID(),
					Message: "cancelled",
					Done:    true,
					Error:   ctx.Err(),
				}
				return nil
			}
			j.Run(ctx, h.progress)
			return nil
		})
	}

	return g.Wait()
}

func (h *Highway) Progress() <-chan Progress {
	return h.progress
}

func (h *Highway) JobCount() int {
	return len(h.jobs)
}
