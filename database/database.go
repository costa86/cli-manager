package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Cli struct {
	Name, Description, Path string
	Id                      int
}

// Deletes all database records
func DeleteAllRecords() {
	_, err := DB.Exec("DELETE FROM cli")
	CheckError(err)
}

// Deletes a single database record by id
func DeleteRecordById(id int) {
	stmt, err := DB.Prepare("DELETE FROM cli WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(id)
	CheckError(err)
	defer stmt.Close()
}

// Updates a database record
func UpdateCli(cli Cli) {

	stmt, err := DB.Prepare("UPDATE cli SET name = ?, description = ?, path = ? WHERE id = ?")
	CheckError(err)

	_, err = stmt.Exec(cli.Name, cli.Description, cli.Path, cli.Id)
	CheckError(err)

	defer stmt.Close()

}

// Creates a database record
func CreateCli(cli Cli) {
	stmt, err := DB.Prepare("INSERT INTO cli(name, description, path) VALUES(?,?,?)")
	CheckError(err)
	_, err = stmt.Exec(cli.Name, cli.Description, cli.Path)
	CheckError(err)

	defer stmt.Close()

}

// Gets database instance
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
	CheckError(err)
	db.Exec(create)

	return db
}

// Exits the program if errors are found
func CheckError(e error) {
	if e != nil {
		print(e)
		os.Exit(1)
	}
}

var DB = getDb()

// Checks whether the database has any records
func HasRecords() bool {
	row := DB.QueryRow("SELECT COUNT(*) FROM cli")
	var count int
	err := row.Scan(&count)
	CheckError(err)
	return count > 0

}

// Gets database records
func GetEntriesContainingText(text string) []Cli {
	rows, err := DB.Query("SELECT name, description, path, id FROM cli ORDER BY name ASC;")

	if text != "" {
		rows, err = DB.Query("SELECT name, description, path, id FROM cli WHERE name LIKE ? OR description LIKE ? ORDER BY name ASC;", "%"+text+"%", "%"+text+"%")
	}

	CheckError(err)
	defer rows.Close()

	var entries []Cli

	for rows.Next() {
		var entry Cli
		err := rows.Scan(&entry.Name, &entry.Description, &entry.Path, &entry.Id)
		CheckError(err)
		entries = append(entries, entry)
	}

	return entries
}
