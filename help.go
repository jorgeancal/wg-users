package main

import "fmt"

/**
 * Print Helps
 */
func printHelp() {
	help := `usage: wg-users [actions] [<users>]
actions: 
	create:	creates the user/users. It will create the config for the WireGuard client in the home root folder.
		Example:
			wg-users create foo bar
	update:	updates the user/users. It will update the credentials of the user/users.
		Example:
			wg-users update foo bar
	delete:	deletes the user/users. It 
		Example:
			wg-users delete foo bar
	list: list the users of we have 
		Example:
			wg-users list`
	fmt.Printf("%s\n", help)
}
