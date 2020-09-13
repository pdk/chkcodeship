package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codeship/codeship-go"
)

func main() {

	login := os.Getenv("CODESHIP_LOGIN")
	password := os.Getenv("CODESHIP_PASSWORD")

	if login == "" || password == "" {
		log.Fatalf("chkcodeship requires CODESHIP_LOGIN and CODESHIP_PASSWORD env vars for authentication")
	}

	if len(os.Args) < 4 {
		log.Fatalf("usage: chkcodeship organization project branch")
	}

	organizationName := os.Args[1]
	projectName := os.Args[2]
	branchName := os.Args[3]

	auth := codeship.NewBasicAuth(login, password)
	client, err := codeship.New(auth)
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}

	ctx := context.Background()

	org, err := client.Organization(ctx, organizationName)
	if err != nil {
		log.Fatalf("failed to switch to org: %v", err)
	}

	projList, _, err := org.ListProjects(ctx)
	if err != nil {
		log.Fatalf("failed to get list of projects: %v", err)
	}

	proj, err := selectProject(projectName, projList)
	if err != nil {
		log.Fatalf("%v", err)
	}

	builds, _, err := org.ListBuilds(ctx, proj.UUID, codeship.PerPage(50))
	if err != nil {
		log.Fatalf("failed to get builds: %v", err)
	}

	var lastSuccessTime time.Time
	var lastSuccessCommit string

	for _, build := range builds.Builds {
		if build.Branch != branchName {
			continue
		}

		if build.Status == "success" && build.FinishedAt.After(lastSuccessTime) {
			lastSuccessCommit = build.CommitSha
			lastSuccessTime = build.FinishedAt
		}
	}

	if lastSuccessCommit == "" {
		log.Fatalf("failed to identify last successful commit for %s", branchName)
	}

	fmt.Printf("%s %s\n", lastSuccessCommit, lastSuccessTime.String())
}

func selectProject(projectName string, projList codeship.ProjectList) (codeship.Project, error) {

	names := []string{}

	for _, proj := range projList.Projects {

		if proj.Name == projectName {
			return proj, nil
		}

		names = append(names, proj.Name)
	}

	return codeship.Project{}, fmt.Errorf("failed to identify project %s in %s", projectName, names)
}
