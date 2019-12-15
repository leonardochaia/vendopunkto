package main

import (
	"github.com/leonardochaia/vendopunkto/internal/cmd"
	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
}
