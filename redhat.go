/*
	Copyright 2022 Carlos de la Torre. All rights reserved.
	Use of this source code is governed by a MIT License
	license that can be found in the LICENSE file.
*/
package redhat

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Struct that represent a complete JSON wrapper, in order to read this
// data later
type jsonWrapper struct {
	Data data `json:"data"`
}

type data struct {
	Repositories []repository `json:"repos"`
}

type repository struct {
	Url         string       `json:"url"`
	Hash        string       `json:"hash"`
	Dockerfiles []dockerfile `json:"dockerfile"`
}

type dockerfile struct {
	Pathfile string   `json:"path"`
	Froms    []string `json:"from"`
}

type RepoInfo struct {
	Url  string
	Hash string
}

// This function is in charge of download a file and put the information in memory
// it put all the data in a []byte array in order to read easly.
func Downloadfile(uri string) (data []byte) {
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

// This
func Readdata(data []byte) (repos []RepoInfo) {
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

func Readrepo(path string, hash string) (files []object.File) {
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
	files = Finddokerfiles(commit.Tree())
	return
}

func Finddokerfiles(tree *object.Tree, err error) (files []object.File) {
	tree.Files().ForEach(func(f *object.File) error {
		match, err := regexp.MatchString(`(?:^|\W)Dockerfile$`, f.Name)
		if f.Mode.IsFile() && match {
			files = append(files, *f)
		}
		return err
	})
	return
}

func Readfile(file object.File) (from []string) {
	lines, err := file.Lines()
	if err != nil {
		log.Fatalf("Sorry, but we haven't be able to read the file %s", err)
	}
	for _, line := range lines {
		match, _ := regexp.MatchString(`^FROM`, line)
		if match {
			line = strings.Split(line, " ")[1]
			from = append(from, line)
		}
	}
	// log.Fatalf("Sorry, we haven't be able to complete this implementation yet, keep tuned")
	return
}

func Defaultimplementation(url *string) (output string) {
	output = "{\n  \"data\": {\n"
	imputFile := Downloadfile(*url)
	repos := Readdata(imputFile)
	for i, element := range repos {
		output = output + "    \"" + element.Url + ":" + element.Hash + "\": {\n"
		dockerfiles := Readrepo(element.Url, element.Hash)
		for j, file := range dockerfiles {
			output = output + "      \"" + file.Name + "\": [\n"
			fromStrings := Readfile(file)
			for k, imageFrom := range fromStrings {
				if k < len(fromStrings)-1 {
					output = output + "        \"" + imageFrom + "\",\n"
				} else {
					output = output + "        \"" + imageFrom + "\"\n"
				}
			}
			if j < len(dockerfiles)-1 {
				output = output + "      ],\n"
			} else {
				output = output + "      ]\n"
			}
		}
		if i < len(repos)-1 {
			output = output + "    },\n"
		} else {
			output = output + "    }\n"
		}
	}
	output = output + "  }\n}"
	return
}

func Jsonimplementation(url *string) (output string) {
	var tempJson jsonWrapper
	var tempData data
	var tempRepo repository
	var tempDocker dockerfile
	imputFile := Downloadfile(*url)
	repos := Readdata(imputFile)
	for _, repo := range repos {
		tempRepo.Url = repo.Url
		tempRepo.Hash = repo.Hash
		dockerfiles := Readrepo(repo.Url, repo.Hash)
		for _, dockerfile := range dockerfiles {
			tempDocker.Pathfile = dockerfile.Name
			fromStrings := Readfile(dockerfile)
			for _, from := range fromStrings {
				tempDocker.Froms = append(tempDocker.Froms, from)
			}
			tempRepo.Dockerfiles = append(tempRepo.Dockerfiles, tempDocker)
		}
		tempData.Repositories = append(tempData.Repositories, tempRepo)
	}
	tempJson.Data = tempData
	json, err := json.Marshal(tempJson)
	if err != nil {
		log.Fatalf("Sorry, but we haven't be able to convert our struct to json format %s", err)
	}
	output = string(json)
	return
}
