package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

/**
 * list Users
 */
func listUsers() {

}

func getUsersList() []User {
	var currentUsers []User

	tsvFile, err := os.Open(FILES[0])

	defer func(tsvFile *os.File) {
		err := tsvFile.Close()
		if err != nil {
			fmt.Printf("There was an error closing the %s file\n", FILES[0])
		}
	}(tsvFile)

	if err != nil {
		_ = fmt.Errorf("error reading the %s file", FILES[0])
		os.Exit(-1)
	}
	scanner := bufio.NewScanner(tsvFile)
	counter := 0
	for scanner.Scan() {
		if counter >= 1 {
			line := scanner.Text()
			columns := strings.Split(line, "\t")
			creationTime, err := time.Parse(time.RFC822, columns[2])
			if err != nil {
				fmt.Print("There is something wrong in the tsv\n")
				os.Exit(-1)
			}
			currentUsers = append(currentUsers, User{columns[0], columns[1], creationTime, columns[3], columns[4], columns[5]})
		}
		counter++

	}
	return currentUsers
}
