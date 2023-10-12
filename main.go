package main

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	_ "github.com/mattn/go-sqlite3"
)

type Cli struct {
	Name, Description, Path string
	Id                      int
}

var print = fmt.Println

func checkError(e error) {
	if e != nil {
		print(e)
		os.Exit(1)
	}
}

func hasRecords() bool {
	row := db.QueryRow("SELECT COUNT(*) FROM cli")

	var count int
	err := row.Scan(&count)
	checkError(err)

	return count > 0

}

func getEntriesContainingName(db *sql.DB, text string) []Cli {
	rows, err := db.Query("SELECT name, description, path, id FROM cli ORDER BY name ASC;")

	if text != "" {
		rows, err = db.Query("SELECT name, description, path, id FROM cli WHERE name LIKE ? OR description LIKE ? ORDER BY name ASC;", "%"+text+"%", "%"+text+"%")
	}

	checkError(err)
	defer rows.Close()

	var entries []Cli

	for rows.Next() {
		var entry Cli
		err := rows.Scan(&entry.Name, &entry.Description, &entry.Path, &entry.Id)
		checkError(err)
		entries = append(entries, entry)
	}

	return entries
}

func createCli(db *sql.DB, cli Cli) {
	stmt, err := db.Prepare("INSERT INTO cli(name, description, path) VALUES(?,?,?)")
	checkError(err)
	_, err = stmt.Exec(cli.Name, cli.Description, cli.Path)
	checkError(err)

	defer stmt.Close()

}

func updateCli(cli Cli) {

	stmt, err := db.Prepare("UPDATE cli SET name = ?, description = ?, path = ? WHERE id = ?")
	checkError(err)

	_, err = stmt.Exec(cli.Name, cli.Description, cli.Path, cli.Id)
	checkError(err)

	defer stmt.Close()

}

func purgeDatabase() {
	pages.AddPage("purgeConfirmation", purgeConfirmation, true, true)
	pages.SwitchToPage("purgeConfirmation")
	purgeConfirmation.ClearButtons()
	purgeConfirmation.SetBackgroundColor(tcell.ColorRed)

	purgeConfirmation.SetText("Are you sure you want to delete ALL the CLI's?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			defer showMenuTui()

			if buttonLabel == "Delete" {
				_, err := db.Exec("DELETE FROM cli")
				checkError(err)
			}
		})
}

func deleteCliById(cli Cli) {
	pages.AddPage("deleteConfirmation", deleteConfirmation, true, true)
	pages.SwitchToPage("deleteConfirmation")
	deleteConfirmation.ClearButtons()
	deleteConfirmation.SetBackgroundColor(tcell.ColorRed)

	deleteConfirmation.SetText(fmt.Sprintf("Are you sure you want to delete the %s CLI?", cli.Name)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			defer showMenuTui()

			if buttonLabel == "Delete" {
				stmt, err := db.Prepare("DELETE FROM cli WHERE id = ?")
				checkError(err)
				_, err = stmt.Exec(cli.Id)
				checkError(err)
				defer stmt.Close()
			}
		})

}

func getDb() *sql.DB {
	const file string = "database.sqlite"
	const create string = `
	CREATE TABLE IF NOT EXISTS cli (
		id INTEGER PRIMARY KEY AUTOINCREMENT,		
		name TEXT,
		description TEXT,
		path TEXT
	);
	`
	db, err := sql.Open("sqlite3", file)
	checkError(err)
	db.Exec(create)

	return db
}

var app = tview.NewApplication()

var menu = tview.NewList()
var list = tview.NewList()
var pages = tview.NewPages()
var actions = tview.NewList()
var addForm = tview.NewForm()
var editForm = tview.NewForm()
var searchForm = tview.NewForm()
var deleteConfirmation = tview.NewModal()
var modal = tview.NewModal()
var purgeConfirmation = tview.NewModal()

func addMenuButton(form *tview.Form) {
	form.AddButton("MENU", func() {
		showMenuTui()
	})

}

func editCliFormTui(cli Cli) {

	pages.AddPage("editForm", editForm, true, true)
	pages.SwitchToPage("editForm")
	editForm.Clear(true)
	editForm.SetTitle(fmt.Sprintf("EDIT %s", cli.Name))
	editForm.SetBorder(true)
	editForm.SetFocus(0)

	editForm.AddInputField("Name", cli.Name, 30, nil, func(text string) {
		cli.Name = text
	})
	editForm.AddInputField("Description", cli.Description, 30, nil, func(text string) {
		cli.Description = text
	})
	editForm.AddInputField("Path", cli.Path, 30, nil, func(text string) {
		cli.Path = text
	})

	editForm.AddButton("SAVE", func() {
		if validateMinChars(cli) {
			updateCli(cli)
			showMenuTui()
		}
	})

	addMenuButton(editForm)
}

