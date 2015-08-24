package mctraining

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
)

var session *r.Session

func DBConnect() {
	var err error
	session, err = r.Connect(r.ConnectOpts{Address: "localhost:30815", Database: "materialscommons"})
	if err != nil {
		fmt.Println("Unable to connect to database on localhost:30815:", err)
		os.Exit(1)
	}
}
