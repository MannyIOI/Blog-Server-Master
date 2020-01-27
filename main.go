package main

import (
	"blogServer/api"
	"blogServer/models"
)

func main() {
	var db models.DBHandler = models.DBHandler{}
	db.Init()
	db.Migrate()

	
	go api.StartAPI(db, "http://10.4.107.133", 8080)
	api.StartAPI(db, "http://10.4.107.133", 8081)

}
