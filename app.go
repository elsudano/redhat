package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Data contain all the information recovered of the repositories
type Data struct {
	data struct {
		repository struct {
		}
	}
	Path string `json:"path"`
}

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
		fmt.Printf("RepoURL: %s, RepoHash: %s\n", repo.Url, repo.Hash)
	}
	return
}

func readRepo(path string, hash string) {
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
	fmt.Printf(commit.String())
	findDokerfile(commit.Tree())
}

func findDokerfile(path_in *object.Tree, err error) (path_out string) {
	for _, entry := range path_in.Entries {
		if !entry.Mode.IsFile() {
			path_out = entry.Name + "/" + findDokerfile(path_in.Tree(entry.Name))
			fmt.Printf("Dir: %s\n", path_out)
		} else {
			path_out = entry.Name
			fmt.Printf("File: %s\n", path_out)
		}
	}
	return
}

func main() {
	url := flag.String("url", "", "You need put the URL from download the file")
	flag.Parse()

	if *url != "" {
		file := downloadFile(*url)
		repos := readData(file)
		for _, element := range repos {
			readRepo(element.Url, element.Hash)
		}
	} else {
		flag.Usage()
	}
}
