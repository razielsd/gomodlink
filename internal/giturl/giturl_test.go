package giturl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGitUrl_Parse(t *testing.T) {
	tests := []struct {
		name   string
		giturl string
		modurl string
	}{
		{
			name: "git",
			giturl: "git@git.ourrepo.tech:path/to/module//customer-manager.git",
			modurl: "git.ourrepo.tech/path/to/module//customer-manager",
		},
		{
			name: "https",
			giturl: "https://git.ourrepo.tech/path/to/module//customer-manager.git",
			modurl: "git.ourrepo.tech/path/to/module//customer-manager",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo, err := Parse(tt.giturl)
			require.NoError(t, err)
			require.NotNil(t, repo)
			require.Equal(t, tt.modurl, repo.ModuleUrl)
		})
	}
}