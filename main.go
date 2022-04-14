package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/razielsd/gomodlink/internal/repo"
)

const (
	formatTxt string = "txt"
	formatSvg string = "svg"
	formatPng string = "png"
)

var (
	from, format, outfile string
	openFile              bool
)

const reportFilePerm = 0o600

func init() {
	flag.StringVar(&from, "from", "", "source file to read from")
	flag.StringVar(&format, "format", "txt", "output format [txt, svg]")
	flag.StringVar(&outfile, "out", "", "output filename")
	flag.BoolVar(&openFile, "open", false, "open image after run, support only macos")
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	if from == "" {
		fmt.Println("Error: require param: from")
		flag.Usage()
		os.Exit(1)
	}
	if _, err := os.Stat(from); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File not found: %s\n", from)
		os.Exit(1)
	}

	repoList := repo.RepoList{}
	err := repoList.LoadFromFile(from)
	if err != nil {
		fmt.Printf("Error load repository configuration: %s\n", err.Error())
		os.Exit(1)
	}
	err = repoList.Load()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	writer := &repo.ReportWriter{
		Repolist: &repoList,
	}
	in, _ := repoList.Intersect()

	switch format {
	case formatTxt:
		out := writer.BuildText(in)
		if outfile == "" {
			fmt.Println(out)
		} else if err := ioutil.WriteFile(outfile, []byte(out), reportFilePerm); err != nil {
			fmt.Printf("ERROR: unable write report - %s\n", err.Error())
		}
	case formatSvg, formatPng:
		if err := writer.BuildGraphviz(in, outfile, format); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		if openFile {
			openbrowser(outfile)
		}
	}
}
