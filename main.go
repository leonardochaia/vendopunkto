package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/leonardochaia/vendopunkto/cmd"
)

func main() {
	cmd.Execute()
}
