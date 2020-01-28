package main

import (
	"blogServer/models"
	"blogServer/network"
)

func main() {
	db1 := new(models.DBHandler)
	db1.Init()
	db1.Migrate()

	master := network.ServerMaster{Address: "localhost", Port: "1234", Database: *db1}
	master.ListenToNodes()
	// api.StartAPI(*db1, master.Address, 1235)
}
