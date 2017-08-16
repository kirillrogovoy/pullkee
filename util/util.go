package util

import "strconv"

// ParseInt removes clutter from converting strings to ints
func ParseInt(input string) int {
	i, err := strconv.ParseInt(input, 10, 64)

	if err != nil {
		panic(err)
	}

	return int(i)
}
