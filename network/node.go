package network

import (
	"blogServer/models"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// NODE ROLES MASTER AND SLAVE
var SlaveList []ServerNode = []ServerNode{}

// ServerMaster comment
type ServerMaster struct {
	Address string
	Port    string
	// SlaveList []ServerNode
	Database models.DBHandler
}

// ServerNode comment
type ServerNode struct {
	Address       string
	Port          string
	MasterAddress string
}

type MessageObject struct {
}

func (server ServerMaster) ListenToNodes() {
	serverMaster := new(ServerMaster)

	err := rpc.Register(serverMaster)

	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}

	rpc.HandleHTTP()

	listener, e := net.Listen("tcp", ":"+server.Port)

	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %s", server.Port)
	// Start accept incoming HTTP connections
	http.Serve(listener, nil)
}

//
func (server ServerMaster) NotifyNodesUser(user models.User, reply *models.User) error {
	fmt.Println("NotifyNodesUser has been called remotely ", user)
	fmt.Println(SlaveList)
	for index := 0; index < len(SlaveList); index++ {
		// TODO notify all clients
		serverNode := SlaveList[index]
		var reply models.User
		fmt.Println(serverNode.Address)
		client, err := rpc.DialHTTP("tcp", serverNode.Address+":"+serverNode.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.CreateUser", user, &reply)
		fmt.Println(reply)
	}
	return nil
}

// AddNode comment
func (server ServerMaster) AddNode(node ServerNode, reply *string) error {
	SlaveList = append(SlaveList, node)

	fmt.Println(SlaveList)
	return nil
}
