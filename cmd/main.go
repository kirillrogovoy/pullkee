package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/kirillrogovoy/pullk/github"
	"github.com/kirillrogovoy/pullk/metric"
	"github.com/kirillrogovoy/pullk/repository"
	"github.com/kirillrogovoy/pullk/util"
)

func main() {
	flags := getFlags()
	creds := getGithubCreds()
	repo := getRepo()
	api, err := getAPI(creds, repo)
	if err != nil {
		reportFetchingError(err)
	}

	r := getRepository(api, repo, flags)

	printRateDetails(api)

	if pulls, err := getPulls(r); err == nil {
		runMetrics(pulls)
	} else {
		reportFetchingError(err)
	}
}

func getAPI(creds *github.Credentials, repo string) (*github.API, error) {
	limiter := time.Tick(time.Millisecond * 75)
	api := &github.API{
		Creds:       creds,
		RateLimiter: &limiter,
	}

	// checking that we can successfully open the requested repo
	_, err := api.Repo(repo)
	return api, err
}

func getRepository(api *github.API, repo string, flags flags) *repository.Repository {
	return &repository.Repository{
		API: api,
		Settings: repository.Settings{
			Repo:        repo,
			Limit:       flags.limit,
			ResetCaches: flags.reset,
		},
	}
}

func getPulls(r *repository.Repository) (github.PullRequests, error) {
	fmt.Println("Getting Pull Request list...")
	return r.Pulls()
}

func reportFetchingError(err error) {
	fmt.Printf("An unexpected error occurred during information fetching\n\n%s\n", err)
	os.Exit(1)
}

func printRateDetails(api *github.API) {
	fmt.Printf(
		"Github rate limit details:\nLimit: %s\nRemaining: %s\nReset: %s\n\n",
		api.LastHeader.Get("X-RateLimit-Limit"),
		api.LastHeader.Get("X-RateLimit-Remaining"),
		time.Unix(int64(util.ParseInt(api.LastHeader.Get("X-RateLimit-Reset"))), 0).String(),
	)
}

func runMetrics(pullRequests github.PullRequests) {
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
