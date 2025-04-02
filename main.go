package main

import (
	"fmt"

	"github.com/lealre/gator/internal/config"
)

func main() {
	name := "renan"
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}

	err = cfg.SetUser(name)
	if err != nil {
		fmt.Print(err)
	}
	cfg, _ = config.Read()

	fmt.Printf("%+v\n", cfg)
}
