package statuscmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/robertguss/mycelium/internal/execrun"
	"github.com/robertguss/mycelium/internal/manifest"
)

// Remote is one GitHub idea-topic repo after optional manifest fetch.
type Remote struct {
	Name       string
	Owner      string
	IsArchived bool
	State      string
	Tier       string
	Revisit    string
}

type ghRepoJSON struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	IsArchived bool   `json:"isArchived"`
	Owner      struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type ghRepoListJSON struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	IsArchived bool   `json:"isArchived"`
}

type ghTopicsJSON struct {
	Names []string `json:"names"`
}

func fetchRemotes(ctx context.Context, runner execrun.Runner) ([]Remote, error) {
	login, err := ghLogin(ctx, runner)
	if err != nil {
		return nil, err
	}
	repos, err := listIdeaRepos(ctx, runner, login)
	if err != nil {
		return nil, err
	}
	out := make([]Remote, 0, len(repos))
	for _, r := range repos {
		rem := Remote{
			Name:       r.Name,
			Owner:      r.Owner,
			IsArchived: r.IsArchived,
		}
		state, tier, revisit, ok := fetchRemoteManifest(ctx, runner, r.Owner, r.Name)
		if !ok {
			rem.State = "unread"
			rem.Tier = "-"
			rem.Revisit = "-"
		} else {
			rem.State = state
			rem.Tier = tier
			rem.Revisit = revisit
		}
		out = append(out, rem)
	}
	return out, nil
}

func ghLogin(ctx context.Context, runner execrun.Runner) (string, error) {
	res, err := runner.Run(ctx, "gh", []string{"api", "user", "--jq", ".login"}, execrun.RunOpts{})
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("gh api user: exit %d", res.ExitCode)
	}
	login := strings.TrimSpace(string(res.Stdout))
	if login == "" {
		return "", fmt.Errorf("gh api user: empty login")
	}
	return login, nil
}

func listIdeaRepos(ctx context.Context, runner execrun.Runner, login string) ([]Remote, error) {
	repos, err := runRepoSearch(ctx, runner, []string{
		"search", "repos", "topic:idea", "user:" + login,
		"--json", "name,url,isArchived,owner",
	})
	if err == nil {
		return repos, nil
	}
	repos, err = runRepoSearch(ctx, runner, []string{
		"search", "repos", "--owner", login, "--topic", "idea",
		"--json", "name,url,isArchived,owner",
	})
	if err == nil {
		return repos, nil
	}
	return listReposViaTopics(ctx, runner, login)
}

func runRepoSearch(ctx context.Context, runner execrun.Runner, args []string) ([]Remote, error) {
	res, err := runner.Run(ctx, "gh", args, execrun.RunOpts{})
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("gh search: exit %d", res.ExitCode)
	}
	return parseSearchRepos(res.Stdout)
}

func parseSearchRepos(raw []byte) ([]Remote, error) {
	var rows []ghRepoJSON
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	out := make([]Remote, 0, len(rows))
	for _, r := range rows {
		owner := strings.TrimSpace(r.Owner.Login)
		if owner == "" {
			owner = ownerFromURL(r.URL)
		}
		if r.Name == "" || owner == "" {
			continue
		}
		out = append(out, Remote{
			Name:       r.Name,
			Owner:      owner,
			IsArchived: r.IsArchived,
		})
	}
	return out, nil
}

func listReposViaTopics(ctx context.Context, runner execrun.Runner, login string) ([]Remote, error) {
	res, err := runner.Run(ctx, "gh", []string{
		"repo", "list", login, "--limit", "1000",
		"--json", "name,url,isArchived",
	}, execrun.RunOpts{})
	if err != nil {
		return nil, err
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("gh repo list: exit %d", res.ExitCode)
	}
	var rows []ghRepoListJSON
	if err := json.Unmarshal(res.Stdout, &rows); err != nil {
		return nil, err
	}
	out := make([]Remote, 0, len(rows))
	for _, r := range rows {
		if r.Name == "" {
			continue
		}
		owner := login
		if o := ownerFromURL(r.URL); o != "" {
			owner = o
		}
		ok, topicErr := repoHasIdeaTopic(ctx, runner, owner, r.Name)
		if topicErr != nil {
			return nil, topicErr
		}
		if !ok {
			continue
		}
		out = append(out, Remote{
			Name:       r.Name,
			Owner:      owner,
			IsArchived: r.IsArchived,
		})
	}
	return out, nil
}

func repoHasIdeaTopic(ctx context.Context, runner execrun.Runner, owner, name string) (bool, error) {
	res, err := runner.Run(ctx, "gh", []string{
		"api", "repos/" + owner + "/" + name + "/topics",
	}, execrun.RunOpts{})
	if err != nil {
		return false, err
	}
	if res.ExitCode != 0 {
		return false, fmt.Errorf("gh api topics: exit %d", res.ExitCode)
	}
	var topics ghTopicsJSON
	if err := json.Unmarshal(res.Stdout, &topics); err != nil {
		return false, err
	}
	for _, n := range topics.Names {
		if n == "idea" {
			return true, nil
		}
	}
	return false, nil
}

func fetchRemoteManifest(ctx context.Context, runner execrun.Runner, owner, name string) (state, tier, revisit string, ok bool) {
	res, err := runner.Run(ctx, "gh", []string{
		"api", "repos/" + owner + "/" + name + "/contents/mycelium.toml",
		"--jq", ".content",
	}, execrun.RunOpts{})
	if err != nil || res.ExitCode != 0 {
		return "", "", "", false
	}
	b64 := strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == ' ' || r == '\t' {
			return -1
		}
		return r
	}, string(res.Stdout))
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", "", "", false
	}
	m, err := manifest.Parse(decoded)
	if err != nil {
		return "", "", "", false
	}
	return m.State, m.Tier, m.Revisit, true
}

func ownerFromURL(url string) string {
	u := strings.TrimSpace(url)
	u = strings.TrimSuffix(u, ".git")
	u = strings.TrimSuffix(u, "/")
	const prefix = "https://github.com/"
	if !strings.HasPrefix(u, prefix) {
		return ""
	}
	rest := strings.TrimPrefix(u, prefix)
	parts := strings.Split(rest, "/")
	if len(parts) >= 1 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

// MergeIdeas joins local scan rows with remote idea repos (§9.4).
// Merge key = local slug vs remote repo name.
func MergeIdeas(local []Idea, remotes []Remote, includeArchived bool) []Idea {
	bySlug := make(map[string]Idea, len(local)+len(remotes))
	for _, loc := range local {
		bySlug[loc.Slug] = loc
	}

	for _, rem := range remotes {
		ownerName := rem.Owner + "/" + rem.Name
		if loc, ok := bySlug[rem.Name]; ok {
			if (loc.State == "archived" || rem.IsArchived) && !includeArchived {
				delete(bySlug, rem.Name)
				continue
			}
			loc.Flag = "ok"
			loc.Github = ownerName
			bySlug[rem.Name] = loc
			continue
		}
		if rem.IsArchived && !includeArchived {
			continue
		}
		bySlug[rem.Name] = Idea{
			Slug:    rem.Name,
			State:   rem.State,
			Tier:    rem.Tier,
			Revisit: rem.Revisit,
			Flag:    "remote",
			Github:  ownerName,
		}
	}

	out := make([]Idea, 0, len(bySlug))
	for _, idea := range bySlug {
		if idea.Flag == "unpublished" && idea.State == "archived" && !includeArchived {
			continue
		}
		out = append(out, idea)
	}
	return out
}
