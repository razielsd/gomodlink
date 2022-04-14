package giturl

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

type GitRepo struct {
	ModuleURL string
}

func Parse(s string) (*GitRepo, error) {
	s = cleanURL(s)
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(u.Host)
	repo := &GitRepo{}
	repo.ModuleURL = strings.TrimSuffix(fmt.Sprintf("%s%s", host, u.Path), ".git")
	return repo, nil
}

func cleanURL(u string) string {
	if i := strings.Index(u, "://"); i > 0 {
		if ai := strings.Index(u, "@"); ai > 0 {
			return u[0:i+3] + u[ai+1:]
		}
		return u[i+3:]
	}
	if i := strings.Index(u, "@"); i > 0 {
		return strings.Replace(u[i+1:], ":", "/", 1)
	}
	return u
}
