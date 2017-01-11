package main

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"
	"log"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	"github.com/freehaha/token-auth"
)

// GetPostController will answer a JSON of the post
// linked to the given id in the URL
// @Title GetPostController
// @Description Get the post matching the given id
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the post"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  Post		 ""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /post
// @Router /post/{id} [get]
func GetPostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	var res = GetPost(bson.ObjectIdHex(postID))
	json.NewEncoder(w).Encode(res)
}

// GetLastestPostsController will answer a JSON of the
// N lastest post. Here N = 50.
// @Title GetLastestPostsController
// @Description Get the N lastest post.
// @Accept  json
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {array}  Post		 ""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /post
// @Router /post [get]
func GetLastestPostsController(w http.ResponseWriter, r *http.Request) {
	var res = GetLastestPosts(50)
	json.NewEncoder(w).Encode(res)
}

// AddPostController will answer a JSON of the
// brand new created post (from the JSON Body)
func AddPostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	decoder.Decode(&post)
	post.Date = time.Now()

	isValid := VerifyAssociationRequest(r, post.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	res := AddPost(post)
	asso := GetAssociation(post.Association)
	json.NewEncoder(w).Encode(res)
	go TriggerNotificationForPost(asso.ID, res.ID, "@" + strings.ToLower(asso.Name) + " a posté une nouvelle news 📰")
}

// UpdatePostController will answer the JSON of the
// modified post (from the JSON Body)
func UpdatePostController(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var post Post
	decoder.Decode(&post)
	vars := mux.Vars(r)
	postID := vars["id"]

	isValid := VerifyAssociationRequest(r, post.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	res := UpdatePost(bson.ObjectIdHex(postID), post)
	json.NewEncoder(w).Encode(res)
}

// DeletePostController will answer a JSON of an
// empty post if the deletation has succeed
func DeletePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := GetPost(bson.ObjectIdHex(vars["id"]))

	isValid := VerifyAssociationRequest(r, post.Association)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	res := DeletePost(post)
	json.NewEncoder(w).Encode(res)
}

// LikePostController will answer a JSON of the
// post and the user that liked the post
// @Title LikePostController
// @Description Return JSON with the post and the user that liked the post
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the post"
// @Param   userID 	 		  path     bson.ObjectId        true        "id of the user"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}   bson.M		 ""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /post
// @Router /post/{id}/like/{userID} [post]
func LikePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := vars["userID"]
	isValid := VerifyUserRequest(r, bson.ObjectIdHex(userID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	post, user := LikePostWithUser(bson.ObjectIdHex(postID), bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(bson.M{"post": post, "user": user})
}

// DislikePostController will answer a JSON of the
// post and the user that disliked the post
// @Title DislikePostController
// @Description Return JSON with the post and the user that disliked the post
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the post"
// @Param   userID 	 		  path     bson.ObjectId        true        "id of the user"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  bson.M		 ""
// @Failure 403 {object}  error   	"Access forbidden"
// @Failure 406 {object}  error     "Request not accepted"
// @Resource /post
// @Router /post/{id}/like/{userID} [delete]
func DislikePostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	userID := vars["userID"]
	isValid := VerifyUserRequest(r, bson.ObjectIdHex(userID))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	post, user := DislikePostWithUser(bson.ObjectIdHex(postID), bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(bson.M{"post": post, "user": user})
}

// CommentPostController will answer a JSON of the post
// @Title CommentPostController
// @Description Return JSON of the post that has been commented
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the post"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  Post		 ""
// @Failure 403 {object}  error    "Access forbidden"
// @Failure 406 {object}  error    "Request not accepted"
// @Resource /post
// @Router /post/{id}/comment [post]
func CommentPostController(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var comment Comment
	if err := json.Unmarshal([]byte(string(body)), &comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{"error": "Mauvais Format"})
		return
	}

	isValid := VerifyUserRequest(r, comment.User)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}

	comment.ID = bson.NewObjectId()
	comment.Date = time.Now()

	vars := mux.Vars(r)
	postID := vars["id"]
	res := CommentPost(bson.ObjectIdHex(postID), comment)
	json.NewEncoder(w).Encode(res)

	for _, tag := range(comment.Tags){
		go TriggerNotificationForUser(comment.User, bson.ObjectIdHex(tag.User), res.ID , "@" + GetUser(comment.User).Username + " t'a taggé sur \"" + res.Title + "\"", comment)
	}
}

// UncommentPostController will answer a JSON of the post
// @Title UncommentPostController
// @Description Return JSON of the post that has been uncommented
// @Accept  json
// @Param   id 	 			    path     bson.ObjectId        true        "id of the post"
// @Param   commentID 	  path     bson.ObjectId        true        "id of the comment"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  Post		 ""
// @Failure 403 {object}  error    "Access forbidden"
// @Failure 406 {object}  error    "Request not accepted"
// @Resource /post
// @Router /post/{id}/comment/{commentID} [delete]
func UncommentPostController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	commentID := vars["commentID"]
	comment, err := GetComment(bson.ObjectIdHex(postID), bson.ObjectIdHex(commentID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Inexistant"})
		return
	}
	post := GetPost(bson.ObjectIdHex(postID))
	isUserValid := VerifyUserRequest(r, comment.User)
	isAssociationValid := VerifyAssociationRequest(r, post.Association)
	if !isUserValid && !isAssociationValid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"error": "Contenu Protégé"})
		return
	}
	res := UncommentPost(bson.ObjectIdHex(postID), bson.ObjectIdHex(commentID))
	json.NewEncoder(w).Encode(res)
}

// @Title ReportCommentController
// @Description Report a comment
// @Accept  json
// @Param   id 	 			path     bson.ObjectId        true        "id of the post"
// @Param   commentID 	  path     bson.ObjectId        true        "id of the comment"
// @Param   token         query    string     true        "#insapptoken"
// @Success 200 {object}  bson.M	 ""
// @Failure 403 {object}  error    "Access forbidden"
// @Failure 406 {object}  error    "Request not accepted"
// @Resource /report
// @Router /report/{id}/comment/{commentID} [put]
func ReportCommentController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	commentID := vars["commentID"]
	token := tauth.Get(r)
	userID := token.Claims("id").(string)
	ReportComment(bson.ObjectIdHex(postID), bson.ObjectIdHex(commentID), bson.ObjectIdHex(userID))
	json.NewEncoder(w).Encode(bson.M{})
}

// // AddImagePostController will set the image of the post and return the post
// func AddImagePostController(w http.ResponseWriter, r *http.Request) {
// 	fileName := UploadImage(r)
// 	if fileName == "error" {
// 		w.Header().Set("status", "400")
// 		fmt.Fprintln(w, "{}")
// 	} else {
// 		vars := mux.Vars(r)
// 		res := SetImagePost(bson.ObjectIdHex(vars["id"]), fileName)
// 		json.NewEncoder(w).Encode(res)
// 	}
// }
