package display

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/rs/zerolog/log"

	"github.com/tanq16/cli-productivity-suite/internal/highway"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var (
	phaseStyle   = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(12)).Bold(true)
	runningStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
	doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(9))
	skipStyle    = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(7))
)

type jobState struct {
	id        string
	status    string // "running", "done", "error", "skipped"
	message   string
	subStatus string
	err       error
}

type Display struct {
	mu   sync.Mutex
	jobs map[string]*jobState
}

func New() *Display {
	return &Display{
		jobs: make(map[string]*jobState),
	}
}

func (d *Display) StartPhase(name string, progress <-chan highway.Progress) {
	d.jobs = make(map[string]*jobState)

	if utils.GlobalDebugFlag {
		d.startPhaseDebug(name, progress)
	} else if utils.GlobalForAIFlag {
		d.startPhaseAI(name, progress)
	} else {
		d.startPhaseNormal(name, progress)
	}
}

func (d *Display) startPhaseDebug(name string, progress <-chan highway.Progress) {
	log.Info().Str("package", "display").Str("phase", name).Msg("starting phase")
	for p := range progress {
		if p.Done {
			if p.Error != nil {
				log.Error().Str("package", "display").Str("job", p.JobID).Err(p.Error).Msg(p.Message)
			} else {
				log.Info().Str("package", "display").Str("job", p.JobID).Msg(p.Message)
			}
		} else {
			log.Debug().Str("package", "display").Str("job", p.JobID).Msg(p.SubStatus)
		}
	}
	log.Info().Str("package", "display").Str("phase", name).Msg("phase complete")
}

func (d *Display) startPhaseAI(name string, progress <-chan highway.Progress) {
	fmt.Printf("[PHASE] %s\n", name)
	for p := range progress {
		if p.Done {
			if p.Error != nil {
				fmt.Printf("[ERROR] %s: %s\n", p.JobID, p.Error)
			} else {
				fmt.Printf("[OK] %s: %s\n", p.JobID, p.Message)
			}
			continue
		}
		if p.Type == highway.ProgressTypeProgress && p.Total > 0 {
			pct := int(float64(p.Current) / float64(p.Total) * 100)
			extra := ""
			if p.Extra != "" {
				extra = " (" + p.Extra + ")"
			}
			fmt.Printf("[INFO] %s: %s %d%%%s\n", p.JobID, p.Message, pct, extra)
		} else if p.Type == highway.ProgressTypeSubStatus {
			sub := p.SubStatus
			if sub == "" {
				sub = p.Message
			}
			fmt.Printf("[INFO] %s: %s\n", p.JobID, sub)
		}
	}
}

func (d *Display) startPhaseNormal(name string, progress <-chan highway.Progress) {
	lipgloss.Println(phaseStyle.Render("▸ " + name))

	done := make(chan struct{})
	go func() {
		for p := range progress {
			d.mu.Lock()
			js, ok := d.jobs[p.JobID]
			if !ok {
				js = &jobState{id: p.JobID, status: "running"}
				d.jobs[p.JobID] = js
			}
			if p.Done {
				if p.Error != nil {
					js.status = "error"
					js.err = p.Error
					js.message = p.Message
				} else if strings.Contains(p.Message, "skipped") || strings.Contains(p.Message, "already at") {
					js.status = "skipped"
					js.message = p.Message
				} else {
					js.status = "done"
					js.message = p.Message
				}
			} else {
				js.subStatus = p.SubStatus
			}
			d.mu.Unlock()
		}
		close(done)
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var lastRenderedLines int
	for {
		select {
		case <-done:
			d.clearLines(lastRenderedLines)
			d.renderFinal()
			return
		case <-ticker.C:
			d.mu.Lock()
			d.clearLines(lastRenderedLines)
			lastRenderedLines = d.renderLive()
			d.mu.Unlock()
		}
	}
}

func (d *Display) clearLines(n int) {
	for range n {
		fmt.Print("\033[A\033[2K")
	}
}

func (d *Display) renderLive() int {
	ids := d.sortedIDs()
	lines := 0
	for _, id := range ids {
		js := d.jobs[id]
		var indicator string
		switch js.status {
		case "running":
			indicator = runningStyle.Render("↻")
		case "done":
			indicator = doneStyle.Render("✓")
		case "error":
			indicator = errStyle.Render("✗")
		case "skipped":
			indicator = skipStyle.Render("·")
		}
		detail := js.subStatus
		if js.status != "running" {
			detail = js.message
		}
		lipgloss.Printf("  %s %s %s\n", indicator, id, skipStyle.Render(detail))
		lines++
	}
	return lines
}

func (d *Display) renderFinal() {
	ids := d.sortedIDs()
	for _, id := range ids {
		js := d.jobs[id]
		if utils.GlobalForAIFlag {
			switch js.status {
			case "done":
				fmt.Printf("[OK] %s: %s\n", id, js.message)
			case "error":
				fmt.Printf("[ERROR] %s: %s\n", id, js.err.Error())
			case "skipped":
				fmt.Printf("[OK] %s: %s\n", id, js.message)
			}
			continue
		}
		switch js.status {
		case "done":
			lipgloss.Printf("  %s %s %s\n", doneStyle.Render("✓"), id, skipStyle.Render(js.message))
		case "error":
			lipgloss.Printf("  %s %s %s\n", errStyle.Render("✗"), id, errStyle.Render(js.err.Error()))
		case "skipped":
			lipgloss.Printf("  %s %s %s\n", skipStyle.Render("·"), id, skipStyle.Render(js.message))
		}
	}
}

func (d *Display) sortedIDs() []string {
	ids := make([]string, 0, len(d.jobs))
	for id := range d.jobs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
