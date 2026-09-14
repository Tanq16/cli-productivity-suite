package utils

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/rs/zerolog/log"
)

type Unit string

const UnitBytes Unit = ""

const (
	barFloor          = 8
	barCeiling        = 30
	rateWindowSpan    = 800 * time.Millisecond
	rateFloor         = 200 * time.Millisecond
	etaCeiling        = 100 * time.Hour
	indeterminateStep = 80 * time.Millisecond
)

var (
	meterMutedStyle  = lipgloss.NewStyle().Foreground(ColorMuted)
	meterChromeStyle = lipgloss.NewStyle().Foreground(ColorChrome)
)

var watchForInterrupt = sync.OnceFunc(func() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Print("\033[?25h")
		os.Exit(130)
	}()
})

type rateSample struct {
	at  time.Time
	val int64
}

type rateWindow struct {
	samples []rateSample
}

func (r *rateWindow) add(at time.Time, val int64) {
	r.samples = append(r.samples, rateSample{at: at, val: val})
	cutoff := at.Add(-rateWindowSpan)
	drop := 0
	for drop < len(r.samples)-2 && !r.samples[drop+1].at.After(cutoff) {
		drop++
	}
	r.samples = r.samples[drop:]
}

func (r *rateWindow) current() float64 {
	if len(r.samples) < 2 {
		return 0
	}
	first, last := r.samples[0], r.samples[len(r.samples)-1]
	if last.at.Sub(first.at) < rateFloor {
		return 0
	}
	delta := float64(last.val - first.val)
	if delta <= 0 {
		return 0
	}
	return delta / last.at.Sub(first.at).Seconds()
}

type Meter struct {
	mu      sync.Mutex
	verb    string
	name    string
	unit    Unit
	total   int64
	current int64
	item    string
	ok      int
	failed  int
	start   time.Time
	window  rateWindow
	drawn   int
	settled bool
	stop    chan struct{}
	wg      sync.WaitGroup
}

func NewMeter(verb, name string, total int64, unit Unit) *Meter {
	m := &Meter{
		verb:  verb,
		name:  name,
		unit:  unit,
		total: total,
		start: time.Now(),
		stop:  make(chan struct{}),
	}
	m.window.add(m.start, 0)

	interval := time.Second
	if m.live() {
		watchForInterrupt()
		fmt.Print("\033[?25l")
		interval = 100 * time.Millisecond
	}
	m.render()

	m.wg.Go(func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-m.stop:
				return
			case <-ticker.C:
				m.render()
			}
		}
	})
	return m
}

func (m *Meter) live() bool {
	return StdoutIsTerminal && !GlobalDebugFlag
}

func (m *Meter) Item(name string) {
	m.mu.Lock()
	m.item = name
	m.mu.Unlock()
}

func (m *Meter) Add(n int64) {
	m.mu.Lock()
	m.current += n
	m.ok += int(n)
	m.window.add(time.Now(), m.current)
	m.mu.Unlock()
}

func (m *Meter) Write(p []byte) (int, error) {
	m.Add(int64(len(p)))
	return len(p), nil
}

func (m *Meter) ItemFailed(name string, err error) {
	m.mu.Lock()
	m.failed++
	m.current++
	m.item = ""
	m.window.add(time.Now(), m.current)
	m.clear()
	PrintIndentedError(fmt.Sprintf("%s: %s", name, oneLine(err)), err)
	m.mu.Unlock()

	m.render()
}

func (m *Meter) Done() {
	m.finish(nil)
}

func (m *Meter) Fail(err error) {
	m.finish(err)
}

func (m *Meter) finish(err error) {
	m.mu.Lock()
	if m.settled {
		m.mu.Unlock()
		return
	}
	m.settled = true
	close(m.stop)
	m.clear()
	if m.live() {
		fmt.Print("\033[?25h")
	}
	line := m.settledLine()
	failed := m.failed
	set := m.isSet()
	m.mu.Unlock()
	m.wg.Wait()

	switch {
	case err != nil:
		PrintIndentedError(fmt.Sprintf("%s: %s", m.name, oneLine(err)), err)
	case failed > 0:
		PrintError(line, nil)
	case set:
		PrintInfo(line)
	default:
		PrintSuccess(line)
	}
}

