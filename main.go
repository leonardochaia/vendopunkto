package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/leonardochaia/vendopunkto/internal/cmd"
)

func main() {
	cmd.Execute()
}
