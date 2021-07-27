package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// DIRS This Variable due to is that static and they are used in other sections
var DIRS = []string{
	"/etc/wireguard/",
	"/etc/wg-users/",
	"/etc/wg-users/config",
}

// FILES This Variable due to is that static and they are used in other sections
var FILES = []string{
	"/etc/wireguard/users.tsv",
	"/etc/wireguard/wg0.conf",
	"/etc/wireguard/serverkey",
	"/etc/wireguard/serverkey.pub",
}

var wg0 map[string]string

/*
	@TODO - Check if it's a good idea use flag package
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

*/
func main() {
	err := getWGServer()
	if err != nil {
		fmt.Printf("There was a problem reading the config.")
		os.Exit(1)
	}
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
		os.Exit(1)
	}

	var actions = os.Args[1]
	var arguments = os.Args[2:]

	routerAction(actions, arguments)
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

	return i == 0, nil
}

func checkingRequiredFiles() (bool, error) {
	for _, filePath := range FILES {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			theFile, fErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if fErr != nil {
				return false, fErr
			}
			if strings.Contains(filePath, "users.tsv") {
				_, err := theFile.WriteString("UserName\tIP\tCreation\tPublic Key\tPrivate Key\tPresharedKey\n")
				if err != nil {
					return false, err
				}
			}
			tErr := theFile.Close()
			if tErr != nil {
				return false, tErr
			}
		}
	}
	return true, nil
}
