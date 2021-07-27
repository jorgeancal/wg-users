package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const wgQuickConf = `[Interface]
Address = $address
PrivateKey = $PrivateKey
DNS = 10.0.0.121,1.1.1.1

[Peer]
PublicKey = $PublicKey
PresharedKey = $PresharedKey
AllowedIPs = 0.0.0.0/0
Endpoint =  $endpoint
PersistentKeepalive = 15
`

/**
 * User Creation
 */
func createUsers(usersToAdd []string) {
	var currentUsers = getUsersList()
	usersToAdd = checkUserList(currentUsers, usersToAdd)
	if len(usersToAdd) > 0 {
		setUpUsersIntoWireGuard(usersToAdd, currentUsers)
	}
}

func setUpUsersIntoWireGuard(usersToAdd []string, users []User) {
	for _, username := range usersToAdd {
		newUser := User{}
		newUser.name = username
		newUser, errGCK := generateClientKeys(newUser)
		if errGCK != nil {
			fmt.Printf("There was an error creating the key for %s \n", username)
			os.Exit(-1)
		}

		ip, errGMNIPA := giveMeNextIPAvailable(users)
		newUser.ip = ip.String()
		if errGMNIPA != nil {
			fmt.Printf("There was an error retrieveing the IP for %s \n", username)
			os.Exit(-1)
		}
		if ip != nil {
			var command = "wg set wg0 peer '" + newUser.publicKey + "' preshared-key /etc/wg-users/" + newUser.name + "/" + newUser.name + ".psk allowed-ips " + ip.String()
			cmd := exec.Command("bash", "-c", command)
			_, errO := cmd.Output()
			if errO != nil {
				fmt.Printf("There was an error setting up %s user - error :\n %v \n", username, errO)
				os.Exit(-1)
			}

			err := registerUserIntoCSV(newUser)
			if err == nil {
				users = append(users, newUser)
			}

			err = createWGQuickConfig(newUser)
			if err == nil {
				fmt.Println("Users added correctly")
				os.Exit(-1)
			}
		}
	}
}

func createWGQuickConfig(user User) error {
	userConfig := wgQuickConf
	userConfig = strings.Replace(userConfig, "$address", user.ip, 1)
	userConfig = strings.Replace(userConfig, "$PrivateKey", user.privateKey, 1)
	userConfig = strings.Replace(userConfig, "$PublicKey", user.publicKey, 1)
	userConfig = strings.Replace(userConfig, "$PresharedKey", user.presharedKey, 1)
	userConfig = strings.Replace(userConfig, "$endpoint", getEndPoint()+":"+wg0["ListenPort"], 1)
	f, err := os.OpenFile("/etc/wg-users/"+user.name+"/"+user.name+"-wg0.conf", os.O_CREATE|os.O_RDWR, 0600)
	if _, err = f.WriteString(userConfig); err != nil {
		return err
	}

	c := exec.Command("bash", "-c", "qrencode -o /etc/wg-users/"+user.name+"/"+user.name+".png < /etc/wg-users/"+user.name+"/"+user.name+"-wg0.conf")
	_, err = c.Output()
	if err != nil {
		return err
	}
	return nil
}

func getEndPoint() string {
	c := exec.Command("cat", DIRS[2]+"/endpoint.conf")
	outputC, err := c.Output()
	if err != nil {
		return ""
	}
	return string(outputC[:len(outputC)-1])
}

func registerUserIntoCSV(user User) error {
	f, err := os.OpenFile(FILES[0], os.O_APPEND|os.O_WRONLY, 0600)
	var lineToWrite = user.name + "\t" + user.ip + "\t" + user.creation.Format(time.RFC822) + "\t" + user.publicKey + "\t" + user.privateKey + "\t" + user.presharedKey + "\n"
	if _, err = f.WriteString(lineToWrite); err != nil {
		return err
	}
	return nil
}

func giveMeNextIPAvailable(userList []User) (net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(wg0["Address"])
	if err != nil {
		return nil, err
	}
	flag := false
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		if ip.String() != "10.51.0.0" && ip.String() != "10.51.0.1" {
			if len(userList) == 0 {
				return ip, nil
			}
			for _, user := range userList {
				if user.ip == ip.String() {
					flag = true
					break
				}
			}
			if flag == false {
				return ip, nil
			} else {
				flag = false
			}
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
func generateClientKeys(user User) (User, error) {
	userFolder := "/etc/wg-users/" + user.name + "/"
	_, err := os.Stat(userFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(userFolder, 0700)
		if errDir != nil {
			return user, errDir
		}
	}

	var command = "wg genkey | tee " + userFolder + user.name + " | wg pubkey > " + userFolder + user.name + ".pub && wg genpsk > " + userFolder + user.name + ".psk"
	creationUserCmd := exec.Command("bash", "-c", command)
	_, err = creationUserCmd.Output()
	if err != nil {
		return user, err
	}

	cmd := exec.Command("cat", userFolder+user.name+".pub")
	output, err := cmd.Output()
	if err != nil {
		return user, err
	}
	user.publicKey = string(output[:len(output)-1])

	md := exec.Command("cat", userFolder+user.name)
	outputP, err := md.Output()
	if err != nil {
		return user, err
	}
	user.privateKey = string(outputP[:len(outputP)-1])

	c := exec.Command("cat", userFolder+user.name+".psk")
	outputC, err := c.Output()
	if err != nil {
		return user, err
	}
	user.presharedKey = string(outputC[:len(outputC)-1])

	return user, nil
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
		if !userFlag {
			cleanUserList = append(cleanUserList, newUser)
		} else {
			userFlag = false
		}
	}

	return cleanUserList
}
