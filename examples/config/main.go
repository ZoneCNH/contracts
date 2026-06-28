package main

import (
	"fmt"
	"time"

	"github.com/ZoneCNH/contracts/pkg/templatex"
)

func main() {
	cfg := templatex.Config{
		Name:    "templatex",
		Timeout: time.Second,
		Secret:  "example",
	}

	fmt.Println(cfg.Sanitize().Secret)
}
