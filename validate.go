package rygelXVI

import (
	"regexp"
)

func validateDomain(domain string) (bool, error) {
	// No http://www.
	httpRegex := regexp.MustCompile(`https?:\/\/(?:www.)?`)
	if httpRegex.MatchString(domain) {
		return false, nil
	}
	return true, nil
}
