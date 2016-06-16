package rygelXVI

import (
	"testing"
)

func TestValidateDomain(t *testing.T) {
	goodDomain := "foo.co"
	isValid, err := validateDomain(goodDomain)
	if err != nil {
		t.Log("Received an unexpected error from domain validation")
		t.Fail()
	} else {
		t.Log("Succesfully validated Domain")
	}
	if isValid != true {
		t.Log("Received a bad validation response")
		t.Fail()
	} else {
		t.Log("Succesfully received the correct validation response")
	}

	badDomain := "http://www.foo.co"
	isValid, err = validateDomain(badDomain)
	if err != nil {
		t.Log("Received an unexpected error from domain validation")
		t.Fail()
	} else {
		t.Log("Succesfully validated Domain")
	}
	if isValid != false {
		t.Log("Received a bad validation response")
		t.Fail()
	} else {
		t.Log("Succesfully received the correct validation response")
	}
}
