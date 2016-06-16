package rygelXVI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/Dystopi/navCrystal"
)

func newStatusTestClient(callback func(http.ResponseWriter, *http.Request)) (*Client, error) {
	mockClient, _ := navCrystal.NewMock(callback)
	testClient, err := NewClient("7es74p1k3y24r3c00l", mockClient.Client)
	if err != nil {
		return nil, err
	}
	return testClient, nil
}

func statusReqHandlerSAD(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"status":[{"domain": "superawesomedomain.sl","zone": "sl","status": "undelegated","summary": "undelegated"}]}`)
}

func statusReqHandlerMiley(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"status":[{"domain":"miley.co","zone":"co","status":"active parked","summary":"parked"}]}`)
}

func statusReqHandlerShouldNotExist(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"status":[{"domain":"thisdomainshouldnotexist.com","zone":"com","status":"undelegated inactive","summary":"inactive"}]}`)
}

func searchReqHandlerMiley(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, `{"results":[{"domain":"miley.co","host":"","subdomain":"miley.","zone":"co","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.co&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.construction","host":"","subdomain":"miley.","zone":"construction","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.construction&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.contractors","host":"","subdomain":"miley.","zone":"contractors","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.contractors&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.consulting","host":"","subdomain":"miley.","zone":"consulting","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.consulting&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.community","host":"","subdomain":"miley.","zone":"community","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.community&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.computer","host":"","subdomain":"miley.","zone":"computer","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.computer&gl=US%2CGardena%2CUS-CA®istrar=&source="}]}`)
}

func multiReqHandlerMiley(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.URL.Path {
	case "/v2/search":
		fmt.Fprintln(w, `{"results":[{"domain":"miley.co","host":"","subdomain":"miley.","zone":"co","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.co&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.construction","host":"","subdomain":"miley.","zone":"construction","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.construction&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.contractors","host":"","subdomain":"miley.","zone":"contractors","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.contractors&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.consulting","host":"","subdomain":"miley.","zone":"consulting","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.consulting&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.community","host":"","subdomain":"miley.","zone":"community","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.community&gl=US%2CGardena%2CUS-CA®istrar=&source="},{"domain":"miley.computer","host":"","subdomain":"miley.","zone":"computer","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=mashape-verticalmass&domain=miley.computer&gl=US%2CGardena%2CUS-CA®istrar=&source="}]}`)
	case "/v2/status":
		fmt.Fprintln(w, `{"status":[{"domain":"miley.co","zone":"co","status":"active parked","summary":"parked"},{"domain":"miley.contractors","zone":"contractors","status":"undelegated inactive","summary":"inactive"},{"domain":"miley.construction","zone":"construction","status":"undelegated inactive","summary":"inactive"},{"domain":"miley.community","zone":"community","status":"undelegated inactive","summary":"inactive"},{"domain":"miley.consulting","zone":"consulting","status":"undelegated inactive","summary":"inactive"},{"domain":"miley.computer","zone":"computer","status":"undelegated inactive","summary":"inactive"}]}`)
	}
}

func TestMakeRequest(t *testing.T) {
	testClient, err := newStatusTestClient(statusReqHandlerSAD)
	testURL, _ := url.Parse("https://domainr.p.mashape.com/v2/status")
	qp := testURL.Query()
	qp.Add("domain", "superawesomedomain.sl")
	testURL.RawQuery = qp.Encode()
	rawJson, err := testClient.makeRequest("GET", testURL)
	if err != nil {
		t.Log("Failed to sucessfully make Domainr request")
		t.Fail()
	} else {
		t.Log("Succesfully made Domainr request")
	}
	var statusResponse DomainrStatusResponse
	json.Unmarshal(rawJson, &statusResponse)
	if statusResponse.Status[0].Status != "undelegated" {
		t.Log("Failed to receive the expected response")
		t.Fail()
	} else {
		t.Log("Succesfully received the expected response")
	}
}

func TestSearchDomainr(t *testing.T) {
	testClient, _ := newStatusTestClient(searchReqHandlerMiley)
	baseDomain := "miley.co"
	domains, err := testClient.SearchDomainr(baseDomain)
	if err != nil {
		t.Log("Failed to sucessfully search Domainr")
		t.Fail()
	} else {
		t.Log("Succesfully searched Domainr")
	}

	if len(domains.Results) != 6 {
		t.Log("Recieved less than expected results")
		t.Fail()
	} else {
		t.Log("Succesfully recieved the minimum expected results")
	}
}

func TestSearchActive(t *testing.T) {
	testClient, _ := newStatusTestClient(multiReqHandlerMiley)
	baseDomain := "miley.co"
	domains, err := testClient.SearchActive(baseDomain)
	if err != nil {
		t.Log("Failed to sucessfully search Domainr")
		t.Fail()
	} else {
		t.Log("Succesfully searched Domainr")
	}

	if len(domains.Results) != 5 {
		t.Log("Recieved less than expected results")
		t.Fail()
	} else {
		t.Log("Succesfully recieved the minimum expected results")
	}
}

func TestDomainrStatus(t *testing.T) {
	testClient, _ := newStatusTestClient(statusReqHandlerMiley)
	takenDomain := "miley.co"
	taken, err := testClient.DomainrStatus([]string{takenDomain})
	if err != nil {
		t.Log("Failed to sucessfully check Domainr")
		t.Fail()
	} else {
		t.Log("Succesfully checked Domainr")
	}
	if taken["miley.co"] != false {
		t.Log("Failed to recieve the expected response for miley.co")
		t.Fail()
	} else {
		t.Log("Succesfully recieved the expected response for miley.co")
	}

	testClient, _ = newStatusTestClient(statusReqHandlerShouldNotExist)
	availableDomain := "thisdomainshouldnotexist.com"
	taken, err = testClient.DomainrStatus([]string{availableDomain})
	if err != nil {
		t.Log("Failed to sucessfully check Domainr")
		t.Fail()
	} else {
		t.Log("Succesfully checked Domainr")
	}
	if taken[availableDomain] != true {
		t.Log("Failed to recieve the expected response for available domain")
		t.Fail()
	} else {
		t.Log("Succesfully recieved the expected response for available domain")
	}
}
