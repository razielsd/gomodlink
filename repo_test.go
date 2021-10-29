package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository_GetModuleUrl(t *testing.T) {
	repo := NewRepository("ssh://git@gitlab.ourrepo.dev:21722/dbg/item-catalog.git", "develop")

	modUrl, err := repo.GetModuleUrl()
	require.NoError(t, err)
	require.Equal(t, "gitlab.ourrepo.dev/dbg/item-catalog", modUrl)
}

func TestRepository_Intersect(t *testing.T) {
	repo := NewRepository("ssh://git@gitlab.ourrepo.dev:21722/dbg/item-catalog.git", "develop")
	repo.Linked = []string{
		"1234",
		"2345/1234",
		"34567/999",
		"gitlab.ourrepo.dev/dbg/item-catalog",
		"gitlab.ourrepo.dev/dbg/item-catalog/pkg",
	}

	search := []string{"1234", "34567"}
	expected := []string{
		"1234",
		"34567/999",
	}
	result, err := repo.Intersect(search)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}
