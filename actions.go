package main

import "strings"

func routerAction(actions string, arguments []string) {
	switch strings.ToLower(actions) {
	case "create":
		createUsers(arguments)
	case "update":
		updateUsers(arguments)
	case "delete":
		deleteUsers(arguments)
	case "list":
		listUsers()
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
