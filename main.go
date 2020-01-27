package main

import (
	"blogServer/api"
	"blogServer/models"
)

func main() {
	var db models.DBHandler = models.DBHandler{}
	db.Init()
	db.Migrate()

	api.StartAPI(db)

}
