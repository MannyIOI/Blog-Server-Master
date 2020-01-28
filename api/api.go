package api

import (
	"blogServer/models"
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
	Database   *models.DBHandler
	ServerNode network.ServerMaster
}

func (server Server) getUser(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	username := ""
	// var err error
	if val, ok := pathParams["username"]; ok {
		username = val
	}

	w.WriteHeader(http.StatusOK)

	var reply models.User
	server.Database.GetUser(username, &reply)
	json.NewEncoder(w).Encode(&reply)
}

func (server Server) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	var reply models.User
	server.ServerNode.NotifyNodesUser(user, &reply)
}

func (server Server) createBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var blog models.Blog
	_ = json.NewDecoder(r.Body).Decode(&blog)
	var reply models.Blog
	server.ServerNode.NotifyNodesBlogCreate(blog, &reply)
}

func (server Server) getBlog(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	blogID := ""
	// var err error
	if val, ok := pathParams["blogIdentifier"]; ok {
		blogID = val
	}

	w.WriteHeader(http.StatusOK)

	val, err := strconv.Atoi(blogID)
	if err == nil {
		var reply models.Blog
		server.Database.GetBlog(val, &reply)
		json.NewEncoder(w).Encode(&reply)
	}

}

func (server Server) getAllBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var reply []models.Blog
	server.Database.GetAllBlogs(&reply)
	json.NewEncoder(w).Encode(&reply)
}

// StartAPI comment
func StartAPI(db models.DBHandler, address string, port int) {
	server := Server{Database: &db}
	r := mux.NewRouter()
	api := r.PathPrefix("").Subrouter()

	api.HandleFunc("/user/{username}/", server.getUser).Methods("GET")
	api.HandleFunc("/user/", server.createUser).Methods("POST")

	api.HandleFunc("/blog/{blogIdentifier}/", server.getBlog).Methods("GET")
	api.HandleFunc("/blog/", server.createBlog).Methods("POST")

	api.HandleFunc("/blogs/", server.getAllBlogs).Methods("GET")

	fmt.Println("Listening on ", address+":"+strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), r))
}
