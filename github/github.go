package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/kirillrogovoy/pullk/util"
)

// Credentials contains user's name and access token to access the Github API
type Credentials struct {
	Username            string
	PersonalAccessToken string
}

// API is the main entry point for calling methods
type API struct {
	Creds              *Credentials
	LastHeader         *http.Header
	RateLimiter        *<-chan time.Time
	abuseMechnismTimer *time.Timer
}

// User is a representation of a Github user (e.g. an author of a Pull Request)
type User struct {
	Login string `json:"login"`
}

// Repo fetches the remote information about a repository
func (a *API) Repo(repo string) (res *http.Response, err error) {
	res, err = a.send(request(fmt.Sprintf("https://api.github.com/repos/%s", repo)))
	return
}

func (a *API) send(req *http.Request) (*http.Response, error) {
	client := client()

	if creds := a.Creds; creds != nil {
		setAuth(req, *creds)
	}

	if a.abuseMechnismTimer != nil {
		<-(*a).abuseMechnismTimer.C
		a.abuseMechnismTimer = nil
	}

	if a.RateLimiter != nil {
		<-*a.RateLimiter
	}

	res, err := sendWithRetry(client, req)

	if err != nil {
		return nil, err
	}

	if abuseMechanismTriggered(res) {
		retryAfter := util.ParseInt(res.Header.Get("Retry-After"))
		duration := time.Duration(retryAfter) * time.Second
		a.abuseMechnismTimer = time.NewTimer(duration)

		res, err = sendWithRetry(client, req)

		if err != nil {
			return nil, err
		}
	}

	a.LastHeader = &res.Header

	log.Printf(
		"DONE - %s: %s",
		req.Method,
		req.URL.String(),
	)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		return nil, composeHTTPError(req, res)
	}

	return res, err
}

func sendWithRetry(client *http.Client, req *http.Request) (*http.Response, error) {
	res, err := client.Do(req)

	retriesMax := 4
	retriesLeft := retriesMax
	for err != nil && retriesLeft > 0 {
		retriesLeft--
		res, err = client.Do(req)
	}

	return res, err
}

func abuseMechanismTriggered(res *http.Response) bool {
	return res.StatusCode == 403 && res.Header.Get("Retry-After") != ""
}

func setAuth(req *http.Request, creds Credentials) {
	req.SetBasicAuth(creds.Username, creds.PersonalAccessToken)
	req.Header.Add("User-Agent", creds.Username)
}

func client() *http.Client {
	return http.DefaultClient
}

func request(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	return req
}

// HTTPBody is a helper to extract the body from a http.Response
func HTTPBody(res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	res.Body.Close()
	return body
}

func composeHTTPError(req *http.Request, res *http.Response) error {
	dump, err := httputil.DumpResponse(res, true)

	if err != nil {
		panic(err)
	}

	return fmt.Errorf(
		"HTTP Request failed.\nURL: %s\n\n%s",
		req.URL.String(),
		dump,
	)
}