func (m *Meter) isSet() bool {
	return m.unit != UnitBytes
}

func (m *Meter) opaque() bool {
	return m.unit != UnitBytes && m.total <= 1
}

func (m *Meter) clear() {
	if !m.live() || m.drawn == 0 {
		return
	}
	fmt.Print(strings.Repeat("\033[1A\033[2K", m.drawn))
	m.drawn = 0
}

func (m *Meter) render() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.settled {
		return
	}

	if GlobalDebugFlag {
		event := log.Info().
			Int64("current", m.current).
			Int64("total", m.total).
			Float64("rate", m.window.current()).
			Str("eta", m.eta())
		if m.total > 0 {
			event = event.Int("percent", m.percent())
		}
		if m.item != "" {
			event = event.Str("item", m.item)
		}
		event.Msg(strings.TrimSpace(m.verb + " " + m.name))
		return
	}

	if !StdoutIsTerminal {
		lipgloss.Println(m.header() + "  " + strings.Join(m.fields(), "  "))
		return
	}

	fields, bar := m.frame(TermWidth())
	lines := []string{m.header(), "  " + strings.TrimLeft(bar+"  "+strings.Join(fields, "  "), " ")}

	var b strings.Builder
	b.WriteString(strings.Repeat("\033[1A\033[2K", m.drawn))
	for _, l := range lines {
		b.WriteString("\r\033[2K" + l + "\n")
	}
	fmt.Print(b.String())
	m.drawn = len(lines)
}

func (m *Meter) header() string {
	label := strings.TrimSpace(m.verb + " " + clip(m.name, 60))
	head := infoStyle.Render("↻ " + label)
	if m.item == "" || strings.HasSuffix(label, m.item) {
		return head
	}
	return head + "  " + meterMutedStyle.Render(clip(m.item, 40))
}

func (m *Meter) percent() int {
	if m.total <= 0 {
		return 0
	}
	return min(max(int(float64(m.current)/float64(m.total)*100), 0), 100)
}

func (m *Meter) eta() string {
	rate := m.averageRate(rateFloor)
	if m.total <= 0 || rate <= 0 {
		return "unknown"
	}
	left := time.Duration(float64(m.total-m.current) / rate * float64(time.Second))
	if left < 0 || left > etaCeiling {
		return "unknown"
	}
	return formatEstimate(left)
}

func (m *Meter) averageRate(floor time.Duration) float64 {
	elapsed := time.Since(m.start)
	if m.current <= 0 || elapsed < floor {
		return 0
	}
	return float64(m.current) / elapsed.Seconds()
}

func (m *Meter) pairReserve() int {
	if !m.isSet() {
		return 14
	}
	return len(strconv.FormatInt(max(m.total, m.current), 10))*2 + 4 + len(m.unit)
}

type meterField struct {
	text     string
	reserved int
}

func (m *Meter) allFields() []meterField {
	if m.opaque() {
		return []meterField{{formatElapsed(time.Since(m.start)), 7}}
	}
	fields := make([]meterField, 0, 5)
	if m.total > 0 {
		fields = append(fields, meterField{fmt.Sprintf("%3d%%", m.percent()), 4})
	}
	fields = append(fields, meterField{formatAmount(m.current, m.total, m.unit), m.pairReserve()})
	if m.isSet() {
		return append(fields, meterField{"eta " + m.eta(), 11})
	}
	return append(fields,
		meterField{formatRate(m.window.current(), m.unit), 11},
		meterField{"eta " + m.eta(), 11},
		meterField{"avg " + formatRate(m.averageRate(rateFloor), m.unit), 15},
	)
}

func (m *Meter) fields() []string {
	all := m.allFields()
	texts := make([]string, len(all))
	for i, f := range all {
		texts[i] = f.text
	}
	return texts
}

func (m *Meter) frame(width int) ([]string, string) {
	all := m.allFields()
	for {
		reserved := 2
		for _, f := range all {
			reserved += f.reserved + 2
		}
		cells := width - reserved
		texts := make([]string, len(all))
		for i, f := range all {
			texts[i] = f.text
		}
		if cells >= barFloor {
			return texts, m.bar(min(cells, barCeiling))
		}
		if len(all) <= 2 {
			return texts, ""
		}
		all = all[:len(all)-1]
	}
}

