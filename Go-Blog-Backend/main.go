package main

import (
	"Go-Blog/database"
	"Go-Blog/model"
	"Go-Blog/router"
	"Go-Blog/utils"
	"log"
	"os"
)

func main() {

	// database initialize
	database.DBInit()

	pwd, err := utils.PasswordHash(os.Getenv("password"))
	if err != nil {
		panic(err)
	}

	id, err := utils.GenerateId(1)
	if err != nil {
		panic(err)
	}

	database.DB.FirstOrCreate(&model.User{}, &model.User{
		Username: "admin",
		Password: pwd,
		Email:    "user@example.com",
		Id:       id,
	})

	model.InitRedis()

	// router initialize
	err = router.SetupRouter().Run("localhost:8080")
	if err != nil {
		log.Fatalf("[Error]: %v", err)
	}
}
