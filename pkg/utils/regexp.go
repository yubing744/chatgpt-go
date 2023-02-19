package utils

import "regexp"

func RegexpExtra(text string, pattern string, groupIndex int) (string, bool) {
	var re = regexp.MustCompile(pattern)

	matchs := re.FindAllStringSubmatch(text, -1)
	if len(matchs) > 0 {
		match := matchs[0]
		if groupIndex >= 0 && groupIndex < len(match) {
			return match[groupIndex], true
		}
	}

	return "", false
}
