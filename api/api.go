package api

import (
	"blogServer/network"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Server comment
type Server struct {
	MasterServer *network.ServerMaster
}

func (server *Server) getAvailableServer(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("here")
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var slaveNode network.ServerNode
	var err error
	server.MasterServer.GetAvailableServer(&slaveNode, &err)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&slaveNode)
	}

}

func enableCors(w *http.ResponseWriter) {

	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// StartAPI comment
func StartAPI(masterNode *network.ServerMaster, address string, port int) {
	server := Server{MasterServer: masterNode}
	r := mux.NewRouter()
	api := r.PathPrefix("").Subrouter()

	api.HandleFunc("/getAvailableServer/", server.getAvailableServer).Methods("GET")
	fmt.Println("Listening on ", address+":"+strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), r))
}
