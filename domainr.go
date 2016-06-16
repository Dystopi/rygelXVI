package rygelXVI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

type DomainrStatusResponse struct {
	Status []DomainrStatus `json:"status"`
}

type DomainrSearchResponse struct {
	Results []DomainrSearch `json:"results"`
}

type DomainrStatus struct {
	Domain  string `json:"domain"`
	Zone    string `json:"zone"`
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

type DomainrSearch struct {
	Domain      string `json:"domain"`
	Host        string `json:"host"`
	Subdomain   string `json:"subdomain"`
	Zone        string `json:"zone"`
	Path        string `json:"path"`
	RegisterURL string `json:"registerURL"`
}

func NewClient(apiKey string, httpClient *http.Client) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API Key must not be empty")
	}
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{apiKey, httpClient}, nil
}

func (c *Client) SearchDomainr(domain string) (*DomainrSearchResponse, error) {
	searchURL, err := url.Parse("https://domainr.p.mashape.com/v2/search")
	if err != nil {
		return nil, err
	}
	qp := searchURL.Query()
	qp.Add("query", domain)
	searchURL.RawQuery = qp.Encode()
	body, err := c.makeRequest("GET", searchURL)
	if err != nil {
		return nil, err
	}
	var searchResults DomainrSearchResponse
	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, err
	}
	return &searchResults, nil
}

func (c *Client) SearchActive(domain string) (*DomainrSearchResponse, error) {
	searchResults, err := c.SearchDomainr(domain)
	if err != nil {
		return nil, err
	}
	var suggestedDomains []string
	dMap := make(map[string]DomainrSearch)
	for _, res := range searchResults.Results {
		suggestedDomains = append(suggestedDomains, res.Domain)
		dMap[res.Domain] = res
	}
	// Return suggested Domains. New method for Search Active
	activeMap, err := c.DomainrStatus(suggestedDomains)
	if err != nil {
		return nil, err
	}
	var activeDomains DomainrSearchResponse
	for key, value := range activeMap {
		if value == true {
			activeDomains.Results = append(activeDomains.Results, dMap[key])
		}
	}

	return &activeDomains, nil
}

func (c *Client) DomainrStatus(domains []string) (map[string]bool, error) {
	statusURL, err := url.Parse("https://domainr.p.mashape.com/v2/status")
	if err != nil {
		return nil, err
	}
	domainString := strings.Join(domains, ",")
	qp := statusURL.Query()
	qp.Set("domain", domainString)
	statusURL.RawQuery = qp.Encode()

	body, err := c.makeRequest("GET", statusURL)
	if err != nil {
		return nil, err
	}

	var statusResponse DomainrStatusResponse
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		return nil, err
	}
	activeMap := make(map[string]bool)

	for _, status := range statusResponse.Status {
		activeMap[status.Domain] = defineActive(status)
	}
	return activeMap, nil
}

func (c *Client) makeRequest(method string, requestURL *url.URL) ([]byte, error) {
	req, err := http.NewRequest(method, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Mashape-Key", c.APIKey)
	req.Header.Add("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received a non-acceptable response code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func defineActive(status DomainrStatus) bool {
	switch status.Status {
	case "inactive":
		return true
	case "undelegated":
		return true
	case "undelegated inactive":
		return true
	case "marketed":
		return true
	case "priced":
		return true
	case "transferable":
		return true
	case "premium":
		return true
	default:
		return false
	}
}
