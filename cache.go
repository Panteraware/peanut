package peanut

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
	"time"
)

func cacheReleaseList(repo Repo) string {
	headers := []CustomHeaders{{Label: "Accept", Value: "application/vnd.github.preview"}}

	if stringLength(repo.Token) > 0 {
		authHeader := CustomHeaders{Label: "Authorization", Value: fmt.Sprintf("Token %s", repo.Token)}
		headers = append(headers, authHeader)
	}

	res := getRequest(repo.Url, headers)

	if res.StatusCode != 200 {
		_ = errors.New(fmt.Sprintf("Tried to cache RELEASES, but failed fetching %s, status %s", repo.Url, res.Status))
	}

	buf := new(strings.Builder)

	_, err := io.Copy(buf, res.Body)

	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	}

	content := buf.String()

	// r := regexp.MustCompile(`/[^ ]*\.nupkg/gim`)

	// matches := r.FindAllStringSubmatch(content, -1)

	// for _, v := range matches {
	// 	const nuPKG = strings.Replace(url, "RELEASES", v[1], 1)
	// 	content = strings.Replace(content, v[1], nuPKG, -1)
	// }

	return content
}

func refreshCache(repo Repo) bool {
	client := github.NewClient(nil)

	if len(repo.Token) > 0 {
		client.WithAuthToken(repo.Token)
	}

	latest, _, err := client.Repositories.GetLatestRelease(context.Background(), repo.Owner, repo.Name)
	if err != nil {
		log.Error().Err(err).Str("owner", repo.Owner).Str("name", repo.Name).Msg("error getting latest release")
	}

	println(latest.Name)

	// repo := fmt.Sprintf("%s/%s", util.Config.Owner, util.Config.Repo)
	// url := fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page_100", repo)
	// println(url)
	// headers := []CustomHeaders{{Label: "Accept", Value: "application/vnd.github.preview"}}

	// if stringLength(util.Config.GitToken) > 0 {
	// 	authHeader := CustomHeaders{Label: "Authorization", Value: fmt.Sprintf("Token %s", util.Config.GitToken)}
	// 	headers = append(headers, authHeader)
	// }

	// res := getRequest(url, headers)

	// if res.StatusCode != 200 {
	// 	_ = errors.New(fmt.Sprintf("Github API responded with %s for url %s", res.Status, url))
	// }

	// var data

	// getJson(res.Body, data)

	// println(data)

	repo.Cache.LastUpdate = time.Now()

	return true
}

func isOutdated(cache Cache) bool {
	if time.Now().UnixMilli()-cache.LastUpdate.UnixMilli() > Config.Interval {
		return true
	}

	return false
}

func loadCache(repo Repo) *github.RepositoryRelease {
	if repo.Cache.LastUpdate == time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC) || isOutdated(repo.Cache) {
		refreshCache(repo)
	}

	return &repo.Cache.Latest
}
