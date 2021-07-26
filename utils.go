package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getWGServer() error {
	configFile, _ := os.Open(FILES[1])

	scanner := bufio.NewScanner(configFile)
	var isGlobalConfig = false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[") {
			if strings.Contains(line, "[Interface]") {
				isGlobalConfig = true
			} else {
				isGlobalConfig = false
			}
		}
		if strings.Contains(line, "=") && isGlobalConfig {
			value := strings.Split(line, "=")
			wg0[value[0]] = value[1]
		}
	}

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Printf("There was an error closing the %s file\n", FILES[0])
		}
	}(configFile)

	return nil
}
