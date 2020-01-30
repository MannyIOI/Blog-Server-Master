package network

import (
	"blogServer/models"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

// NODE ROLES MASTER AND SLAVE

// SlaveList comment
var SlaveList []ServerNode = []ServerNode{}
var currentSlaveServerIndex int = 0

// ServerMaster comment
type ServerMaster struct {
	Address string
	Port    string
}

// ServerNode comment
type ServerNode struct {
	Address       string
	Port          string
	MasterAddress string
}

// MessageObject comment
type MessageObject struct {
}

// ListenToNodes comment
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

// NotifyNodesUser comment
func (server ServerMaster) NotifyNodesUser(user models.User, reply *models.User) error {
	// fmt.Println("NotifyNodesUser has been called remotely ", user)
	// fmt.Println(SlaveList)
	for index := 0; index < len(SlaveList); index++ {
		// TODO notify all clients
		serverNode := SlaveList[index]
		var reply models.User
		// fmt.Println(serverNode.Address)
		client, err := rpc.DialHTTP("tcp", serverNode.Address+":"+serverNode.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.CreateUser", user, &reply)
		// fmt.Println(reply)
	}
	return nil
}

// NotifyNodesBlogCreate comment
func (server ServerMaster) NotifyNodesBlogCreate(blog models.Blog, reply *models.Blog) error {
	// fmt.Println("NotifyNodesBlogCreate has been called remotely ", blog)
	for index := 0; index < len(SlaveList); index++ {
		// TODO notify all clients
		serverNode := SlaveList[index]
		var reply models.User
		// fmt.Println(serverNode.Address)
		client, err := rpc.DialHTTP("tcp", serverNode.Address+":"+serverNode.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.CreateBlog", blog, &reply)
		// fmt.Println(reply)
	}
	return nil
}

// NotifyNodesBlogUpdate comment
func (server ServerMaster) NotifyNodesBlogUpdate(blog models.Blog, reply *models.Blog) error {
	for index := 0; index < len(SlaveList); index++ {
		// TODO notify all clients
		serverNode := SlaveList[index]
		var reply models.User
		client, err := rpc.DialHTTP("tcp", serverNode.Address+":"+serverNode.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.UpdateBlogContent", blog, &reply)
	}
	return nil
}

// NotifyNodesBlogTitleUpdate comment
func (server ServerMaster) NotifyNodesBlogTitleUpdate(blog models.Blog, reply *models.Blog) error {
	for index := 0; index < len(SlaveList); index++ {
		// TODO notify all clients
		serverNode := SlaveList[index]
		var reply models.User
		client, err := rpc.DialHTTP("tcp", serverNode.Address+":"+serverNode.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.UpdateBlogTitle", blog, &reply)
	}
	return nil
}

// AddNode comment
func (server ServerMaster) AddNode(node ServerNode, reply *string) error {
	addNode := true
	for index := 0; index < len(SlaveList); index++ {
		if SlaveList[index] == node {
			addNode = false
		}
	}

	if addNode && len(SlaveList) > 0 {
		serverNode := SlaveList[0]
		var replyUsers []models.User
		var replyBlogs []models.Blog
		client, err := rpc.DialHTTP("tcp", serverNode.Address+":"+serverNode.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.GetAllBlogs", "", &replyBlogs)
		client.Call("DBHandler.GetAllUsers", "", &replyUsers)

		client, err = rpc.DialHTTP("tcp", node.Address+":"+node.Port)
		if err != nil {
			log.Fatal("error")
		}
		client.Call("DBHandler.BatchInsertBlogs", replyBlogs, nil)
		client.Call("DBHandler.BatchInsertUsers", replyUsers, nil)
	}

	if addNode {
		SlaveList = append(SlaveList, node)
	}

	return nil
}

// GetAvailableServer comment
func (server ServerMaster) GetAvailableServer(slaveNode *ServerNode, err *error) {
	if len(SlaveList) > 0 {
		err = nil
		currentSlaveNode := SlaveList[currentSlaveServerIndex]
		port, _ := strconv.Atoi(currentSlaveNode.Port)

		// fmt.Println(currentSlaveNode.Address + ":" + strconv.Itoa(port))
		_, err2 := rpc.DialHTTP("tcp", currentSlaveNode.Address+":"+currentSlaveNode.Port)
		if err2 != nil {
			err2 = rpc.ServerError("Connection refused to slave node " + currentSlaveNode.Address + ":" + strconv.Itoa(port))
		} else {
			port = port + 1
			currentSlaveNode.Port = strconv.Itoa(port)
			fmt.Println(currentSlaveNode)
			*slaveNode = currentSlaveNode
		}
		currentSlaveServerIndex = (currentSlaveServerIndex + 1) % len(SlaveList)
	} else {
		*err = rpc.ServerError("No server has been started yet")
	}
}
