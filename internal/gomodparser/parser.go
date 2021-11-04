package gomodparser

import (
	"bufio"
	"strings"
)

type Parser struct {
	deps []string
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(modFile string) error {
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
			deps := strings.Fields(line)[0]
			p.deps = append(p.deps, deps)
		}
	}

	return nil
}

func (p *Parser) GetDeps() []string {
	return p.deps
}

