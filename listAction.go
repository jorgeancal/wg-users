package main

import (
	"bufio"
	"fmt"
	"net"
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
	}

	scanner := bufio.NewScanner(tsvFile)
	counter := 0
	for scanner.Scan() {
		if counter > 0 {
			line := scanner.Text()
			columns := strings.Split(line, "\t")
			creationTime, err := time.Parse(time.RFC822, columns[2])
			if err != nil {
				fmt.Print("There is something wrong in the tsv")
			}
			currentUsers = append(currentUsers, User{columns[0], net.IP(columns[1]), creationTime})
		}
		counter++
	}

	return currentUsers
}
