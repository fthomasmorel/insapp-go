package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"github.com/freehaha/token-auth"
	"github.com/gorilla/mux"
)

func GetMyAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetMyAssociations(bson.ObjectIdHex(assocationID))
	json.NewEncoder(w).Encode(res)
}

// GetAssociationController will answer a JSON of the association
// linked to the given id in the URL
// @Title GetAssociationController
// @Description Return JSON of the association linked to the given id in the URL
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the association"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  Association		""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /association
// @Router /association/{id} [get]
func GetAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetAssociation(bson.ObjectIdHex(assocationID))
	json.NewEncoder(w).Encode(res)
}

// GetAllAssociationsController will answer a JSON of all associations
// @Title GetAllAssociationsController
// @Description Return JSON of all the associations
// @Accept  json
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {array}   Association		""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /association
// @Router /association [get]
func GetAllAssociationsController(w http.ResponseWriter, r *http.Request) {
	var res = GetAllAssociation()
	json.NewEncoder(w).Encode(res)
}

func CreateUserForAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assocationID := vars["id"]
	var res = GetAssociation(bson.ObjectIdHex(assocationID))

	decoder := json.NewDecoder(r.Body)
	var user AssociationUser
	decoder.Decode(&user)

	isValid := VerifyAssociationRequest(r, bson.ObjectIdHex(assocationID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	user.Association = res.ID
	user.Username = res.Email
	user.Password = GetMD5Hash(user.Password)
	AddAssociationUser(user)
	json.NewEncoder(w).Encode(res)
}

// AddAssociationController will answer a JSON of the
// brand new created association (from the JSON Body)
func AddAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var association Association
	decoder.Decode(&association)
	isValid := VerifyAssociationRequest(r, association.ID)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := AddAssociation(association)
	json.NewEncoder(w).Encode(res)
}

// UpdateAssociationController will answer the JSON of the
// modified association (from the JSON Body)
func UpdateAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var association Association
	decoder.Decode(&association)
	vars := mux.Vars(r)
	assocationID := vars["id"]
	isValid := VerifyAssociationRequest(r, bson.ObjectIdHex(assocationID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := UpdateAssociation(bson.ObjectIdHex(assocationID), association)
	json.NewEncoder(w).Encode(res)
}

// DeleteAssociationController will answer a JSON of an
// empty association if the deletation has succeed
func DeleteAssociationController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assoID := vars["id"]
	isValid := VerifyAssociationRequest(r, bson.ObjectIdHex(assoID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := DeleteAssociation(bson.ObjectIdHex(assoID))
	json.NewEncoder(w).Encode(res)
}

func VerifyAssociationRequest(r *http.Request, associationId bson.ObjectId) bool {
	token := tauth.Get(r)
	id := token.Claims("id").(string)
	if bson.ObjectIdHex(id) != associationId {
		result := GetAssociationUser(bson.ObjectIdHex(id))
		return result.Master
	}
	return true
}
