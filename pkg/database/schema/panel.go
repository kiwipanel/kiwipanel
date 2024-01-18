package schema

import "gorm.io/gorm"

type Panel struct {
	gorm.Model
	Version   string
	Passcode  string
	Firsttime bool
	Lock      bool
}
