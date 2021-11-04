package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository_GetModuleUrl(t *testing.T) {
	repo := NewRepository("ssh://git@gitlab.ourrepo.dev:21722/dbg/item-catalog.git", "develop", "item")

	modUrl, err := repo.GetModuleUrl()
	require.NoError(t, err)
	require.Equal(t, "gitlab.ourrepo.dev/dbg/item-catalog", modUrl)
}

func TestRepository_Intersect(t *testing.T) {
	repo := NewRepository("ssh://git@gitlab.ourrepo.dev:21722/dbg/item-catalog.git", "develop", "item")
	repo.Deps = []string{
		"1234",
		"2345/1234",
		"34567/999",
		"gitlab.ourrepo.dev/dbg/item-catalog",
		"gitlab.ourrepo.dev/dbg/item-catalog/pkg",
	}
	search := []*Repository{
		NewRepository("1234", "develop", "1234"),
		NewRepository("34567", "develop", "34567"),
		NewRepository("94567", "develop", "94567"),
	}
	expected := []*Repository{
		search[0],
		search[1],
	}
	result, err := repo.Intersect(search)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestRepoList_Load(t *testing.T) {

}