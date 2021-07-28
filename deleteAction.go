package main

import (
	"fmt"
	"os"
	"os/exec"
)

/**
 * Delete Users
 * we will execute the next command `wg set wg0 peer $pubKey remove
 */
func deleteUsers(users []string) {
	currentUsers := getUsersList()
	for _, currentUser := range currentUsers {
		for _, user := range users {
			if currentUser.name == user {
				var command = "wg set wg0 peer " + currentUser.publicKey + " remove"
				cmd := exec.Command("bash", "-c", command)
				_, errO := cmd.Output()
				if errO != nil {
					fmt.Printf("There was an error delete %s user \n", currentUser.name)
					os.Exit(-1)
				}
			}
		}
	}
}
