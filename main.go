package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
)

var from string

func init() {
	flag.StringVar(&from, "from", "", "source file to read from")
}

func showHelp() {
	fmt.Println("Params:")
	fmt.Println("  from - repository configuration, see example.json")
	fmt.Println("Example: gomodlink --from example.json")
	fmt.Println()
}

func main() {
	flag.Parse()
	if from == "" {
		fmt.Println("Error: require param: from")
		showHelp()
		os.Exit(1)
	}
	if _, err := os.Stat(from); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File not found: %s\n", from)
		os.Exit(1)
	}

	repoList := RepoList{}
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

	out, _ := repoList.Intersect()

	keyList := make([]string, 0)
	for i := range out {
		keyList = append(keyList, i)
	}

	sort.Slice(keyList, func(i, j int) bool {
		fr := out[keyList[i]]
		sr := out[keyList[j]]
		return len(fr) > len(sr)
	})

	depCounter := 0
	for _, name := range keyList {
		v := out[name]
		fmt.Printf("Repository: %s (%d)\n", name, len(v))
		for _, line := range v {
			depCounter++
			fmt.Printf("      %s\n", line)
		}
		fmt.Println()
	}
	avg := 0.0
	if len(out) > 1 {
		avg = float64(depCounter) / float64(len(out) - 1)
	}

	fmt.Printf("Total repository: %d\n", len(out))
	fmt.Printf("Total dependencies: %d\n", depCounter)
	fmt.Printf("AVG dependencies: %.2f\n", avg)
}
