package github

import "strings"

func FindStringByLineNumber(content string, findingStr string) []int {

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, findingStr) {
			return []int{i + 1}
		}
	}
	return []int{}

}
