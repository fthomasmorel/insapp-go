package main

import (
	"encoding/json"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
)


// AddUserController will answer a JSON of the
// brand new created user (from the JSON Body)

// @Title UpdateNotificationUserController
// @Description Update notification
// @Accept  json
// @Param   token     query    string     true        "#insapptoken"
// @Success 200 {object} bson.M
// @Failure 403 {object} error  "Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /notification
// @Router /notification [post]
func UpdateNotificationUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user NotificationUser
	decoder.Decode(&user)
	isValid := VerifyUserRequest(r, user.UserId)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	CreateOrUpdateNotificationUser(user)
	json.NewEncoder(w).Encode(bson.M{"status": "ok"})
}

// @Title GetNotificationController
// @Description Return the notifiaction for the user id
// @Accept  json
// @Param   userID    path     bson.ObjectId     true        "id of the user"
// @Param   token     query    string     			 true        "#insapptoken"
// @Success 200 {object} bson.M
// @Failure 403 {object} error    "Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /notification
// @Router /notification/{userID} [get]
func GetNotificationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	isValid := VerifyUserRequest(r, bson.ObjectIdHex(userID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := GetNotificationsForUser(bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(bson.M{"notifications": res})
}

// @Title DeleteNotificationController
// @Description Delete a notification
// @Accept  json
// @Param   userID    path     bson.ObjectId     true        "id of the user"
// @Param   id    	  path     bson.ObjectId     true        "id of the notification"
// @Param   token     query    string     			 true        "#insapptoken"
// @Success 200 {object} bson.M
// @Failure 403 {object} error    "Forbidden access"
// @Failure 406 {object} error    "Request not accepted"
// @Resource /notification
// @Router /notification/{userID}/{id} [delete]
func DeleteNotificationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	notifID := vars["id"]
	isValid := VerifyUserRequest(r, bson.ObjectIdHex(userID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := ReadNotificationForUser(bson.ObjectIdHex(userID), bson.ObjectIdHex(notifID))
	json.NewEncoder(w).Encode(bson.M{"notifications": res})
}
