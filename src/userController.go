package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"github.com/freehaha/token-auth"
	"github.com/gorilla/mux"
)


// @Title GetUserController
// @Description Return the user associated with the given id in the URL
// @Accept  json
// @Param   id     path    int     true        "#userid"
// @Success 200 {object} User
// @Failure 403 {object} error  "Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /user
// @Router /user/{id} [get]
func GetUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	var res = GetUser(bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(res)
}

func GetAllUserController(w http.ResponseWriter, r *http.Request) {
	var res = GetAllUser()
	json.NewEncoder(w).Encode(res)
}

// AddUserController will answer a JSON of the
// brand new created user (from the JSON Body)
func AddUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	decoder.Decode(&user)
	res := AddUser(user)
	json.NewEncoder(w).Encode(res)
}

// @Title UpdateUserController
// @Description Update the user from the JSON body and return the modified user
// @Accept  json
// @Param   id 	     path    int     true        "#userid"
// @Param   user     body    User     true        "Modification to give to the user"
// @Success 200 {object} User
// @Failure 403 {object} error  "Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /user
// @Router /user/{id} [put]
func UpdateUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	decoder.Decode(&user)
	vars := mux.Vars(r)
	userID := vars["id"]
	isValidUser := VerifyUserRequest(r, bson.ObjectIdHex(userID))
	isValidAssociation := VerifyAssociationRequest(r, bson.ObjectIdHex(userID))
	if !isValidUser && !isValidAssociation {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := UpdateUser(bson.ObjectIdHex(userID), user)
	json.NewEncoder(w).Encode(res)
}


// DeleteUserController will answer a JSON of an
// empty user if the deletion has succeed

// @Title DeleteUserController
// @Description Delete the given user and answer an empty user if success
// @Accept  json
// @Param   id 	     path    int     true        "#useridid of the user to delete"
// @Success 200 {object} User			"Return empty User when success"
// @Failure 403 {object} error  	"Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /user
// @Router /user/{id} [delete]
func DeleteUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	isUserValid := VerifyUserRequest(r, bson.ObjectIdHex(userID))
	isAssociationValid := VerifyAssociationRequest(r, bson.ObjectIdHex(userID))
	if !isUserValid && !isAssociationValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	user := GetUser(bson.ObjectIdHex(userID))
	res := DeleteUser(user)
	json.NewEncoder(w).Encode(res)
}

// @Title SearchUserController
// @Description Return the user(s) matching the entry (username)
// @Accept  json
// @Param   username 	   path    string     true        "username of the user to look for"
// @Success 200 {array}  User			"Return the list of User matching"
// @Failure 403 {object} error  	"Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /search
// @Router /search/users/{username} [get]
func SearchUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users := SearchUser(vars["username"])
	json.NewEncoder(w).Encode(bson.M{"users": users})
}

// @Title ReportUserController
// @Description Report a user. Return empty JSON
// @Accept  json
// @Param   id 	 			    path    int     true        "if of the user to report"
// @Success 200 {object}  bson.M		"Empty JSON"
// @Failure 403 {object}  error   	"Forbidden access"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /report
// @Router /report/user/{id} [put]
func ReportUserController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	token := tauth.Get(r)
	reporterID := token.Claims("id").(string)
	ReportUser(bson.ObjectIdHex(userID), bson.ObjectIdHex(reporterID))
	json.NewEncoder(w).Encode(bson.M{})
}

func VerifyUserRequest(r *http.Request, userId bson.ObjectId) bool {
	token := tauth.Get(r)
	id := token.Claims("id").(string)
	return bson.ObjectIdHex(id) == userId
}
