package main

import (
	"blogServer/network"
)

func main() {
	master := network.ServerMaster{Address: "localhost", Port: "1234"}
	master.ListenToNodes()
}
