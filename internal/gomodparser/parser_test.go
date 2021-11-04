package gomodparser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewParser(t *testing.T) {
	require.NotNil(t, NewParser())
}

func TestParser_Parse(t *testing.T) {

	gomodfile := `module github.com/razielsd/gomodlink

go 1.17

require (
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
`
	parser := NewParser()
	err := parser.Parse(gomodfile)
	require.NoError(t, err)

	deps := []string{
		"github.com/stretchr/testify",
		"golang.org/x/sync",
		"github.com/davecgh/go-spew",
		"github.com/pmezard/go-difflib",
		"gopkg.in/yaml.v3",
	}
	require.Equal(t, deps, parser.GetDeps())
}
