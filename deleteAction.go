package main

import (
	"fmt"
	"io"
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

	removeLineInCSV(users)
}

// we are going to create a wg0-new.conf then mv the old one to the new one.
func removeLineInCSV(users []string) {
	copyConfigFile()
}

func copyConfigFile() {
	sourceFileStat, err := os.Stat(FILES[1])
	if err != nil {
		fmt.Errorf("Error trying to get %s FileInf\n", FILES[1])
	}

	if !sourceFileStat.Mode().IsRegular() {
		fmt.Errorf("%s is not a regular file\n", FILES[1])
	}

	source, err := os.Open(FILES[1])
	if err != nil {
		fmt.Errorf("%s failed to open\n", FILES[1])
	}
	defer source.Close()

	destination, err := os.Create(DIRS[1] + "wg0-new.conf")
	if err != nil {
		fmt.Errorf("%s failed to open\n", DIRS[1]+"wg0-new.conf")
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		fmt.Errorf("%s failed to copy\n", DIRS[1]+"wg0-new.conf")
	}
}
