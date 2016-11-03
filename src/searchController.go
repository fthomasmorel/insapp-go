package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
  "github.com/gorilla/mux"
)

func SearchUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users := SearchUser(vars["name"])
	json.NewEncoder(w).Encode(bson.M{"users": users})
}

func SearchPostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	posts := SearchPost(vars["name"])
	json.NewEncoder(w).Encode(bson.M{"posts": posts})
}

func SearchEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	events := SearchEvent(vars["name"])
	json.NewEncoder(w).Encode(bson.M{"events": events})
}

func SearchAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assos := SearchAssociation(vars["name"])
	json.NewEncoder(w).Encode(bson.M{"associations": assos})
}

func SearchUniversalController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  users := SearchUser(vars["name"])
  posts := SearchPost(vars["name"])
  events := SearchEvent(vars["name"])
	assos := SearchPost(vars["name"])
	json.NewEncoder(w).Encode(bson.M{"associations": assos, "users": users, "posts": posts, "events": events})
}
