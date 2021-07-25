package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

/**
 * User Creation
 */
func createUsers(usersToAdd []string) {
	var currentUsers = getUsersList()
	usersToAdd = checkUserList(currentUsers, usersToAdd)

	setUpUsersIntoWireGuard(usersToAdd, currentUsers)

}

func setUpUsersIntoWireGuard(usersToAdd []string, users []User) {
	for _, u := range usersToAdd {
		key, errGCK := generateClientKeys(u)
		if errGCK != nil {
			fmt.Printf("There was an error creating the key for %s \n", u)
		}

		ip, errGMNIPA := giveMeNextIPAvailable(users)
		if errGMNIPA != nil {
			fmt.Printf("There was an error creating the key for %s \n", u)
		}

		cmd := exec.Command("wg", "set", "wg0", "peer", key, "allowed-ips", string(ip))

		_, errO := cmd.Output()
		if errO != nil {
			fmt.Printf("There was an error setting up %s user\n", u)
		}

		newUser := User{u, ip, time.Now()}

		err := registerUserIntoCSV(newUser)
		if err == nil {
			users = append(users, newUser)
		}

	}
}

func registerUserIntoCSV(user User) error {
	return nil
}

func giveMeNextIPAvailable(userList []User) (net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(CIDR)
	if err != nil {
		return nil, err
	}
	flag := false
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		for _, user := range userList {
			if user.ip.Equal(ip) {
				flag = true
				break
			}
		}
		if !flag {
			return ip, nil
		}
	}

	return nil, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

/*
 * we are going to generate the Private and the Public Key for the User  and we will retrieve the public
 */
func generateClientKeys(u string) (string, error) {
	creationUserCmd := exec.Command("wg", "wg", "genkey", "|", "tee", "/root/wg-user/"+u, "|", "wg", "pubkey", ">", "/root/wg-user/"+u+".pub")
	_, err := creationUserCmd.Output()

	if err != nil {
		return "", err
	}

	cmd := exec.Command("cat", "/root/wg-user/"+u+".pub")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	key := string(output[:len(output)-1])

	return key, nil
}

/*
 * This function is going to check if any user that is in the list is already in the WireGuard config
 * and it will remove it from the list and it will print that user is already in the server
 */
func checkUserList(currentUsers []User, users []string) []string {
	if len(currentUsers) < 1 {
		return users
	}

	var cleanUserList []string
	userFlag := false

	for _, newUser := range users {
		for _, user := range currentUsers {
			if user.name == newUser {
				userFlag = true
				break
			}
		}
		if userFlag {
			cleanUserList = append(cleanUserList, newUser)
		} else {
			userFlag = false
		}
	}

	return cleanUserList
}
