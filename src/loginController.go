package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
  "os/exec"
	"strings"

	"github.com/freehaha/token-auth/memory"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Login struct {
	Username 		string 	`json:"username"`
	Password 		string 	`json:"password"`
	Device 			string 	`json:"device"`
	EraseDevice bool 		`json:"erase"`
}

type Credentials struct {
	ID 				bson.ObjectId	`bson:"_id,omitempty"`
	Username 	string				`json:"username"`
	AuthToken string				`json:"authtoken"`
	User 			bson.ObjectId	`json:"user" bson:"user"`
	Device 		string	 			`json:"device"`
}

type AssociationUser struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Username    string        `json:"username"`
	Association bson.ObjectId `json:"association" bson:"association"`
	Password    string        `json:"password"`
	Master      bool          `json:"master"`
	Owner       bson.ObjectId `json:"owner" bson:"owner,omitempty"`
}

func LogAssociationController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login Login
	decoder.Decode(&login)
	auth, master, err := checkLoginForAssociation(login)
	if err == nil {
		sessionToken := logAssociation(auth, master)
		json.NewEncoder(w).Encode(bson.M{"token": sessionToken.Token, "master": master, "associationID": auth})
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": "Failed to authentificate"})
	}
}

func LogUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var credentials Credentials
	decoder.Decode(&credentials)
	cred, err := checkLoginForUser(credentials)
	if err == nil {
		sessionToken := logUser(cred.User)
		user := GetUser(cred.User)
		json.NewEncoder(w).Encode(bson.M{"credentials": credentials, "sessionToken": sessionToken, "user": user})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"error": err})
	}
}

func SignInUserController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var login Login
	decoder.Decode(&login)

	if login.Username == "fthomasm" {
		login.Username = "fthomasm" + RandomString(4)
	}

	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(bson.M{"error": "De manière temporaire, les inscriptions sont désactivées. Réessaye Lundi 😊" })
	return

	isValid, err := verifyUser(login)
	if isValid {
		session, _ := mgo.Dial("127.0.0.1")
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		db := session.DB("insapp").C("user")
		count, _ := db.Find(bson.M{"username": login.Username}).Count()
		var user User
		if count == 0 {
			user = AddUser(User{Name: "", Username: login.Username, Description: "", Email: "", EmailPublic: false, Promotion: "", Events: []bson.ObjectId{}, PostsLiked: []bson.ObjectId{}})
		}else{
			db.Find(bson.M{"username": login.Username}).One(&user)
		}
		token := generateAuthToken()
		credentials := Credentials{AuthToken: token, User: user.ID, Username: user.Username, Device: login.Device}
		result := addCredentials(credentials)
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(bson.M{"error": err})
	}
}

func generateAuthToken() (string){
	out, _ := exec.Command("uuidgen").Output()
	return strings.TrimSpace(string(out))
}

func DeleteCredentialsForUser(id bson.ObjectId){
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("credentials")
	db.Remove(bson.M{"user": id})
}

func addCredentials(credentials Credentials) (Credentials){
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("credentials")
	var cred Credentials
	db.Find(bson.M{"username": credentials.Username}).One(&cred)
	db.RemoveId(cred.ID)
	db.Insert(credentials)
	var result Credentials
	db.Find(bson.M{"username": credentials.Username}).One(&result)
	return result
}

func checkLoginForAssociation(login Login) (bson.ObjectId, bool, error) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("association_user")
	var result []AssociationUser
	db.Find(bson.M{"username": login.Username, "password": GetMD5Hash(login.Password)}).All(&result)
	if len(result) > 0 {
		return result[0].Association, result[0].Master, nil
	}
	return bson.ObjectId(""), false, errors.New("Failed to authentificate")
}

func verifyUser(login Login) (bool, error){
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("user")
	count, err := db.Find(bson.M{"username": login.Username}).Count()
	if count > 0 || err != nil {
		return false || login.EraseDevice, errors.New("User Already Exist")
	}
	return true, nil
}

func checkLoginForUser(credentials Credentials) (Credentials, error) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("credentials")
	var result []Credentials
	db.Find(bson.M{"username": credentials.Username, "authtoken": credentials.AuthToken}).All(&result)
	if len(result) > 0 {
		return result[0], nil
	}
	return Credentials{}, errors.New("Wrong Credentials")
}

func logAssociation(id bson.ObjectId, master bool) *memstore.MemoryToken {
	if master {
		memStoreUser.NewToken(id.Hex())
		memStoreAssociationUser.NewToken(id.Hex())
		return memStoreSuperUser.NewToken(id.Hex())
	}
	memStoreUser.NewToken(id.Hex())
	return memStoreAssociationUser.NewToken(id.Hex())
}

func logUser(id bson.ObjectId) *memstore.MemoryToken {
	return memStoreUser.NewToken(id.Hex())
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
