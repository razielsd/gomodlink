package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type Repository struct {
	Url    string
	Branch string
	Linked []string
}

func NewRepository(url, branch string) *Repository {
	return &Repository{
		Url:    url,
		Branch: branch,
		Linked: make([]string, 0),
	}
}

func (r *Repository) Load() error {
	modFile, err := r.download()
	if err != nil {
		return err
	}
	err = r.parseGoMod(modFile)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetModuleUrl() (string, error) {
	u, err := url.Parse(r.Url)
	if err != nil {
		return "", err
	}
	host, _, _ := net.SplitHostPort(u.Host)
	mod := strings.TrimSuffix(fmt.Sprintf("%s%s", host, u.Path), ".git")
	return mod, nil
}

func (r *Repository) Intersect(l []string) ([]string, error) {
	result := make([]string, 0)
	modUrl, err := r.GetModuleUrl()
	if err != nil {
		return nil, err
	}
	for _, s := range l {
		for _, v := range r.Linked {
			if (strings.HasPrefix(v, s+"/") || (v == s)) && !strings.HasPrefix(v, modUrl) {
				result = append(result, v)
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

func (r *Repository) parseGoMod(modFile string) error {
	scanner := bufio.NewScanner(strings.NewReader(modFile))
	isReq := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "require") {
			isReq = true
			continue
		} else if line == ")" {
			isReq = false
		}
		if isReq {
			repoLink := strings.Fields(line)[0]
			r.Linked = append(r.Linked, repoLink)
		}
	}
	return nil
}
