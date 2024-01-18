package models

import (
	"fmt"

	"github.com/kiwipanel/scaffolding/pkg/database"
	"github.com/kiwipanel/scaffolding/pkg/database/schema"
)

func ReadPanel() (schema.Panel, error) {
	var panel schema.Panel
	result := database.DB.First(&panel, 1)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	} else if result.RowsAffected == 0 {
		fmt.Println("No record found")
		return panel, nil
	} else {
		// Record found, do something with 'panel'
		fmt.Println("ID:", panel.ID)
		fmt.Println(panel.Passcode)
		return panel, nil
	}
	return panel, nil
}
