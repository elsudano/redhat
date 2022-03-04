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
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Data contain all the information recovered of the repositories
type Data struct {
	data struct {
		repository struct {
		}
	}
	Path         string `json:"path"`
	Denomination string `json:"displayName"`
	Description  string `json:"annotation"`
	// Image        string `json:"image"`
	CPU struct {
		Processors int `json:"processors"`
	}
	PowerStatus string `json:"power_state"`
	Memory      int    `json:"memory"`
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

func readRepo(path string) {
	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// Length of the HEAD history
	Info("git rev-list HEAD --count")

	// ... retrieving the HEAD reference
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits
	var cCount int
	err = cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})
	CheckIfError(err)

	fmt.Println(cCount)
}

func main() {
	url := flag.String("url", "", "You need put the URL from download the file")
	flag.Parse()

	if *url != "" {
		file := downloadFile(*url)
		repos := readData(file)
		for _, element := range repos {
			readRepo(element.Url)
		}
	} else {
		flag.Usage()
	}
}
