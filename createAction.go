package main

/**
 * User Creation
 */
func createUsers(usersToAdd []string) {
	var currentUsers []User = getUsersList()
	usersToAdd = checkUserList(currentUsers, usersToAdd)
}

/*
 * This function is going to check if any user that is in the list is already in the WireGuard config
 * and it will remove it from the list and it will print that user is already in the server
 */
func checkUserList(currentUsers []User, users []string) []string {

	return nil
}
