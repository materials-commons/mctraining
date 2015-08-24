package mctraining

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/mcstore/pkg/db/dai"
	"github.com/materials-commons/mcstore/pkg/db/model"
	"github.com/materials-commons/mcstore/pkg/db/schema"
)

var (
	createProjectCommand = cli.Command{
		Name:    "project",
		Aliases: []string{"proj", "p"},
		Usage:   "Create a new project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "template, t",
				Value: "training_template",
				Usage: "The template project to use",
			},
			cli.StringFlag{
				Name:  "owner",
				Usage: "Project owner",
			},
		},
		Action: createProjectCLI,
	}

	rusers = dai.NewRUsers(session)
	rprojs = dai.NewRProjects(session)
	rfiles = dai.NewRFiles(session)
)

func createProjectCLI(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Println("You must provide a project name")
		os.Exit(1)
	}

	projectName := c.Args()[0]
	owner := c.String("owner")
	if owner == "" {
		owner = projectName + "@mctraining.org"
	}

	templateProjectID := c.String("template")

	createUser(owner)
	projectID := createProject(projectName, owner)
	addFilesFromTemplateProject(templateProjectID, projectID, owner)
}

func createUser(user string) {
	if _, err := rusers.ByID(user); err != nil {
		// User not found so create a new one
		runMCUser(user)
	}

}

func runMCUser(user string) {
	out, err := exec.Command("mcuser.py", "--email="+user, "--password='"+user+"'").Output()
	if err != nil {
		fmt.Println("Unable to add user:", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}

func createProject(name, owner string) string {
	proj := schema.NewProject(name, owner)
	p, err := rprojs.Insert(&proj)
	if err != nil {
		fmt.Printf("Unable to create project %s: %s\n", name, err)
		os.Exit(1)
	}
	return p.ID
}

func addFilesFromTemplateProject(templateProjectID, projectID, owner string) {
	projectDirID := getProjectDirID(projectID)
	rql := r.Table("project2datafile").GetAllByIndex("project_id", templateProjectID)
	var projectFiles []schema.Project2DataFile
	if err := model.ProjectFiles.Qs(session).Rows(rql, &projectFiles); err != nil {
		fmt.Printf("Unable to get files for template project %s: %s\n", templateProjectID, err)
		os.Exit(1)
	}

	for _, p2f := range projectFiles {
		f, err := rfiles.ByID(p2f.DataFileID)
		if err != nil {
			fmt.Printf("Unable to retrieve file %s for project %s: %s\n", p2f.DataFileID, projectID, err)
		}
		f.UsesID = f.ID
		f.ID = ""
		f.Owner = owner
		if _, err := rfiles.Insert(f, projectDirID, projectID); err != nil {
			fmt.Printf("Failed to insert file %s into project %s: %s\n", f.Name, projectID, err)
		}
	}
}

func getProjectDirID(projectID string) string {
	var projectDir schema.Project2DataDir
	rql := r.Table("project2datadir").GetAllByIndex("project_id", projectID)
	if err := model.GetRow(rql, session, &projectDir); err != nil {
		fmt.Printf("Unable to get project directory for %s: %s\n", projectID, err)
		os.Exit(1)
	}

	return projectDir.DataDirID
}