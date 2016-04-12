// Fetch prints the content found at a URL.
package main

import (
	"fmt"
	"log"
	"os"
	"time"
    "strings"

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
    return 0
}
	//

func GetReNo(repo string) int{
	client := github.NewClient(nil)
    k := time.Now()
    d, _ := time.ParseDuration("-240h")
	commitopts := &github.CommitsListOptions{Path: reno, Since: k.Add(d)}

	lastcommits, err := GetLatestCommit("openstack", repo, client, commitopts)
	if err != nil {
		return -1
	}
	fmt.Println(lastcommits)
    for f := range lastcommits {
        GetCommitDetail(client, repo, lastcommits[f])
    }
	return 0
}

func GetCommitDetail(sgc *github.Client, repo, sha string) {
	respcommit, res, err := sgc.Repositories.GetCommit("openstack", repo, sha)
	if err != nil {
		log.Printf("err: %s res: %s", err, res)
	}
	for f := range respcommit.Files {
        fn := *respcommit.Files[f].Filename
        if strings.Contains(fn, "releasenotes/") {
            fmt.Println(sha, fn)
        }
	}
}

func GetLatestCommit(owner, repo string, sgc *github.Client, opts *github.CommitsListOptions) ([]string, error) {
	commits, res, err := sgc.Repositories.ListCommits(owner, repo, opts)

    var comms []string

	if err != nil {
		log.Printf("err: %s res: %s", err, res)
		return comms, err
	}


	log.Printf("last commit length is: %d", len(commits))

    if len(commits) > 0 {
        for c := range commits {
	        log.Printf("commit : %s", *commits[c].SHA)
            comms = append(comms, *commits[c].SHA)
        }
    }
    return comms, nil
}

func main() {
//    fmt.Println(os.Args[1])
//  args 1 is repo name
	fmt.Println(GetReNo(os.Args[1]))
}
