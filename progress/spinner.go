package progress

import (
	"fmt"
	"runtime"
	"time"

	"github.com/git-lfs/git-lfs/git/githistory/log"
)

// Indeterminate progress indicator 'spinner'
type Spinner struct {
	stage int
	msg   string

	updates chan *log.Update
}

var spinnerChars = []byte{'|', '/', '-', '\\'}

// Print a spinner (stage) followed by msg (no linefeed)
func (s *Spinner) Print(msg string) {
	s.msg = msg
	s.Spin()
}

// Just spin the spinner one more notch & use the last message
func (s *Spinner) Spin() {
	s.stage = (s.stage + 1) % len(spinnerChars)
	s.update(string(spinnerChars[s.stage]), s.msg)
}

// Finish the spinner with a completion message & newline
func (s *Spinner) Finish(finishMsg string) {
	s.msg = finishMsg
	s.stage = 0

	var sym string
	if runtime.GOOS == "windows" {
		// Windows console(s) cannot display UTF-8 check marks except in
		// ConEmu (not cmd or git bash).
		sym = "*"
	} else {
		sym = fmt.Sprintf("%c", '\u2714')
	}

	s.update(sym, finishMsg)
	close(s.updates)
}

func (s *Spinner) Updates() <-chan *log.Update { return s.updates }
func (s *Spinner) Throttled() bool             { return false }

func (s *Spinner) update(prefix, msg string) {
	s.updates <- &log.Update{
		S:  fmt.Sprintf("%s %s", prefix, msg),
		At: time.Now(),
	}
}

func NewSpinner() *Spinner {
	return &Spinner{
		updates: make(chan *log.Update),
	}
}
