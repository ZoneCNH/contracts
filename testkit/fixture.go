package testkit

import (
	"time"

	"github.com/ZoneCNH/contracts/pkg/templatex"
)

func Config(name string) templatex.Config {
	return templatex.Config{
		Name:    name,
		Timeout: time.Second,
	}
}
