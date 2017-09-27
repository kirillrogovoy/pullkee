<p align="center">
  <img src="https://maxcdn.icons8.com/Android_L/PNG/512/Programming/pull_request-512.png" width="200" alt="Pull request icon">
  <h3 align="center">pullkee</h3>
  <p align="center">A simple Pull Requests analyzer.</p>
  <p align="center">
    <a href="https://github.com/kirillrogovoy/pullkee/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/kirillrogovoy/pullkee.svg?style=flat-square"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/kirillrogovoy/pullkee"><img alt="Travis" src="https://img.shields.io/travis/kirillrogovoy/pullkee.svg?style=flat-square"></a>
    <a href="https://codecov.io/gh/kirillrogovoy/pullkee"><img alt="Codecov branch" src="https://img.shields.io/codecov/c/github/kirillrogovoy/pullkee/master.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/kirillrogovoy/pullkee"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kirillrogovoy/pullkee?style=flat-square"></a>
    <a href="http://godoc.org/github.com/kirillrogovoy/pullkee"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
    <a href="https://github.com/goreleaser"><img alt="Powered By: GoReleaser" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square"></a>
  </p>
</p>

## Why?

It's always been fun for me to browse pages like this: https://github.com/facebook/react/graphs/contributors.

Although it can't possibly give one a meaningful insight, I've been curious about a number of other metrics in my work project.
Who is producing more code? Who's being picked as a reviewer more ofter? How long does it take for us on average to merge
a pull request? Who writes more (or longer) comments?

Again, it's not something you can strongly base your decisions on, but it's just plain curiosity.
Also, maybe a combination of such metrics could actually mean something.

**So,** the single purpose of this project is to provide that kind of insights given a Github repository name.

Another great motivator for me was to learn Golang as this project presents a big deal of different challenges.

## Install

If you have the Golang environment set up on your computer, just run:
```
go get github.com/kirillrogovoy/pullkee
```
and you are all set.

Otherwise, you can manually download the binary from the [Releases page](https://github.com/kirillrogovoy/pullkee/releases).
In order to install it automatically, run:
```
curl https://raw.githubusercontent.com/kirillrogovoy/pullkee/master/godownloader.sh | bash
```

## Usage

Just run `pullkee` to see the usage. Here's a copy for convenience:
```
Usage:
    pullkee [flags] [repo]
    repo - Github repository path as "username/reponame"

    Flags:
    --limit - Only use N last pull requests

    Environment variables:
    GITHUB_CREDS - API credentials in the format "username:personal_access_token"
```

For example, to get the reports for the last 500 merged pull requests of the React repo, run this:
```sh
GITHUB_CREDS="your_name:your_key" pullkee --limit 500 facebook/react
```

## API rate limits and cache

Strongly consider using the `--limit` parameter on big repos since
you have a limited number of requests to make to the Github API. For me, it's currently 5000 per 1 hour.
Also, always provide the `GITHUB_CREDS` env var, otherwise you only have 60 requests per 1 hour without it.

Don't have a token yet? [Say no more](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/).

That said, pullkee always uses a per-PR local cache in order to avoid
repetitive requests for the data of the same pull request.

It means, even if you ran out of requests, you still can wait for them to renew and continue.

## Metrics

The current list of metrics is "baked in" into the project and cannot be changed from outside.
I'd prefer to keep it that way unless someone is explicitely interested in that.

Just fork the repo to change or add metrics.

## Contribute

Please, contribute in any way if you feel like it.
Start from the [docs](https://godoc.org/github.com/kirillrogovoy/pullkee) to get a high-level overview of the code.
Let me know if you can't do something. **Keep the test coverage > 95%**.
