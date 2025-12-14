package app

import (
	"fmt"
	"net/http"

	"github.com/kiwipanel/kiwipanel/config"
)

func Boot(environment string) {
	cfg, err := config.Load("", environment)
	if err != nil {
		panic(err)
	}

	fmt.Println("KiwiPanel listening on", cfg.Server.Port)
	fmt.Println("KiwiPanel listening on", cfg)

	r := Register(environment)
	Migrate()
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	http.ListenAndServe(port, r)

}
