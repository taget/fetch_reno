// Fetch prints the content found at a URL.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

//RenoStruct stands for a Reno
type RenoStruct struct {
	SHA         *string
	Type        int
	FileContent *[]byte
	FileName    *string
}

var (
	myRepos = []string{"nova", "neutron"}
	//Reno path
	Reno = "releasenotes/notes"
	//Period is the period time from now
	Period = "-240h"
	//RenoStruct list
	myRenos = []RenoStruct{}
	//TmpFileDir to keep old SHA
	TmpFileDir = ".data/"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

//GetReNo get reno by sepcify a period , return last commit of reno dir
func GetReNo(client *github.Client, repo, lastcommit string) (string, error) {

	k := time.Now()
	d, _ := time.ParseDuration(Period)
	var commitopts *github.CommitsListOptions
	log.Println("Time period: ", k.Add(d), k)

	if len(lastcommit) > 0 {
		fmt.Println("----------lastcommit--------")
		fmt.Println(lastcommit)
		commitopts = &github.CommitsListOptions{Path: Reno, SHA: lastcommit}
	} else {
		commitopts = &github.CommitsListOptions{Path: Reno, Since: k.Add(d)}
	}
	fmt.Println(*commitopts)
	lastcommits, err := GetLatestCommit("openstack", repo, client, commitopts)
	if err != nil {
		return "", err
	}

	for f := range lastcommits {
		GetCommitDetail(client, repo, lastcommits[f])
	}
	return lastcommits[0], nil
}

//GetCommitDetail get commit details of a repo by specify SHA
func GetCommitDetail(sgc *github.Client, repo, sha string) {
	respcommit, res, err := sgc.Repositories.GetCommit("openstack", repo, sha)
	if err != nil {
		log.Printf("err: %s res: %s", err, res)
	}
	for f := range respcommit.Files {
		fn := *respcommit.Files[f].Filename
		if strings.Contains(fn, Reno) {

			r, err := sgc.Repositories.DownloadContents("openstack", repo, fn, &github.RepositoryContentGetOptions{})
			if err != nil {
				log.Printf("err: %s res: %s", err, res)
				continue
			}
			//fixme 2000 seems is not enough
			p := make([]byte, 2000)
			r.Read(p)
			r.Close()
			myRenos = append(myRenos, RenoStruct{SHA: &sha, Type: 1, FileName: &fn, FileContent: &p})
		}
	}
}

//GetLatestCommit get last commits of a repo by specify opts
func GetLatestCommit(owner, repo string, sgc *github.Client, opts *github.CommitsListOptions) ([]string, error) {
	commits, res, err := sgc.Repositories.ListCommits(owner, repo, opts)

	var comms []string

	if err != nil {
		log.Printf("err: %s res: %s", err, res)
		return comms, err
	}

	log.Printf("Last commit length is: %d", len(commits))

	if len(commits) > 0 {
		for c := range commits {
			log.Printf("commit : %s", *commits[c].SHA)
			comms = append(comms, *commits[c].SHA)
		}
	}
	return comms, nil
}

//GetOldSHA to get last SHA, else return empty
func GetOldSHA(filename string) string {
	if _, err := os.Stat(filename); err == nil {
		sha := make([]byte, 100)
		sha, err = ioutil.ReadFile(filename)
		check(err)
		return string(sha)
	}
	return ""
}

//WriteNewSHA to get last SHA, else return empty
func WriteNewSHA(filename, sha string) error {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()
	_, err = f.WriteString(sha)
	check(err)
	f.Sync()
	return err
}

//Init create TmpDir
func Init() {
	//Igonre errors
	os.Mkdir(TmpFileDir, 0777)
}
func main() {
	//args 1 is repo name
	if len(os.Args) < 2 {
		fmt.Println("You need to specify repo name, for example : `nova`")
		os.Exit(1)
	}
	Init()
	client := github.NewClient(nil)
	repo := os.Args[1]
	oldshafile := TmpFileDir + repo

	lastcommit := GetOldSHA(oldshafile)
	newsha, err := GetReNo(client, os.Args[1], lastcommit)
	check(err)
	fmt.Println("newsha")
	fmt.Println(newsha)
	err = WriteNewSHA(oldshafile, newsha)
	check(err)

	fmt.Println("------------------------")
	for re := range myRenos {
		fmt.Println(*myRenos[re].FileName)
		fmt.Println(string(*myRenos[re].FileContent))
	}
}
