package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

type RepoInfo struct {
	Url  string
	Hash string
}

func downloadFile(uri string) (data []byte) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to get URL %s, please make sure that the URL is correct", err)
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read Body %s, please make sure that the URL is correct", err)
	}
	return
}

func readData(data []byte) (repos []RepoInfo) {
	var repo RepoInfo
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		temp := strings.SplitN(scanner.Text(), " ", 2)
		repo.Url = temp[0]
		repo.Hash = temp[1]
		repos = append(repos, repo)
		// fmt.Printf("RepoURL: %s, RepoHash: %s\n", repo.Url, repo.Hash)
	}
	return
}

func readRepo(path string, hash string) (files []object.File) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: path,
	})
	if err != nil {
		log.Fatalf("Sorry, but we haven't be able to open the repository %s", err)
	}
	commit, err := repo.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		log.Fatalf("Sorry, but we haven't be able to read the commit %s", err)
	}
	// fmt.Printf(commit.String())
	// for _, file := range findDokerfiles(commit.Tree()) {
	// 	fmt.Printf("File: %s\n", file.Name)
	// }
	files = findDokerfiles(commit.Tree())
	return
}

func findDokerfiles(tree *object.Tree, err error) (files []object.File) {
	tree.Files().ForEach(func(f *object.File) error {
		match, err := regexp.MatchString(`(?:^|\W)Dockerfile$`, f.Name)
		if f.Mode.IsFile() && match {
			files = append(files, *f)
		}
		return err
	})
	return
}

func readFile(file object.File) (from []string) {
	from = append(from, "from string")
	// log.Fatalf("Sorry, we haven't be able to complete this implementation yet, keep tuned")
	return
}

func defaultImplementation(url *string) {
	jsonOutput := "{\n  \"data\": {\n"
	file := downloadFile(*url)
	repos := readData(file)
	for i, element := range repos {
		jsonOutput = jsonOutput + "    \"" + element.Url + ":" + element.Hash + "\": {\n"
		dockerfiles := readRepo(element.Url, element.Hash)
		for j, file := range dockerfiles {
			jsonOutput = jsonOutput + "      \"" + file.Name + "\": [\n"
			fromStrings := readFile(file)
			for k, imageFrom := range fromStrings {
				if k < len(fromStrings)-1 {
					jsonOutput = jsonOutput + "        \"" + imageFrom + "\",\n"
				} else {
					jsonOutput = jsonOutput + "        \"" + imageFrom + "\"\n"
				}
			}
			if j < len(dockerfiles)-1 {
				jsonOutput = jsonOutput + "      ],\n"
			} else {
				jsonOutput = jsonOutput + "      ]\n"
			}
		}
		if i < len(repos)-1 {
			jsonOutput = jsonOutput + "    },\n"
		} else {
			jsonOutput = jsonOutput + "    }\n"
		}
	}
	jsonOutput = jsonOutput + "  }\n}"
	fmt.Printf(jsonOutput)
}

func fixedImplementation(url *string) {
	log.Fatalf("Sorry, we haven't be able to complete this implementation yet, keep tuned")
}

func main() {
	url := flag.String("url", "", "You need put the URL from download the file")
	fix := flag.Bool("fix", false, "By default it's false, but if you want to see the implementation with teh correct JSON change at true")
	flag.Parse()

	if *url != "" && *fix {
		fixedImplementation(url)
	} else if *url != "" {
		defaultImplementation(url)
	} else {
		flag.Usage()
	}
}
