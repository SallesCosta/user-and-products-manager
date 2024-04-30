package main

import (
	"fmt"

	"github.com/sallescosta/crud-api/configs"
)

func main() {
	config, _ := configs.LoadConfig(".")
	fmt.Println(config.DBDriver)

}
