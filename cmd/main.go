package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"github.com/kirillrogovoy/pullk/cache"
	"github.com/kirillrogovoy/pullk/github"
	client "github.com/kirillrogovoy/pullk/github/client"
	github_util "github.com/kirillrogovoy/pullk/github_util"
	"github.com/kirillrogovoy/pullk/metric"
)

func main() {
	flags := getFlags()

	client := getHTTPClient(getGithubCreds())
	api := getAPI(&client, getRepo())

	// check that we can at least successfully fetch repository's meta information
	if _, err := api.Repository(); err != nil {
		reportErrorAndExit(err)
	}

	printRateDetails(client)

	cache := getCache()
	pulls := getPulls(flags, api, cache)

	runMetrics(pulls)
}

func getHTTPClient(creds *client.Credentials) client.Client {
	rateLimiter := time.Tick(time.Millisecond * 75)
	return client.New(http.DefaultClient, client.Options{
		Credentials: creds,
		RateLimiter: &rateLimiter,
		MaxRetries:  3,
		Log:         func(msg string) { fmt.Println(msg) },
	})

}

func getAPI(client client.HTTPClient, repo string) github.APIv3 {
	api := github.APIv3{
		RepoName:   repo,
		HTTPClient: client,
	}

	return api
}

func getCache() cache.Cache {
	return cache.FSCache{
		CachePath: path.Join(os.TempDir(), "pullk_cache"),
		FS:        RealFS{},
	}
}

func getPulls(f flags, a github.API, c cache.Cache) []github.PullRequest {
	fmt.Println("Getting Pull Request list...")

	pulls, err := github_util.Pulls(a, f.limit)
	if err != nil {
		reportErrorAndExit(err)
	}

	if err := github_util.FillDetails(a, c, pulls); err != nil {
		reportErrorAndExit(err)
	}

	return pulls
}

func reportErrorAndExit(err error) {
	fmt.Printf("An unexpected error occurred:\n\n%s\n", err)
	os.Exit(1)
}

func printRateDetails(c client.Client) {
	l := c.LastResponse.Header
	resetAt, _ := strconv.Atoi(l.Get("X-RateLimit-Reset"))

	fmt.Printf(
		"Github rate limit details:\nLimit: %s\nRemaining: %s\nReset: %s\n\n",
		l.Get("X-RateLimit-Limit"),
		l.Get("X-RateLimit-Remaining"),
		time.Unix(int64(resetAt), 0).String(),
	)
}

func runMetrics(pullRequests []github.PullRequest) {
	for _, m := range metric.Metrics() {
		name := reflect.TypeOf(m).Elem().Name()
		description := m.Description()
		err := m.Calculate(pullRequests)
		fmt.Printf("Metric '%s' (%s)\n", name, description)
		if err != nil {
			log.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", m.String())
		}
	}
}
