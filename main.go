package progbar

// "self_utilities/progress_bar"

import (
	"fmt"
	"os"
	"errors"
)

const empty = byte(' ')
const filled = byte('=')

// Total is the number of jobs to be completed for the progress bar
// Length is the number of segments for the bar to have
// Bar will have the bar segments
type Progress struct {
	Total	int
	Length	int
	Bar	[]byte
}

var ErrorNotTTY = errors.New("This interface is not a TTY")

// Returns a [Progress] object if the current terminal interface is a TTY
// The returned error will be nil or [ErrorNotTTY]
func NewProgress(total, length int) (p *Progress, err error) {
	if isTTY() {
		p = &Progress{Total: total, Length: length, Bar: func(l int) []byte{
			v := make([]byte, l)
			for i:=0; i<len(v); i++ {
				v[i] = empty
			}
			return v
		}(length) }
	} else {
		err = ErrorNotTTY
	}
	return
}

// Returns the ratio of n/d, unless d is 0, in which case it returns 0 to prevent a div/0 panic
func R(n, d int) float64 {
	if d == 0 {
		return 0
	}
	return float64(n)/float64(d)
}

// Prevents printing nonsense if a progress bar isn't properly initialized
func(p *Progress) Valid() bool {
	if p.Total < 1 || p.Length < 1 {
		return false
	}
	return true
}

// Prints a string with a carriage return if the output device is a TTY
func Status(s string) {
	if isTTY() {
		fmt.Printf("%s \r", s)
	}
}

// Generates a status line with a header, progress bar, status, and value / total counter
func(p *Progress) Get (header, status string, value int) string {
	p.Update(value)

	return fmt.Sprintf("%s[%s] %s %d / %d \r", header, p.Bar, status, value, p.Total)
}

// Generates and draws a status line from [Get]
func(p *Progress) Draw (header, status string, value int) {
	if p.Valid() {
		Status(p.Get(header, status, value))
	}
}

// Updates the segments of the bar to represent the current percentage of jobs done
// It will not show the bar as full unless int(R(count/total)) == 1
func (p *Progress) Update (count int) {
	for i := 0; i < p.Length; i++ {
		p.Bar[i] = empty
		if R(i,p.Length) < R(count,p.Total) {
			p.Bar[i] = filled
		}
	}
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

