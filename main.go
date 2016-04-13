// Fetch prints the content found at a URL.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
    "strconv"

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
func GetReNo(client *github.Client, repo, lastcommit string) (string, time.Time, error) {

	var commitopts *github.CommitsListOptions
    var since, current time.Time

	if len(lastcommit) > 0 {
        //fixme (eliqiao) passing lastcommit to commitopts doesn't work at all
        //commitopts = &github.CommitsListOptions{Path: Reno, SHA: lastcommit}
        //get commit's commiter date
	    respcommit, _ , err := client.Repositories.GetCommit("openstack", repo, lastcommit)
        check(err)
        since = *respcommit.Commit.Committer.Date
        // incrence 1s to avoid get current commit
	    d, _ := time.ParseDuration("1s")
        since = since.Add(d)
	} else {
		current = time.Now()
	    d, _ := time.ParseDuration(Period)
        since = current.Add(d)
	}

	log.Println("Time period: ", since, current)

	commitopts = &github.CommitsListOptions{Path: Reno, Since: since}

	lastcommits, err := GetLatestCommit("openstack", repo, client, commitopts)
	if err != nil {
		return lastcommit, since, err
	}

	for f := range lastcommits {
		GetCommitDetail(client, repo, lastcommits[f])
	}

    if len(lastcommits) > 0 {
        return lastcommits[0], since, nil

    }
    return lastcommit, since, nil
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
		return strings.Replace(string(sha), "\n", "", -1)
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

//


func main() {
	//args 1 is repo name
	if len(os.Args) < 2 {
		fmt.Println("You need to specify repo name, for example : `nova` <option days>")
		os.Exit(1)
	}

	repo := os.Args[1]
	oldshafile := TmpFileDir + repo

    if len(os.Args) > 1 {
        day, err := strconv.Atoi(os.Args[2])
        if err != nil {
		    fmt.Println("You need to specify repo name, for example : `nova` <option days>")
		    os.Exit(1)
        }
        Period = "-" + strconv.Itoa(day * 24) + "h"
        // remove cookie
        os.Remove(oldshafile)
    }
	Init()
    client := github.NewClient(nil)

	lastcommit := GetOldSHA(oldshafile)

	newsha, since, err := GetReNo(client, os.Args[1], lastcommit)
	check(err)
	err = WriteNewSHA(oldshafile, newsha)
	check(err)

	fmt.Printf("------------------------------------[updates releasenotes for %s]---------------------------------\n", repo)
	fmt.Printf("---------[from %s to %s ]---------\n", since, time.Now())
	for re := range myRenos {
		fmt.Println(*myRenos[re].FileName)
		fmt.Println(string(*myRenos[re].FileContent))
	}
}
