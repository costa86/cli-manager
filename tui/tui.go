package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
	db "github.com/costa86/cli-manager/database"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Prompts for confirmation before deleting all CLI database records
func purgeDatabase() {
	pages.AddPage("purgeConfirmation", purgeConfirmation, true, true)
	pages.SwitchToPage("purgeConfirmation")
	purgeConfirmation.ClearButtons()
	purgeConfirmation.SetBackgroundColor(tcell.ColorRed)

	purgeConfirmation.SetText("Are you sure you want to delete ALL the CLIs?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			defer ShowMenuTui()

			if buttonLabel == "Delete" {
				db.DeleteAllRecords()
			}
		})
	purgeConfirmation.SetFocus(1)
	app.SetFocus(purgeConfirmation)
}

// Prompts for confirmation before deleting a CLI database record
func deleteCliById(cli db.Cli) {
	pages.AddPage("deleteConfirmation", deleteConfirmation, true, true)
	pages.SwitchToPage("deleteConfirmation")
	deleteConfirmation.ClearButtons()
	deleteConfirmation.SetBackgroundColor(tcell.ColorRed)

	deleteConfirmation.SetText(fmt.Sprintf("Are you sure you want to delete the %s CLI?", cli.Name)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			defer ShowMenuTui()

			if buttonLabel == "Delete" {
				db.DeleteRecordById(cli.Id)
			}
		})
	deleteConfirmation.SetFocus(1)
	app.SetFocus(deleteConfirmation)
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

// Adds a button to a form to take the user back to the main menu
func addMenuButton(form *tview.Form) {
	form.AddButton("MENU", func() {
		ShowMenuTui()
	})
}

// Shows a form to edit a CLI
func editCliFormTui(cli db.Cli) {

	pages.AddPage("editForm", editForm, true, true)
	pages.SwitchToPage("editForm")
	editForm.Clear(true)
	editForm.SetTitle(fmt.Sprintf("EDIT %s", cli.Name))
	editForm.SetBorder(true)

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
			db.UpdateCli(cli)
			ShowMenuTui()
		}
	})

	addMenuButton(editForm)
	app.SetFocus(editForm)
	editForm.SetFocus(0)
}

// Shows a form to create a new CLI
func addCliFormTui() {
	cli := db.Cli{}
	pages.AddPage("add", addForm, true, false)
	pages.SwitchToPage("add")
	addForm.Clear(true)
	addForm.SetTitle("ADD A NEW CLI")
	addForm.SetBorder(true)

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
			db.CreateCli(cli)
			ShowMenuTui()
		}
	})

	addMenuButton(addForm)
	app.SetFocus(addForm)
	addForm.SetFocus(0)
}

// Notifies the user if the fields in a CLI don't match the mininum characters' quantity
func validateMinChars(cli db.Cli) bool {
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

// Shows a form to search for CLIs by name/description
func searchFormTui() {
	var name string
	pages.AddPage("search", searchForm, true, false)
	pages.SwitchToPage("search")
	searchForm.Clear(true)
	searchForm.SetBorder(true)
	searchForm.SetTitle("SEARCH FOR A CLI BY NAME OR DESCRIPTION. IF NO RESULTS ARE FOUND, YOU'LL BE SENT BACK TO THE MENU")

	searchForm.AddInputField("Name/Description", "", 20, nil, func(text string) {
		name = text
	})

	searchForm.AddButton("SEARCH", func() {
		showCliListTui(name)
	})

	addMenuButton(searchForm)
	searchForm.SetFocus(0)
	app.SetFocus(searchForm)
}

// Shows main menu
func ShowMenuTui() {
	pages.SwitchToPage("menu")
	menu.SetBorder(true)
	menu.SetTitle("MAIN MENU")
	menu.Clear()
	list.Clear()

	if db.HasRecords() {
		menu.AddItem("View all CLIs", "", 'v', func() {
			showCliListTui("")
		})

		menu.AddItem("Search for CLIs", "", 's', func() {
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

// Shows a list of CLIs
func showCliListTui(name string) {
	cliList := db.GetEntriesContainingText(name)

	if len(cliList) == 0 {
		ShowMenuTui()
		return
	}

	pages.AddPage("list", list, true, true)
	pages.SwitchToPage("list")
	list.SetBorder(true)
	title := fmt.Sprintf("ALL CLIs: %d", len(cliList))

	if name != "" {
		title = fmt.Sprintf("ALL CLIs containing %s: %d", name, len(cliList))
	}

	list.SetTitle(title)

	for _, v := range cliList {
		v := v
		shortcut := rune(v.Name[0])
		list.AddItem(v.Name, v.Description, shortcut, func() {
			getCliActionTui(v)
		})
	}
	list.AddItem("Menu", "", 'm', func() {
		ShowMenuTui()
	})

}

// Shows a list of actions to be performed on a CLI
func getCliActionTui(cli db.Cli) {
	pages.AddPage("actions", actions, true, true)
	pages.SwitchToPage("actions")
	actions.Clear()
	actions.SetBorder(true)
	actions.SetTitle(fmt.Sprintf("SELECT AN ACTION FOR %s", cli.Name))

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
		ShowMenuTui()
	})

}

// Runs TUI
func LoadPages() {
	pages.AddPage("menu", menu, true, true)
	err := app.SetRoot(pages, true).EnableMouse(true).Run()
	db.CheckError(err)
}
