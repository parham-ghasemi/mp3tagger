package compare

import (
	"strings"
)

// lowWeight is the value given to tokens that should have less impact.
const lowWeight = 0.3

var lowWeightTokens = map[string]float64{
	"official": lowWeight,
	"part":     lowWeight,
	"pt":       lowWeight,
	"feat":     lowWeight,
	"ft":       lowWeight,
	"audio":    lowWeight,
}

func Normalize(s string) string {
	s = strings.ToLower(s)
	replacer := strings.NewReplacer(".", "", ",", "", "(", "", ")", "", "-", " ", "_", " ")
	return replacer.Replace(s)
}

func getTokenWeight(token string) float64 {
	if weight, exists := lowWeightTokens[token]; exists {
		return weight
	}
	return 1.0
}

func Score(s1, s2 string) int {
	s1t := strings.Fields(Normalize(s1))
	s2t := strings.Fields(Normalize(s2))

	set1 := make(map[string]bool)
	for _, token := range s1t {
		set1[token] = true
	}

	set2 := make(map[string]bool)
	for _, token := range s2t {
		set2[token] = true
	}

	union := make(map[string]bool)
	for token := range set1 {
		union[token] = true
	}
	for token := range set2 {
		union[token] = true
	}

	var intersectionWeight float64
	var unionWeight float64

	for token := range union {
		weight := getTokenWeight(token)
		unionWeight += weight

		if set1[token] && set2[token] {
			intersectionWeight += weight
		}
	}

	if unionWeight == 0 {
		return 0
	}

	similarity := int((intersectionWeight / unionWeight) * 100)
	return similarity
}
