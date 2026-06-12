package removecommontokens

import (
	"strings"
)

func removeId(s string)string {
	s = strings.ToLower(s)

	tagOpenInd := strings.LastIndex(s, "[")
	if tagOpenInd != -1 {
		s = s[:tagOpenInd]
	}

	return strings.TrimSpace(s)
}

func removeTokens(tokens []string, name string)string{
	commonSet := map[string]bool {}

	for _, token := range tokens {
		commonSet[token] = true
	}

	nameTokens := strings.Fields(name)

	result := []string{}

	for _, token := range nameTokens {
		if !commonSet[token] {
			result = append(result, token)
		}
	}

	return strings.Join(result, " ")
}

func RemoveCommonTokens(files []string, thresholdPercent int) []string{

	counts := map[string]int {}

	for _, file := range files {
		cleanFile := removeId(file)
		fileTokens := strings.Fields(cleanFile)

		seen := map[string]bool {}

		for _, token := range fileTokens {
			if !seen[token] {
				counts[token]++
				seen[token] = true
			}
		}
	}

	cTokens := []string{}
	for key, val := range counts {
		if val >= len(files)*thresholdPercent / 100 {
			cTokens = append(cTokens, key)
		}
	}

	res := []string {}

	for _, file :=range files {
		noIdFile := removeId(file)
		cleanFile := removeTokens(cTokens, noIdFile)

		res = append(res, cleanFile)
	}

	return res
}