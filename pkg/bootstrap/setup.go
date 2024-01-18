package bootstrap

import (
	"fmt"

	"github.com/kiwipanel/scaffolding/pkg/database"
	"github.com/kiwipanel/scaffolding/pkg/database/schema"
	"github.com/kiwipanel/scaffolding/pkg/helpers"
)

func FirstSetup() {
	passcode := helpers.GenerateRandomString(9)
	database.DB.Create(&schema.Panel{Version: "0.0.1", Passcode: passcode, Firsttime: true, Lock: false})
}

func Setup() {
	var panel schema.Panel
	result := database.DB.First(&panel, 1)

	if result.Error != nil {
		// Handle the error (record not found, connection issue, etc.)
		fmt.Println("Cannot connect to the database. Error:", result.Error)
	} else if result.RowsAffected == 0 {
		// Handle the case where no record was found
		FirstSetup()
		fmt.Println("No record found. Starting to setup")
	} else {
		// Record found, do something with 'panel'
		fmt.Println("ID:", panel.ID)
		fmt.Println(panel.Passcode)
		fmt.Println("Already setup")
	}

}
