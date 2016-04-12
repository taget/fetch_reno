// Fetch prints the content found at a URL.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
)

var (
	myRepos        = []string{"nova", "neutron"}
	reno    string = "releasenotes/notes"
)

func getclient() int {
	client := github.NewClient(nil)
	var allRepos []github.Repository
	listopts := github.ListOptions{Page: 1, PerPage: 10000}
	opt := &github.RepositoryListByOrgOptions{Type: "all", ListOptions: listopts}

	repos, _, err := client.Repositories.ListByOrg("openstack", opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	num := 0
	for i := range repos {
		for r := range myRepos {
			if myRepos[r] == *repos[i].Name {
				num++
				fmt.Println(num)
				fmt.Println(*repos[i].Name)
				allRepos = append(allRepos, repos[i])
			}
		}
	}

	//
	commitopts := &github.CommitsListOptions{Path: reno, Since: time.Time}
	time.Duration(Hour)
	lastcommit, err := GetLatestCommit("openstack", "nova", client, commitopts)
	if err != nil {
		return -1
	}
	fmt.Println(lastcommit)

	respcommit, res, err := client.Repositories.GetCommit("openstack", "nova", lastcommit)
	if err != nil {
		log.Printf("err: %s res: %s", err, res)
		return -1
	}
	for f := range respcommit.Files {
		fmt.Println(*respcommit.Files[f].Filename)

	}

	return 0
}

func GetLatestCommit(owner, repo string, sgc *github.Client, opts *github.CommitsListOptions) (string, error) {
	commits, res, err := sgc.Repositories.ListCommits(owner, repo, opts)

	if err != nil {
		log.Printf("err: %s res: %s", err, res)
		return "", err
	}

	log.Printf("last commit: %s", *commits[0].SHA)

	return *commits[0].SHA, nil
}

func main() {
	fmt.Println(getclient())
}