func addCliFormTui() {
	cli := Cli{}
	pages.AddPage("add", addForm, true, false)
	pages.SwitchToPage("add")
	addForm.Clear(true)
	addForm.SetTitle("ADD A NEW CLI")
	addForm.SetBorder(true)
	addForm.SetFocus(0)

	addForm.AddInputField("Name", "", 30, nil, func(text string) {
		cli.Name = text
	})

	addForm.AddInputField("Description", "", 30, nil, func(text string) {
		cli.Description = text
	})

	addForm.AddInputField("Path", "", 30, nil, func(text string) {
		cli.Path = text
	})

	addForm.AddButton("SAVE", func() {
		if validateMinChars(cli) {
			createCli(db, cli)
			showMenuTui()
		}
	})

	addMenuButton(addForm)

}

func validateMinChars(cli Cli) bool {
	minChars := 3
	minCharsReached := len(cli.Name) >= minChars && len(cli.Description) >= minChars && len(cli.Path) >= minChars

	if !minCharsReached {
		modal.ClearButtons()
		pages.AddPage("modal", modal, true, true)
		pages.SwitchToPage("modal")
		modal.SetText(fmt.Sprintf("All the fields require %d+ characters", minChars))
		modal.SetBorder(true)
		modal.AddButtons([]string{"OK"})
		modal.SetFocus(0)
		modal.SetBackgroundColor(tcell.ColorRed)
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.SwitchToPage("add")
			addForm.SetFocus(0)
		})
		return false
	}
	return true
}

func searchFormTui() {
	var name string
	pages.AddPage("search", searchForm, true, false)
	pages.SwitchToPage("search")
	searchForm.Clear(true)
	searchForm.SetBorder(true)
	searchForm.SetTitle("SEARCH FOR A CLI BY NAME OR DESCRIPTION. IF NO RESULTS ARE FOUND, YOU'LL BE BACK TO THE MENU")

	searchForm.AddInputField("Name/Description", "", 20, nil, func(text string) {
		name = text
	})

	searchForm.AddButton("SEARCH", func() {
		showCliListTui(name)
	})

	addMenuButton(searchForm)
}

func showMenuTui() {
	pages.SwitchToPage("menu")
	menu.SetBorder(true)
	menu.SetTitle("MAIN MENU")
	menu.Clear()
	list.Clear()

	if hasRecords() {
		menu.AddItem("View all CLI's", "", 'v', func() {
			showCliListTui("")
		})

		menu.AddItem("Search for CLI's", "", 's', func() {
			searchFormTui()
		})

		menu.AddItem("Purge database", "", 'p', func() {
			purgeDatabase()
		})
	}

	menu.AddItem("Add a new CLI", "", 'a', func() {
		addCliFormTui()
	})

	menu.AddItem("Quit program", "", 'q', func() {
		app.Stop()
	})

}

func showCliListTui(name string) {
	cliList := getEntriesContainingName(db, name)

	if len(cliList) == 0 {
		showMenuTui()
		return
	}

	pages.AddPage("list", list, true, true)
	pages.SwitchToPage("list")
	list.SetBorder(true)
	title := fmt.Sprintf("ALL CLI'S: %d", len(cliList))

	if name != "" {
		title = fmt.Sprintf("ALL CLI'S containing %s: %d", name, len(cliList))
	}

	list.SetTitle(title)

	for _, v := range cliList {
		shortcut := rune(v.Name[0])
		v := v
		list.AddItem(v.Name, v.Description, shortcut, func() {
			getCliActionTui(v)
		})
	}
	list.AddItem("Menu", "", 'm', func() {
		showMenuTui()
	})

}

func getCliActionTui(cli Cli) {
	pages.AddPage("actions", actions, true, true)
	pages.SwitchToPage("actions")
	actions.Clear()
	actions.SetBorder(true)
	actions.SetTitle(fmt.Sprintf("PICK AND ACTION FOR %s", cli.Name))

	actions.AddItem("Delete", "", 'd', func() {
		deleteCliById(cli)
	})

	actions.AddItem("Edit", "", 'e', func() {
		editCliFormTui(cli)
	})

	actions.AddItem(fmt.Sprintf("Get path for %s to your clipboard", cli.Name), "This will exit the program. Then you may use (Ctrl + v) to run the CLI", 'g', func() {
		clipboard.WriteAll(cli.Path)
		app.Stop()
	})

	actions.AddItem("Menu", "", 'm', func() {
		showMenuTui()
	})

}

func loadPages() {
	pages.AddPage("menu", menu, true, true)
	err := app.SetRoot(pages, true).EnableMouse(true).Run()
	checkError(err)
}

var db = getDb()

func main() {
	defer db.Close()

	showMenuTui()
	loadPages()

}
