package input

import "regexp"

func validateUrl(url string) bool {
	match, _ := regexp.MatchString(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,63}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`, url)
	return match
}
