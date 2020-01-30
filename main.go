package main

import (
	"blogServer/api"
	"blogServer/network"
)

func main() {
	master := network.ServerMaster{Address: "localhost", Port: "1234"}

	go master.ListenToNodes()
	api.StartAPI(&master, "localhost", 1235)
}
