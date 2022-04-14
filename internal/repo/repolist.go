package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/sync/errgroup"
)

type RepoList struct { //nolint:revive
	repoList []*Repository
}

type RepoJSON struct { //nolint:revive
	Repo []struct {
		Branch string `json:"branch"`
		URL    string `json:"url"`
		Name   string `json:"name"`
	} `json:"repo"`
}

func (l *RepoList) AddRepository(repo *Repository) {
	repo.Key = fmt.Sprintf("n%d", len(l.repoList))
	l.repoList = append(l.repoList, repo)
}

func (l *RepoList) GetAll() []*Repository {
	return l.repoList
}

func (l *RepoList) Load() error {
	g := &errgroup.Group{}
	for _, repo := range l.repoList {
		repo := repo
		g.Go(func() error {
			return repo.Load()
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (l *RepoList) Intersect() (map[string][]*Repository, error) {
	result := make(map[string][]*Repository)
	for _, repo := range l.repoList {
		intersect, err := repo.Intersect(l.repoList)
		if err != nil {
			return nil, err
		}
		result[repo.Key] = intersect
	}
	return result, nil
}

func (l *RepoList) LoadFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data := &RepoJSON{}
	err = json.Unmarshal(content, data)
	if err != nil {
		return err
	}
	for _, v := range data.Repo {
		if v.Name == "" {
			v.Name = MakeDefaultRepoName(v.URL)
		}
		l.AddRepository(NewRepository(v.URL, v.Branch, v.Name))
	}
	return nil
}

func MakeDefaultRepoName(u string) string {
	parts := strings.Split(u, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}
