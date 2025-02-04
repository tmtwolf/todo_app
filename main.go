package main

import (
	"fmt"
	"log"
	"todo_app/app/controllers"
	"todo_app/app/models"
)

func main() {
	fmt.Println(models.Db)
	fmt.Println(models.Encrypt("password"))

	if err := controllers.StartMainServer(); err != nil {
		log.Fatalln(err)
	}
}
