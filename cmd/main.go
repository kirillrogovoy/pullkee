// Package cmd is the entry point of the app
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"github.com/kirillrogovoy/pullkee/cache"
	"github.com/kirillrogovoy/pullkee/github"
	"github.com/kirillrogovoy/pullkee/github/client"
	"github.com/kirillrogovoy/pullkee/github/util"
	"github.com/kirillrogovoy/pullkee/metric"
	"github.com/kirillrogovoy/pullkee/progress"
)

// Main is the entry function called by the "main" package
func Main() {
	flags := getFlags()

	client := getHTTPClient(getGithubCreds())
	repo := getRepo()
	api := getAPI(&client, repo)

	// check that we can at least successfully fetch repository's meta information
	if _, err := api.Repository(); err != nil {
		reportErrorAndExit(err)
	}

	printRateDetails(client)

	cache := getCache(repo)
	pulls := getPulls(flags, api, cache)

	runMetrics(pulls)
}

func getHTTPClient(creds *client.Credentials) client.Client {
	rateLimiter := time.Tick(time.Millisecond * 100)
	return client.New(http.DefaultClient, client.Options{
		Credentials: creds,
		RateLimiter: &rateLimiter,
		MaxRetries:  3,
	})
}

func getAPI(client client.HTTPClient, repo string) github.APIv3 {
	return github.APIv3{
		RepoName:   repo,
		HTTPClient: client,
	}
}

func getCache(repo string) cache.Cache {
	return cache.FSCache{
		CachePath: path.Join(os.TempDir(), "pullkee_cache", repo),
		FS:        RealFS{},
	}
}

func getPulls(f flags, a github.API, c cache.Cache) []github.PullRequest {
	fmt.Println("Getting Pull Request list...")

	pulls, err := util.Pulls(a, f.limit)
	if err != nil {
		reportErrorAndExit(err)
	}

	fmt.Println("Attaching details...")

	bar := progress.Bar{
		Len: 50,
		OnChange: func(v string) {
			fmt.Printf("\r%s", v)
		},
	}
	bar.Set(0)

	ch := util.FillDetails(
		a,
		c,
		pulls,
	)

	for i := range pulls {
		if err := <-ch; err != nil {
			reportErrorAndExit(err)
		}
		bar.Set(float64(i+1) / float64(len(pulls)))
	}
	fmt.Print("\n\n")

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
