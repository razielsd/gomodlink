package repo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/razielsd/gomodlink/internal/giturl"
	"github.com/razielsd/gomodlink/internal/gomodparser"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Repository struct {
	Url    string
	Branch string
	Deps   []string
	Key    string
	Name   string
	modUrl string
}

func NewRepository(url, branch, name string) *Repository {
	b := make([]byte, 4)
	rand.Read(b)
	return &Repository{
		Url:    url,
		Branch: branch,
		Name:   name,
		Deps:   make([]string, 0),
		modUrl: "",
		Key:    hex.EncodeToString(b),
	}
}

func (r *Repository) Load() error {
	modFile, err := r.download()
	if err != nil {
		return err
	}
	parser := gomodparser.NewParser()
	err = parser.Parse(modFile)
	if err != nil {
		return err
	}
	r.Deps = parser.GetDeps()
	return nil
}

func (r *Repository) GetModuleUrl() (string, error) {
	if r.modUrl != "" {
		return r.modUrl, nil
	}
	u, err := giturl.Parse(r.Url)
	if err != nil {
		return "", err
	}
	return u.ModuleUrl, nil
}

func (r *Repository) Intersect(l []*Repository) ([]*Repository, error) {
	// check build all moduleUrl
	for _, repo := range l {
		_, err := repo.GetModuleUrl()
		if err != nil {
			return nil, err
		}
	}

	result := make([]*Repository, 0)
	modUrl, err := r.GetModuleUrl()
	if err != nil {
		return nil, err
	}
	for _, repo := range l {
		s, _ := repo.GetModuleUrl()
		for _, v := range r.Deps {
			if (strings.HasPrefix(v, s+"/") || (v == s)) && !strings.HasPrefix(v, modUrl) {
				result = append(result, repo)
				continue
			}
		}
	}
	return result, nil
}

func (r *Repository) download() (string, error) {
	cmdArg := strings.Fields(
		fmt.Sprintf("archive --remote=%s %s go.mod", r.Url, r.Branch),
	)
	gitCmd := exec.Command("git", cmdArg...)
	tarCmd := exec.Command("tar", "-xO")
	pipe, _ := gitCmd.StdoutPipe()
	tarCmd.Stdin = pipe
	tarCmd.Stdout = os.Stdout

	var out bytes.Buffer
	tarCmd.Stdout = &out

	err := gitCmd.Start()
	if err != nil {
		return "", err
	}
	err = tarCmd.Start()
	if err != nil {
		return "", err
	}
	err = gitCmd.Wait()
	if err != nil {
		return "", err
	}
	err = tarCmd.Wait()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