func (m *Meter) bar(cells int) string {
	if m.total <= 0 || m.opaque() {
		w := max(3, cells/5)
		span := cells - w
		if span <= 0 {
			return infoStyle.Render(strings.Repeat("─", cells))
		}
		pos := int(time.Since(m.start)/indeterminateStep) % (2 * span)
		if pos > span {
			pos = 2*span - pos
		}
		return meterChromeStyle.Render(strings.Repeat("─", pos)) +
			infoStyle.Render(strings.Repeat("─", w)) +
			meterChromeStyle.Render(strings.Repeat("─", cells-pos-w))
	}
	filled := m.percent() * cells / 100
	if filled == 0 {
		return meterChromeStyle.Render(strings.Repeat("─", cells))
	}
	if filled == cells {
		return infoStyle.Render(strings.Repeat("─", cells))
	}
	return infoStyle.Render(strings.Repeat("─", filled)) +
		meterChromeStyle.Render(strings.Repeat("─", cells-filled))
}

func (m *Meter) settledLine() string {
	parts := []string{clip(m.name, 60)}
	if m.failed > 0 {
		parts = append(parts, fmt.Sprintf("%d ok, %d failed", m.ok, m.failed))
	} else {
		parts = append(parts, formatCount(m.current, m.unit))
	}
	parts = append(parts, formatElapsed(time.Since(m.start)))
	if rate := m.averageRate(0); rate > 0 && !m.isSet() {
		parts = append(parts, "avg "+formatRate(rate, m.unit))
	}
	return strings.Join(parts, "  ")
}

func oneLine(err error) string {
	return clip(strings.Join(strings.Fields(err.Error()), " "), 120)
}

func clip(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit-1]) + "…"
}

func scaleBytes(n float64) (float64, string) {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	idx := 0
	for n >= 1024 && idx < len(units)-1 {
		n /= 1024
		idx++
	}
	return n, units[idx]
}

func formatNumber(v float64) string {
	if v >= 100 {
		return strconv.FormatFloat(v, 'f', 0, 64)
	}
	return strconv.FormatFloat(v, 'f', 1, 64)
}

func formatAmount(current, total int64, unit Unit) string {
	if unit != UnitBytes {
		if total > 0 {
			return fmt.Sprintf("%d / %d %s", current, total, unit)
		}
		return fmt.Sprintf("%d %s", current, unit)
	}
	if total > 0 {
		scaled, name := scaleBytes(float64(total))
		factor := float64(total) / scaled
		return fmt.Sprintf("%s / %s %s", formatNumber(float64(current)/factor), formatNumber(scaled), name)
	}
	return formatCount(current, unit)
}

func formatCount(n int64, unit Unit) string {
	if unit != UnitBytes {
		noun := string(unit)
		if n == 1 {
			noun = strings.TrimSuffix(noun, "s")
		}
		return fmt.Sprintf("%d %s", n, noun)
	}
	scaled, name := scaleBytes(float64(n))
	return fmt.Sprintf("%s %s", formatNumber(scaled), name)
}

func formatRate(rate float64, unit Unit) string {
	if unit != UnitBytes {
		return fmt.Sprintf("%s %s/s", formatNumber(rate), unit)
	}
	scaled, name := scaleBytes(rate)
	return fmt.Sprintf("%s %s/s", formatNumber(scaled), name)
}

func formatElapsed(d time.Duration) string {
	if d < time.Minute {
		return strconv.FormatFloat(d.Seconds(), 'f', 1, 64) + "s"
	}
	return formatEstimate(d)
}

func formatEstimate(d time.Duration) string {
	switch {
	case d < time.Minute:
		return strconv.Itoa(int(d.Seconds())) + "s"
	case d < time.Hour:
		return fmt.Sprintf("%dm%02ds", int(d.Minutes()), int(d.Seconds())%60)
	default:
		return fmt.Sprintf("%dh%02dm", int(d.Hours()), int(d.Minutes())%60)
	}
}
