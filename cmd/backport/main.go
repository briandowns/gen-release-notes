package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/gen-release-notes/repository"
	"github.com/briandowns/gen-release-notes/token"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

var (
	name    string
	version string
	gitSHA  string
)

const httpTimeout = time.Second * 10

const usage = `version: %s
Usage: %[2]s [-t token] [-r repo] [-m milestone] [-p prev milestone]
Options:
    -h                   help
    -v                   show version and exit
    -t                   github token (optional)
    -r repo              repository that should be used
	-i issue id          original issue id
	-c commit            commit id that is being bacported
	-b branch(es)        branches issue is being backported to
Examples: 
	# generate release notes for RKE2 for milestone v1.21.5
    %[2]s -r k3s -m v1.21.5+k3s1 -p v1.21.4+k3s1 
`

// retrieveOriginalIssue
func retrieveOriginalIssue(ctx context.Context, client *github.Client) (*github.Issue, error) {
	// org, err := repository.OrgFromRepo(repo)
	// if err != nil {
	// 	return err
	// }

	issue, _, err := client.Issues.Get(ctx, "briandowns", "wings", int(issueID))
	if err != nil {
		return nil, err
	}

	return issue, nil
}

const (
	issueTitle = "[%s] - %s"
	issueBody  = "Backport fix for %s\n\n* #%d"
)

// createBackportIssues
func createBackportIssues(ctx context.Context, client *github.Client, origIssue *github.Issue, branch string) (*github.Issue, error) {
	title := fmt.Sprintf(issueTitle, strings.Title(branch), origIssue.GetTitle())
	body := fmt.Sprintf(issueBody, origIssue.GetTitle(), *origIssue.Number)

	issue, _, err := client.Issues.Create(ctx, "briandowns", "wings", &github.IssueRequest{
		Title:    github.String(title),
		Body:     github.String(body),
		Assignee: origIssue.GetAssignee().Login,
	})
	if err != nil {
		return nil, err
	}

	return issue, nil
}

var (
	vers     bool
	ghToken  string
	repo     string
	commitID string
	issueID  uint
	branches string
)

func main() {
	flag.Usage = func() {
		w := os.Stderr
		for _, arg := range os.Args {
			if arg == "-h" {
				w = os.Stdout
				break
			}
		}
		fmt.Fprintf(w, usage, version, name)
	}

	flag.BoolVar(&vers, "v", false, "")
	flag.StringVar(&ghToken, "t", "", "")
	flag.StringVar(&repo, "r", "", "")
	flag.StringVar(&commitID, "c", "", "")
	flag.UintVar(&issueID, "i", 0, "")
	flag.StringVar(&branches, "b", "", "")
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stdout, "version: %s - git sha: %s\n", version, gitSHA)
		return
	}

	if !repository.IsValidRepo(repo) {
		fmt.Println("error: please provide a valid repository")
		os.Exit(1)
	}

	if ghToken == "" {
		fmt.Println("error: please provide a token")
		os.Exit(1)
	}

	if commitID == "" {
		fmt.Println("error: please provide a commit id")
		os.Exit(1)
	}

	if issueID == 0 {
		fmt.Println("error: please provide a valid issue id")
		os.Exit(1)
	}

	backportBranches := strings.Split(branches, ",")
	if len(backportBranches) < 1 || backportBranches[0] == "" {
		fmt.Println("error: please provide at least one branch to perform the backport")
		os.Exit(1)
	}

	ctx := context.Background()

	ts := token.TokenSource{
		AccessToken: ghToken,
	}
	oauthClient := oauth2.NewClient(ctx, &ts)
	oauthClient.Timeout = httpTimeout
	client := github.NewClient(oauthClient)

	origIssue, err := retrieveOriginalIssue(ctx, client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, branch := range backportBranches {
		ni, err := createBackportIssues(ctx, client, origIssue, branch)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Backport issue created: " + ni.GetHTMLURL())
	}

	os.Exit(0)
}
