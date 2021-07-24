package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var DIRS = []string{
	"/etc/wireguard/",
	"/root/wg-user/",
}
var FILES = []string{
	"/etc/wireguard/users.tsv",
	"/etc/wireguard/wg0.conf",
}

func main() {
	if result, err := isRunningInRoot(); result == false {
		fmt.Printf("This program must be run as root!\n Error - %s", err)
		os.Exit(1)
	}

	if result, err := checkingRequiredFolder(); result == false {
		fmt.Printf("%s \nThere was an error creating the files", err)
		os.Exit(1)
	}

	if result, err := checkingRequiredFiles(); result == false {
		fmt.Printf("%s \nThere was an error creating the files", err)
		os.Exit(1)
	}

	if len(os.Args) <= 1 {
		printHelp()
	}
}

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

func checkingRequiredFolder() (bool, error) {
	for _, dirPath := range DIRS {
		_, err := os.Stat(dirPath)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(dirPath, 0700)
			if errDir != nil {
				return false, errDir
			}
		}
	}
	return true, nil
}

func isRunningInRoot() (bool, error) {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	i, uErr := strconv.Atoi(string(output[:len(output)-1]))

	if uErr != nil {
		return false, uErr
	}

	return i == 0, err
}

func checkingRequiredFiles() (bool, error) {
	for _, filePath := range FILES {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			theFile, fErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if fErr != nil {
				return false, fErr
			}
			tErr := theFile.Close()
			if tErr != nil {
				return false, tErr
			}
		}
	}
	return true, nil
}
