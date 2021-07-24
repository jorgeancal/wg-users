package main

import "strings"

func routerAction(actions string, users []string) {
	switch strings.ToLower(actions) {
	case "create":
		createUsers(users)
	case "update":
		updateUsers(users)
	case "delete":
		deleteUsers(users)
	case "list":
		listUsers()
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
