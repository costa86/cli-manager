package main

import (
	db "github.com/costa86/cli-manager/database"
	"github.com/costa86/cli-manager/tui"
)

func main() {
	defer db.DB.Close()
	tui.ShowMenuTui()
	tui.LoadPages()
}
