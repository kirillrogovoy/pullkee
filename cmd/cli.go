package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kirillrogovoy/pullk/github"
)

const usage = `Usage:
	pullk [flags] [repo]
	repo - Github repository path as "username/reponame"

	Flags:
	--limit - Only use N last pull requests
	--reset - Reset the local cache and fetch all the data through the Github API again

	Environment variables:
	GITHUB_CREDS - API credentials in the format "username:personal_access_token"`

type flags struct {
	limit int
	reset bool
}

func getFlags() flags {
	flags := flags{}

	flag.IntVar(&flags.limit, "limit", 0, "")
	flag.BoolVar(&flags.reset, "reset", false, "")

	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()

	return flags
}

func getGithubCreds() *github.Credentials {
	env := os.Getenv("GITHUB_CREDS")

	if env == "" {
		return nil
	}

	if matches, _ := regexp.MatchString(`^[\w-]+:[\w-]+$`, env); !matches {
		fmt.Printf("Invalid format of the GITHUB_CREDS environment variable!\n\n%s\n", usage)
		os.Exit(1)
	}

	creds := strings.Split(env, ":")
	return &github.Credentials{
		Username:            creds[0],
		PersonalAccessToken: creds[1],
	}
}

func getRepo() string {
	repo := flag.Arg(0)

	if repo == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	if matches, _ := regexp.MatchString(`^[\w-\.]+/[\w-\.]+$`, repo); !matches {
		fmt.Printf("Invalid format of the repo!\n\n%s\n", usage)
		os.Exit(1)
	}

	if len(flag.Args()) > 1 {
		fmt.Println("Excess arguments after the repo")
		fmt.Println(usage)
		os.Exit(1)
	}

	return repo
}
