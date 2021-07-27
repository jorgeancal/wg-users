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
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Printf("There was an error closing the %s file\n", FILES[0])
		}
	}(configFile)

	return nil
}

func configEndPoint(endpoint string) {
	if endpoint != "" {
		var command = "echo '" + endpoint + "' > " + DIRS[2] + "/endpoint.conf"
		cmd := exec.Command("bash", "-c", command)
		_, errO := cmd.Output()
		if errO != nil {
			fmt.Printf("Error setting endpoint \n %v \n", errO)
			os.Exit(-1)
		}
	}
}
