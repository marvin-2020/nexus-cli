package main

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Almost completely ripped off https://www.socketloop.com/tutorials/golang-natural-string-sorting-example

type Compare func(str1, str2 string) bool

func (cmp Compare) Sort(strs []string) {
	strSort := &strSorter{
		strs: strs,
		cmp:  cmp,
	}
	sort.Sort(strSort)
}

type strSorter struct {
	strs []string
	cmp  func(str1, str2 string) bool
}

// str like 0.2.77.20240520-184-g17ff52e
func extractNumberFromString(str string) (num uint64) {
	index := strings.LastIndex(str, "-")
	if index > 0 {
		str = str[0:index]
	}
	lastSuf := ""
	index = strings.LastIndex(str, "-")
	if index > 0 {
		lastSuf = str[index+1:]
		str = str[0:index]
	}
	strSliceLastSuf := make([]string, 0)
	for _, v := range lastSuf {
		if unicode.IsDigit(v) {
			strSliceLastSuf = append(strSliceLastSuf, string(v))
		}
	}
	strSlice := make([]string, 0)
	for _, v := range str {
		if unicode.IsDigit(v) {
			strSlice = append(strSlice, string(v))
		}
	}

	// If the tag was all non-digits, the strSlice would be empty (e.g., 'latest')
	// therefore just throw it to the end (1 << 32 is maxint)
	if len(strSlice) == 0 {
		return 1 << 63
	}

	num, err := strconv.ParseUint(strings.Join(strSlice, ""), 10, 64)
	if err != nil {
		log.Println(str, " strSlice ", num)
		log.Println(err)
	}
	num2, err := strconv.ParseUint(strings.Join(strSliceLastSuf, ""), 10, 64)
	if err != nil {
		log.Println(str, " strSliceLastSuf ", num2)
		log.Println(err)
	}
	return num + num2
}

func (s *strSorter) Len() int { return len(s.strs) }

func (s *strSorter) Swap(i, j int) { s.strs[i], s.strs[j] = s.strs[j], s.strs[i] }

func (s *strSorter) Less(i, j int) bool { return s.cmp(s.strs[i], s.strs[j]) }
