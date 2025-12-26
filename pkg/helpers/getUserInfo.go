package helpers

import (
	"fmt"
	"log"
	"os/user"
)

func GetUserInfo() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Username: %s\n", currentUser.Username)
	fmt.Printf("UID: %s\n", currentUser.Uid)
	fmt.Printf("GID: %s\n", currentUser.Gid)
	fmt.Printf("Home Directory: %s\n", currentUser.HomeDir)
	fmt.Printf("Display Name: %s\n", currentUser.Name)
	return "user info"
}
