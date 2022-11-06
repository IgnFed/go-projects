package main

import (
	"os"

	"github.com/ignfed/web-dev/src"
)

func main(){
  src.Server.Start()
}

func init() {
	os.Setenv("DBDRIVER", "mysql")
	os.Setenv("USERDB", "ignfed")
	os.Setenv("PASSDB", "")
	os.Setenv("DBNAME", "recordings")
}

