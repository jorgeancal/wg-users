package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getWGServer() error {
	configFile, _ := os.Open(FILES[1])
	wg0 = make(map[string]string)
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
			wg0[strings.TrimSpace(value[0])] = strings.TrimSpace(value[1])
		}
	}
	wg0["ServerPublicKey"] = getServerPublicKey()
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Printf("There was an error closing the %s file\n", FILES[0])
		}
	}(configFile)

	return nil
}

func getServerPublicKey() string {
	c := exec.Command("cat", FILES[3])
	outputC, err := c.Output()
	if err != nil {
		fmt.Println("Error reading the server public key")
		os.Exit(-1)
	}
	return string(outputC[:len(outputC)-1])
}

func configEndPoint(arguments []string) {
	var isEndPoint = false
	for _, argument := range arguments {
		if strings.Contains(argument, "endpoint") {
			isEndPoint = true
		}
		if isEndPoint {
			var command = "echo '" + argument + "' > " + DIRS[2] + "/endpoint.conf"
			cmd := exec.Command("bash", "-c", command)
			_, errO := cmd.Output()
			if errO != nil {
				fmt.Printf("Error setting endpoint \n %v \n", errO)
				os.Exit(-1)
			}
		}
	}
}
