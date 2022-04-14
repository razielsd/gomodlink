package repo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/razielsd/gomodlink/internal/giturl"
	"github.com/razielsd/gomodlink/internal/gomodparser"
)

const tryDownload = 5

type Repository struct {
	URL    string
	Branch string
	Deps   []string
	Key    string
	Name   string
	modURL string
}

func NewRepository(url, branch, name string) *Repository {
	const parts = 4
	b := make([]byte, parts)
	rand.Read(b) //nolint:gosec
	return &Repository{
		URL:    url,
		Branch: branch,
		Name:   name,
		Deps:   make([]string, 0),
		modURL: "",
		Key:    hex.EncodeToString(b),
	}
}

func (r *Repository) Load() error {
	modFile, err := r.download()
	if err != nil {
		return err
	}
	parser := gomodparser.NewParser()
	if err := parser.Parse(modFile); err != nil {
		return err
	}
	r.Deps = parser.GetDeps()
	return nil
}

func (r *Repository) GetModuleURL() (string, error) {
	if r.modURL != "" {
		return r.modURL, nil
	}
	u, err := giturl.Parse(r.URL)
	if err != nil {
		return "", err
	}
	return u.ModuleURL, nil
}

func (r *Repository) Intersect(l []*Repository) ([]*Repository, error) {
	// check build all moduleUrl
	for _, repo := range l {
		_, err := repo.GetModuleURL()
		if err != nil {
			return nil, err
		}
	}

	result := make([]*Repository, 0)
	modURL, err := r.GetModuleURL()
	if err != nil {
		return nil, err
	}
	for _, repo := range l {
		s, _ := repo.GetModuleURL()
		for _, v := range r.Deps {
			if (strings.HasPrefix(v, s+"/") || (v == s)) && !strings.HasPrefix(v, modURL) {
				result = append(result, repo)
				continue
			}
		}
	}
	return result, nil
}

func (r *Repository) download() (string, error) {
	var err error
	var modFile string
	for i := tryDownload; i > 0; i-- {
		modFile, err = r.tryDownload()
		if err == nil {
			return modFile, nil
		}
	}
	return "", fmt.Errorf("error download go.mod(%s)", r.URL)
}

func (r *Repository) tryDownload() (string, error) {
	cmdArg := strings.Fields(
		fmt.Sprintf("archive --remote=%s %s go.mod", r.URL, r.Branch),
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
