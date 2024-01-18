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

func Migrate() {
	database.DB.AutoMigrate(&schema.Panel{}, &schema.Product{}, &schema.User{})
}

func Setup() {
	Migrate()
	var panel schema.Panel
	result := database.DB.First(&panel, 1)

	if result.Error != nil {
		// Handle the error (record not found, connection issue, etc.)
		fmt.Println("Error:", result.Error)
	} else if result.RowsAffected == 0 {
		// Handle the case where no record was found
		FirstSetup()
		fmt.Println("No record found")
	} else {
		// Record found, do something with 'panel'
		fmt.Println("ID:", panel.ID)
		fmt.Println(panel.Passcode)
	}

}
