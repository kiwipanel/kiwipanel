package schema

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Code  string
	Price int
}
