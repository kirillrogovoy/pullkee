package client

import "net/http"

// Credentials struct contains user's name and access token to access the Github API
type Credentials struct {
	Username            string
	PersonalAccessToken string
}

// authenticating is a HTTPClient which sets authorization headers for Github API
type authenticating struct {
	HTTPClient // "back-end" HTTPClient to use for actual HTTP queries
	*Credentials
}

// Do is HTTPClient.Do
func (c authenticating) Do(req *http.Request) (*http.Response, error) {
	if creds := c.Credentials; creds != nil {
		req.SetBasicAuth(creds.Username, creds.PersonalAccessToken)
		req.Header.Add("User-Agent", creds.Username)
	}
	return c.HTTPClient.Do(req)

}
