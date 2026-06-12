package tokencompare

import (
	"slices"
	"strings"
)

func TokenCompare(s1, s2 string) int {
	s1t := strings.Fields(s1)
	s2t := strings.Fields(s2)

	common := []string{}

	for _, token := range s1t {
		if slices.Contains(s2t, token) {
			// counts[token] = true
			common = append(common, token)
		}
	}

	similarity := int(float64(len(common)) / (float64(len(s1t)) + float64(len(s2t)) - float64(len(common))) * 100)
	return similarity
}