# CLI Manager

![search-img](search-img.gif)

You know when you want to search for a graphical (GUI) program in your computer and you hit the main menu and start typing its name or visually search for it? This program is the same, but for CLIs! 😎

# Features

* Add new CLIs
* Edit existing CLIs
* Delete CLIs
* Search/filter CLIs
* Copy CLI command to clipboard
* View all CLIs

# Installation
No installation required. Just run the executable via command-line

|OS|Executable|
|--|--|
|Windows|[cli-manager-windows.exe](cli-manager-windows.exe)|
|Linux|[cli-manager-linux](cli-manager-linux)|
|MacOS|[cli-manager-darwin](cli-manager-darwin)|

# Persistence
The CLI records are stored on a local SQLite database file, named `database.sqlite`

# Build from source
Will generate the executables for Windows, Linux and MacOS

Requirements:

* Go (check version in [go.mod](go.mod))
* Make

Command:

    make build