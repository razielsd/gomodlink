package repo

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
)

type ReportWriter struct {
	Repolist *RepoList
}

var (
	ErrCreateTmpFile = errors.New("unable to create tmp file")
	ErrWriteTmpFile  = errors.New("unable write to tmp file")
)

func (r *ReportWriter) BuildText(in map[string][]*Repository) string {
	out := bytes.NewBufferString("")

	keyList := make([]string, 0)
	for i := range in {
		keyList = append(keyList, i)
	}
	mapName := make(map[string]string)
	for _, repo := range r.Repolist.GetAll() {
		mapName[repo.Key] = repo.Name
	}

	sort.Slice(keyList, func(i, j int) bool {
		fr := in[keyList[i]]
		sr := in[keyList[j]]
		return len(fr) > len(sr)
	})

	depCounter := 0
	for _, name := range keyList {
		v := in[name]
		_, _ = fmt.Fprintf(out, "Repository: %s (%d)\n", mapName[name], len(v))
		for _, line := range v {
			depCounter++
			_, _ = fmt.Fprintf(out, "      %s\n", line.Name)
		}
		_, _ = fmt.Fprintf(out, "\n")
	}
	avg := 0.0
	if len(in) > 1 {
		avg = float64(depCounter) / float64(len(in)-1)
	}

	_, _ = fmt.Fprintf(out, "Total repository: %d\n", len(in))
	_, _ = fmt.Fprintf(out, "Total dependencies: %d\n", depCounter)
	_, _ = fmt.Fprintf(out, "AVG dependencies: %.2f\n", avg)

	return out.String()
}

func (r *ReportWriter) BuildGraphviz(in map[string][]*Repository, filename string, format string) error {
	out := bytes.NewBufferString("digraph gomodlink { \n")

	for _, repo := range r.Repolist.GetAll() {
		_, _ = fmt.Fprintf(out,
			"%s [label=\"%s\", URL=\"%s\", tooltip=\"%s\"];\n", repo.Key, repo.Name, repo.Url, repo.Name,
		)
	}

	for k, v := range in {
		for _, repo := range v {
			_, _ = fmt.Fprintf(out, "%s -> %s;\n", k, repo.Key)
		}
	}

	_, _ = fmt.Fprintf(out, "}")

	destFile, err := ioutil.TempFile("/tmp", "gml-")
	if err != nil {
		return ErrCreateTmpFile
	}

	if _, err := destFile.Write(out.Bytes()); err != nil {
		return ErrWriteTmpFile
	}

	destFile.Close()

	cmdArg := strings.Fields(
		fmt.Sprintf("-T%s %s -o %s", format, destFile.Name(), filename),
	)
	cmd := exec.Command("dot", cmdArg...)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}
