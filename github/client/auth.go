package github

import "net/http"

// Credentials contains user's name and access token to access the Github API
type Credentials struct {
	Username            string
	PersonalAccessToken string
}

// auth is a HTTPClient which sets authorization headers for Github API
type auth struct {
	HTTPClient // "back-end" HTTPClient to use for actual HTTP queries
	*Credentials
}

// Do is HTTPClient.Do
func (c auth) Do(req *http.Request) (*http.Response, error) {
	if creds := c.Credentials; creds != nil {
		req.SetBasicAuth(creds.Username, creds.PersonalAccessToken)
		req.Header.Add("User-Agent", creds.Username)
	}
	return c.HTTPClient.Do(req)

}
