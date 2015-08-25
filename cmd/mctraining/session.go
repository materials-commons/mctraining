package mctraining

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/mcstore/pkg/db/dai"
)

var (
	session *r.Session
	rusers dai.Users
	rprojs dai.Projects
	rfiles dai.Files
)


func DBConnect() {
	var err error
	session, err = r.Connect(r.ConnectOpts{Address: "localhost:30815", Database: "materialscommons"})
	if err != nil {
		fmt.Println("Unable to connect to database on localhost:30815:", err)
		os.Exit(1)
	}
	rusers = dai.NewRUsers(session)
	rprojs = dai.NewRProjects(session)
	rfiles = dai.NewRFiles(session)
}
