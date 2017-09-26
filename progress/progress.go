// Package progress provides an utility to print a simple progress bar
// in the command line
package progress

import (
	"fmt"
	"math"
)

// Bar contains the information about the progress bar and the hook to call
// once it receives the new state
type Bar struct {
	Len      int
	OnChange func(s string)
}

// Set expects `v` to be from 0 to 1 and updates the progress bar according to the value
func (b Bar) Set(v float64) {
	if v < 0 || v > 1 {
		panic(fmt.Sprintf("The value passed should be between 0 and 1, %f given", v))
	}

	b.OnChange(render(v, b.Len))
}

func render(v float64, len int) string {
	bar := ""

	cellsFilled := int(math.Floor(v * float64(len)))

	for i := 0; i < cellsFilled; i++ {
		bar += "#"
	}

	for i := 0; i < len-cellsFilled; i++ {
		bar += " "
	}

	return fmt.Sprintf("[%s]", bar)
}
