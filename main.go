package main

import (
	"fmt"

	"github.com/DeniesKresna/ecommerceapi/Configs"
	"github.com/DeniesKresna/ecommerceapi/Routers"
	check "github.com/asaskevich/govalidator"
)

func main() {
	check.SetFieldsRequiredByDefault(true)
	if err := Configs.DatabaseInit(); err != nil {
		fmt.Println("status ", err)
	}

	Configs.DatabaseMigrate()

	r := Routers.SetupRouter()
	r.Run(":8090")
}
