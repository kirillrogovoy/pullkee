package metric

import (
	"fmt"
	"sort"
)

// average is a helper for calculating an average value of something
type average struct {
	Sum   float64
	Count int
}

// Calculate divides Sum by Count
func (a average) Calculate() float64 {
	return a.Sum / float64(a.Count)
}

// averageMap is a map which contains an average value for each dev
type averageMap map[string]*average

// toList converts the map to a sortable list of AverageItem
func (a *averageMap) toList() averageList {
	var d averageList
	for name, average := range *a {
		d = append(d, averageItem{name, average.Calculate()})
	}
	return d
}

func (a *averageMap) reset() {
	(*a) = averageMap{}
}

func (a *averageMap) add(value float64, dev string) {
	if _, ok := (*a)[dev]; !ok {
		(*a)[dev] = &average{}
	}
	av := (*a)[dev]
	av.Sum += value
	av.Count++
}

func (a *averageMap) setCount(i int) {
	for key := range *a {
		(*a)[key].Count = i
	}
}

func (a *averageMap) string(unit string) string {
	d := a.toList()

	averageTotal := 0.0
	for _, dev := range d {
		averageTotal += dev.Value
	}

	result := ""
	result += fmt.Sprintf("Total average: %.2f %s\n", averageTotal/float64(len(d)), unit)

	sort.Sort(sort.Reverse(d))

	for _, dev := range d {
		result += fmt.Sprintf("Average for %s: %.2f %s\n", dev.Name, dev.Value, unit)
	}

	return result
}

// averageItem is a simple container for dev's name and some average value
type averageItem struct {
	Name  string
	Value float64
}

// averageList is used to sort a slice of AverageItems by Value
type averageList []averageItem

func (a averageList) Len() int {
	return len(a)
}
func (a averageList) Swap(x, y int) {
	a[x], a[y] = a[y], a[x]
}
func (a averageList) Less(x, y int) bool {
	return a[x].Value < a[y].Value
}

func (a averageList) string(unit string) string {
	averageTotal := 0.0
	for _, dev := range a {
		averageTotal += dev.Value
	}

	result := ""
	result += fmt.Sprintf("Total average: %.2f %s\n", averageTotal/float64(len(a)), unit)

	sort.Sort(sort.Reverse(a))

	for _, dev := range a {
		result += fmt.Sprintf("Average for %s: %.2f %s\n", dev.Name, dev.Value, unit)
	}

	return result
}

type counterMap map[string]*counter

func (c counterMap) string() string {
	var cs counters

	for _, i := range c {
		cs = append(cs, *i)
	}

	sort.Sort(sort.Reverse(cs))

	result := ""
	for _, i := range cs {
		result += fmt.Sprintf("For %s: %d\n", i.Name, i.Count)
	}

	return result
}

type counter struct {
	Name  string
	Count int
}
type counters []counter

func (a counters) Len() int {
	return len(a)
}
func (a counters) Swap(x, y int) {
	a[x], a[y] = a[y], a[x]
}
func (a counters) Less(x, y int) bool {
	return a[x].Count < a[y].Count
}
