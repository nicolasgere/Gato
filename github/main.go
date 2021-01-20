package github

import (
	"context"
	"github.com/getlantern/systray"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

const (
	Queued    int = 2
	Cancelled int = 0
	Done      int = 1
	Running   int = 3
	Fail      int = 4
)

type Action struct {
	Branch  string
	Status  int
	Started *github.Timestamp
	Menu    *systray.MenuItem
}

func GetLastAction(repo string, owner string, token string, author string) (actions []*Action, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	rep, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, &github.ListWorkflowRunsOptions{
		Actor: author,
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	})
	if err != nil {
		return
	}
	duplicate := map[string]bool{}

	for _, e := range rep.WorkflowRuns {
		status := 0
		switch *e.Status {
		case "in_progress":
			status = Running
		case "queued":
			status = Running
		case "completed":
			switch *e.Conclusion {
			case "success":
				status = Done
			case "failure":
				status = Fail
			case "cancelled":
				status = Cancelled

			}

		}
		a := Action{
			Status:  status,
			Branch:  *e.HeadBranch,
			Started: e.CreatedAt,
		}
		if duplicate[*e.HeadBranch] {
			continue
		}
		actions = append(actions, &a)
		duplicate[*e.HeadBranch] = true
		if len(actions) > 2 {
			break
		}
	}
	return
}
