package main

import (
	"ecom-be/app/auth"
	"ecom-be/app/database"
	"ecom-be/app/handle"
	myimg "ecom-be/app/image"
	"ecom-be/app/order"
	"ecom-be/app/product"
	"log"
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
	or := order.NewService(storage)

	router, err := handle.Handler(auth, img, pr, or)
	if err != nil {
		print(err)
		log.Fatal("Router error: ", err)
		return
	}
	router.Run()
}
