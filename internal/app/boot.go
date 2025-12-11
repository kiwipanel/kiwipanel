package app

import (
	"fmt"

	"github.com/kiwipanel/kiwipanel/config"
)

func Boot(environment string) {
	cfg, err := config.Load("", environment)
	if err != nil {
		panic(err)
	}

	fmt.Println("KiwiPanel listening on", cfg.Server.Port)
	fmt.Println("KiwiPanel listening on", cfg)

	Register(environment)
	Migrate()
	Setup()
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	r.Logger.Fatal(r.Start(port))
}
