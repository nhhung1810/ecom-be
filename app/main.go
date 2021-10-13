package main

import (
	"ecom-be/app/auth"
	"ecom-be/app/database"
	"ecom-be/app/handle"
	"log"

	"os"
)

func main() {
	// models.Connect() //check here
	storage, err := database.NewStorage()
	if err != nil {
		log.Fatal("Database error: ", err)
		return
	}
	auth := auth.NewService(storage)

	router, err := handle.Handler(auth)
	if err != nil {
		print(err)
		log.Fatal("Router error: ", err)
		return
	}
	router.Run()
}

func setEnvironment(key string, value string) {
	os.Setenv(key, value)
	// HOW TO GET THE ENV VARIABLE
	// os.Getenv(Key)
}
