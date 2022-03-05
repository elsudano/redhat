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

// Data structure that stores the root node of the JSON format
// required to complete the exercise
type JsonWrapper struct {
	Data Data `json:"data"`
}

// Data structure that stores an array of repositories to make
// it easier to iterate through them when we get the output of the
// JSON output
type Data struct {
	Repositories []Repository `json:"repos"`
}

// Data structure that keeps the data in a repository in order,
// so that the data can be accessed from URLs and the
// URL and the HASH required by the input file.
// also an array with the Dockerfiles associated with the repository is maintained.
type Repository struct {
	Url         string       `json:"url"`
	Hash        string       `json:"hash"`
	Dockerfiles []Dockerfile `json:"dockerfile"`
}

// Data structure that keeps the information in the Dockerfiles in order,
// so that you can query the data when you get the JSON output.
// the data when the JSON output is obtained.
type Dockerfile struct {
	Pathfile string   `json:"path"`
	Froms    []string `json:"from"`
}

// Data structure that keeps the data in a repository in order,
// so that the URL and HASH data required by the input file can be accessed.
type RepoInfo struct {
	Url  string
	Hash string
}

// Function in charge of downloading the input file and saving
// the data in memory without the need to store the data on disk.
func DownloadFile(uri string) (data []byte) {
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

// Function in charge of reading the list of repositories and
// storing it in the appropriate structure for later processing
// in the search functions.
func ReadData(data []byte) (repos []RepoInfo) {
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

// Function in charge of reading the information of each one
// of the repositories stored with the ReadData() function.
func ReadRepo(path string, hash string) (files []object.File) {
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
	files = FindDokerfiles(commit.Tree())
	return
}

// Auxiliary function in charge of searching the dockerfiles in
// the different repositories indicating the path of the file
// in each repository.
func FindDokerfiles(tree *object.Tree, err error) (files []object.File) {
	tree.Files().ForEach(func(f *object.File) error {
		match, err := regexp.MatchString(`(?:^|\W)Dockerfile$`, f.Name)
		if f.Mode.IsFile() && match {
			files = append(files, *f)
		}
		return err
	})
	return
}

// Function in charge of reading a Dockerfile to extract the
// information of the images configured in the file.
func ReadFile(file object.File) (from []string) {
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

// Default implementation, in the exercise it is indicated that
// it is necessary to obtain a specific output when the application
// is executed, this function is in charge of respecting the design
// requests regarding the output.
func DefaultImplementation(url *string) (output string) {
	output = "{\n  \"data\": {\n"
	imputFile := DownloadFile(*url)
	repos := ReadData(imputFile)
	for i, element := range repos {
		output = output + "    \"" + element.Url + ":" + element.Hash + "\": {\n"
		dockerfiles := ReadRepo(element.Url, element.Hash)
		for j, file := range dockerfiles {
			output = output + "      \"" + file.Name + "\": [\n"
			fromStrings := ReadFile(file)
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

// Implementation in enhanced JSON mode, this implementation does
// not take as much into account the output requested in the exercise
// but tries to format the output so that the JSON format is more
// manageable when navigating through the children.
func JsonImplementation(url *string) (output string) {
	var tempJson JsonWrapper
	var tempData Data
	var tempRepo Repository
	var tempDocker Dockerfile
	imputFile := DownloadFile(*url)
	repos := ReadData(imputFile)
	for _, repo := range repos {
		tempRepo.Url = repo.Url
		tempRepo.Hash = repo.Hash
		dockerfiles := ReadRepo(repo.Url, repo.Hash)
		for _, dockerfile := range dockerfiles {
			tempDocker.Pathfile = dockerfile.Name
			fromStrings := ReadFile(dockerfile)
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
