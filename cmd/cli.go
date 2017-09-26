package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kirillrogovoy/pullkee/github/client"
)

const usage = `Usage:
	pullkee [flags] [repo]
	repo - Github repository path as "username/reponame"

	Flags:
	--limit - Only use N last pull requests

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

func getGithubCreds() *client.Credentials {
	creds := os.Getenv("GITHUB_CREDS")

	if creds == "" {
		return nil
	}

	if matches, _ := regexp.MatchString(`^[\w-]+:[\w-]+$`, creds); !matches {
		fmt.Printf("Invalid format of the GITHUB_CREDS environment variable!\n\n%s\n", usage)
		os.Exit(1)
	}

	split := strings.Split(creds, ":")
	return &client.Credentials{
		Username:            split[0],
		PersonalAccessToken: split[1],
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
