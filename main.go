package main

import (
	"ecom-be/app/auth"
	"ecom-be/app/database"
	"ecom-be/app/handle"
	myimg "ecom-be/app/image"
	"ecom-be/app/product"
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
	img := myimg.NewService(storage)
	pr := product.NewService(storage)
	router, err := handle.Handler(auth, img, pr)
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
