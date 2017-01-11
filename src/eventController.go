package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// GetEventController will answer a JSON of the event
// from the given "id" in the URL. (cf Routes in routes.go)
// @Title GetEventController
// @Description Return JSON of the event linked to the association's id in the URL
// @Accept  json
// @Param   id 	 			    path      bson.ObjectId        true        "id of the association"
// @Param   token         query     string     true        "#insapptoken"
// @Success 200 {object}  Event	    ""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /event
// @Router /event/{id} [get]
func GetEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetEvent(bson.ObjectIdHex(assocationID))
	json.NewEncoder(w).Encode(res)
}

// GetFutureEventsController will answer a JSON
// containing all future events from "NOW"
// @Title GetFutureEventsController
// @Description Return JSON of all future events from "NOW"
// @Accept  json
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {array}   Event	   ""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /event
// @Router /event [get]
func GetFutureEventsController(w http.ResponseWriter, r *http.Request) {
	var res = GetFutureEvents()
	json.NewEncoder(w).Encode(res)
}

// AddEventController will answer the JSON
// of the brand new created event from the JSON body
func AddEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)

	isValid := VerifyAssociationRequest(r, event.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	res := AddEvent(event)
	asso := GetAssociation(event.Association)
	json.NewEncoder(w).Encode(res)
	go TriggerNotificationForEvent(asso.ID, res.ID, "@" + strings.ToLower(asso.Name) + " t'invite à " + res.Name + " 📅")
}

// UpdateEventController will answer the JSON
// of the brand new modified event from the JSON body
func UpdateEventController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	decoder.Decode(&event)
	vars := mux.Vars(r)
	eventID := vars["id"]

	isValid := VerifyAssociationRequest(r, event.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	res := UpdateEvent(bson.ObjectIdHex(eventID), event)
	json.NewEncoder(w).Encode(res)
}

// DeleteEventController will answer an empty JSON
// if the deletation has succeed
func DeleteEventController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := GetEvent(bson.ObjectIdHex(vars["id"]))

	isValid := VerifyAssociationRequest(r, event.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	res := DeleteEvent(event)
	json.NewEncoder(w).Encode(res)
}

// AddParticipantController will answer the JSON
// of the event with the given partipant added
// @Title AddParticipantController
// @Description Add the participant userID to this event id
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the event"
// @Param   userID 	 			path     bson.ObjectId        true        "id of the user"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  bson.M		""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /event
// @Router /event/{id}/participant/{userID} [post]
func AddParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	isValid := VerifyUserRequest(r, userID)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	event, user := AddParticipant(eventID, userID)
	json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}

// RemoveParticipantController will answer the JSON
// of the event without the given partipant added
// @Title RemoveParticipantController
// @Description Remove the participant userID to this event id
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the event"
// @Param   userID 	 			path     bson.ObjectId        true        "id of the user"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  bson.M		""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /event
// @Router /event/{id}/participant/{userID} [delete]
func RemoveParticipantController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := bson.ObjectIdHex(vars["id"])
	userID := bson.ObjectIdHex(vars["userID"])
	isValid := VerifyUserRequest(r, userID)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	event, user := RemoveParticipant(eventID, userID)
	json.NewEncoder(w).Encode(bson.M{"event": event, "user": user})
}
