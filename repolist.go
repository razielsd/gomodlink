package main

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/sync/errgroup"
)

type RepoList struct {
	repoList []*Repository
}

type RepoJson struct {
	Repo []struct {
		Branch string `json:"branch"`
		URL    string `json:"url"`
	} `json:"repo"`
}

func (l *RepoList) AddRepository(repo *Repository) {
	l.repoList = append(l.repoList, repo)
}

func (l *RepoList) Load() error {
	g, _ := errgroup.WithContext(context.Background())
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

func (l *RepoList) Intersect() (map[string][]string, error) {
	result := make(map[string][]string)
	urlList := make([]string, 0)
	for _, repo := range l.repoList {
		link, err := repo.GetModuleUrl()
		if err != nil {
			return nil, err
		}
		urlList = append(urlList, link)
	}
	for _, repo := range l.repoList {
		modUrl, err := repo.GetModuleUrl()
		if err != nil {
			return nil, err
		}
		result[modUrl], err = repo.Intersect(urlList)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (l *RepoList) LoadFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data := &RepoJson{}
	err = json.Unmarshal(content, data)
	if err != nil {
		return err
	}
	for _, v := range data.Repo {
		l.AddRepository(NewRepository(v.URL, v.Branch))
	}
	return nil
}
