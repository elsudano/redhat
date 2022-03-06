## RedHat Test

I will try do my best, in order to develop a tiny tool in order to show my skills with GO-LANG

You can see in this repository a one implementation of the tool with the following funcionalities and restrictions.


## The challenge

Your task is to build a tool that given a list of repositories, it identifies all the Dockerfile files inside each repository, extracts the image names from the FROM
statement, and returns a json with the aggregated information for all the repositories.

You can find the details of the FROM command here: https://docs.docker.com/engine/reference/builder/#from

- The input will be provided as a URL pointing to a plaintext file.
- Each line will have two fields separated by a space: - the https url of the github public repository - the commit SHA to verify.
- You can skip any line that doesn't match this pattern.

Example input: https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt

Example output:
```json
{
  "data": {
    "https://github.com/app-sre/qontract-reconcile.git:30af65af14a2dce962df923446afff24dd8f123e": {
      "dockerfiles/Dockerfile": [ 
          "quay.io/app-sre/qontract-reconcile-base:0.2.1"
      ]
    },
    "https://github.com/app-sre/container-images.git:c260deaf135fc0efaab365ea234a5b86b3ead404": {
      "jiralert/Dockerfile": [
        "registry.access.redhat.com/ubi8/go-toolset:latest",
        "registry.access.redhat.com/ubi8-minimal:8.2"
      ],
      "qontract-reconcile-base/Dockerfile": [
        "registry.access.redhat.com/ubi8/ubi:8.2",
        "registry.access.redhat.com/ubi8/ubi:8.2",
        "registry.access.redhat.com/ubi8/ubi:8.2"
      ]
    }
  }
}
```

## Deliverables

- URL to private GitHub repository with the code.
- A README.md file detailing your implementation and any additional features added.

## Bonus points:

Since we are a cloud-native team, we want to run this as a Kubernetes Job. If you already know kubernetes, that is excellent. If you don't, we will hugely value
you taking the time to check out minikube and figuring out how to use Jobs. The list of repository urls should be provided to the Job with the
REPOSITORY_LIST_URL environment variable, which should point at an url.

Please feel free to implement any additional features that make this project more production ready. Do make sure to document them in the README

## Explain my implementation

- I chose GO because I like this language
- I have done two types of implementations because when I finished with the first one, the default JSON format was difficult to handle for searching.
- The second implementation does not maintain the format of the output requested in the exercise, but in exchange, it allows to manage the JSON in an easier way.
- In order to test the 2 implementations it can be executed as follows:
  - default:
  ```go
  go run . -url <URL of text file> | jq .
  ```
  - the improved one:
  ```go 
  go run . -url <URL of text file> -fix | jq .
  ```
- If you want, you can use the [jq](https://stedolan.github.io/jq/) command in order to see formmated the output

<!-- if you want to generate again the documentation run: gomarkdoc -e -o README.md -->
<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# redhat

```go
import "github.com/elsudano/redhat"
```

Copyright 2022 Carlos de la Torre\. All rights reserved\. Use of this source code is governed by a MIT License license that can be found in the LICENSE file\.

## Index

- [func DefaultImplementation(url *string) (output string)](<#func-defaultimplementation>)
- [func DownloadFile(uri string) (data []byte)](<#func-downloadfile>)
- [func FindDokerfiles(tree *object.Tree, err error) (files []object.File)](<#func-finddokerfiles>)
- [func JsonImplementation(url *string) (output string)](<#func-jsonimplementation>)
- [func ReadFile(file object.File) (from []string)](<#func-readfile>)
- [func ReadRepo(path string, hash string) (files []object.File)](<#func-readrepo>)
- [type Data](<#type-data>)
- [type Dockerfile](<#type-dockerfile>)
- [type JsonWrapper](<#type-jsonwrapper>)
- [type RepoInfo](<#type-repoinfo>)
  - [func ReadData(data []byte) (repos []RepoInfo)](<#func-readdata>)
- [type Repository](<#type-repository>)


## func [DefaultImplementation](<https://github.com/elsudano/redhat/blob/main/redhat.go#L145>)

```go
func DefaultImplementation(url *string) (output string)
```

Default implementation\, in the exercise it is indicated that it is necessary to obtain a specific output when the application is executed\, this function is in charge of respecting the design requests regarding the output\.

## func [DownloadFile](<https://github.com/elsudano/redhat/blob/main/redhat.go#L58>)

```go
func DownloadFile(uri string) (data []byte)
```

Function in charge of downloading the input file and saving the data in memory without the need to store the data on disk\.

## func [FindDokerfiles](<https://github.com/elsudano/redhat/blob/main/redhat.go#L112>)

```go
func FindDokerfiles(tree *object.Tree, err error) (files []object.File)
```

Auxiliary function in charge of searching the dockerfiles in the different repositories indicating the path of the file in each repository\.

## func [JsonImplementation](<https://github.com/elsudano/redhat/blob/main/redhat.go#L182>)

```go
func JsonImplementation(url *string) (output string)
```

Implementation in enhanced JSON mode\, this implementation does not take as much into account the output requested in the exercise but tries to format the output so that the JSON format is more manageable when navigating through the children\.

## func [ReadFile](<https://github.com/elsudano/redhat/blob/main/redhat.go#L125>)

```go
func ReadFile(file object.File) (from []string)
```

Function in charge of reading a Dockerfile to extract the information of the images configured in the file\.

## func [ReadRepo](<https://github.com/elsudano/redhat/blob/main/redhat.go#L90>)

```go
func ReadRepo(path string, hash string) (files []object.File)
```

Function in charge of reading the information of each one of the repositories stored with the ReadData\(\) function\.

## type [Data](<https://github.com/elsudano/redhat/blob/main/redhat.go#L27-L29>)

Data structure that stores an array of repositories to make it easier to iterate through them when we get the output of the JSON output

```go
type Data struct {
    Repositories []Repository `json:"repos"`
}
```

## type [Dockerfile](<https://github.com/elsudano/redhat/blob/main/redhat.go#L44-L47>)

Data structure that keeps the information in the Dockerfiles in order\, so that you can query the data when you get the JSON output\. the data when the JSON output is obtained\.

```go
type Dockerfile struct {
    Pathfile string   `json:"path"`
    Froms    []string `json:"from"`
}
```

## type [JsonWrapper](<https://github.com/elsudano/redhat/blob/main/redhat.go#L20-L22>)

Data structure that stores the root node of the JSON format required to complete the exercise

```go
type JsonWrapper struct {
    Data Data `json:"data"`
}
```

## type [RepoInfo](<https://github.com/elsudano/redhat/blob/main/redhat.go#L51-L54>)

Data structure that keeps the data in a repository in order\, so that the URL and HASH data required by the input file can be accessed\.

```go
type RepoInfo struct {
    Url  string
    Hash string
}
```

### func [ReadData](<https://github.com/elsudano/redhat/blob/main/redhat.go#L75>)

```go
func ReadData(data []byte) (repos []RepoInfo)
```

Function in charge of reading the list of repositories and storing it in the appropriate structure for later processing in the search functions\.

## type [Repository](<https://github.com/elsudano/redhat/blob/main/redhat.go#L35-L39>)

Data structure that keeps the data in a repository in order\, so that the data can be accessed from URLs and the URL and the HASH required by the input file\. also an array with the Dockerfiles associated with the repository is maintained\.

```go
type Repository struct {
    Url         string       `json:"url"`
    Hash        string       `json:"hash"`
    Dockerfiles []Dockerfile `json:"dockerfile"`
}
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->