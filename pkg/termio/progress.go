package termio

import (
	"fmt"
	"math"

	"github.com/muesli/goprogressbar"
)

func init() {
	goprogressbar.Stdout = Stderr
}

// ProgressBar is a wrapper around goprogressbar.
type ProgressBar struct {
	Bar    *goprogressbar.ProgressBar
	Hidden bool
}

// NewProgressBar creates a new progress bar.
func NewProgressBar(text string, total int64) *ProgressBar {
	return &ProgressBar{
		Bar: &goprogressbar.ProgressBar{
			Total:   total,
			Current: 0,
            Text:    text,
			Width:   60,
			PrependTextFunc: func(p *goprogressbar.ProgressBar) string {
				cur := p.Current
				max := p.Total
				digits := int(math.Log10(float64(max))) + 1
				// Log10(0) is undefined
				if max < 1 {
					digits = 1
				}
				return fmt.Sprintf(fmt.Sprintf(" %%%dd / %%%dd ", digits, digits), cur, max)
			},
		},
	}
}

// Add adds the given amount to the progress.
func (p *ProgressBar) Add(v int64) {
	cur := p.Bar.Current + v
	p.Bar.Current = min(cur, p.Bar.Total)
	p.print()
}

// Inc adds one to the progress.
func (p *ProgressBar) Inc() {
	if p.Bar.Current < p.Bar.Total {
		p.Bar.Current++
	}
	p.print()
}

// Set sets an arbitrary progress.
func (p *ProgressBar) Set(v int64) {
	p.Bar.Current = min(v, p.Bar.Total)
	p.print()
}

// SetText changes the status text.
// Does not display the updated progress bar, so set the status
// text before calling an update function such as `Inc`.
func (p *ProgressBar) SetText(text string) {
    p.Bar.Text = text
}

// Done finalizes the progress bar.
func (p *ProgressBar) Done() {
	if p.Hidden {
		return
	}

	fmt.Fprintln(Stderr, "")
}

// Clear removes the progress bar.
func (p *ProgressBar) Clear() {
	p.Bar.Clear()
}

// print will print the progress bar, if necessary.
func (p *ProgressBar) print() {
	if ! p.Hidden {
        p.Bar.LazyPrint()
	}
}

// `current` gets the current progress value (mostly for testing)
func (p *ProgressBar) current() int64 {
    return p.Bar.Current
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}
