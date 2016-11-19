package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Post defines the model of a Post
// swagger:model
type Post struct {
	// id for the Post
  // required: true
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	// title for the Post
  // required: true
	Title       string          `json:"title"`
	// association for the Post
  // required: true
	Association bson.ObjectId   `json:"association"`
	// description for the Post
  // required: true
	Description string          `json:"description"`
	// date for the Post
  // required: true
	Date        time.Time       `json:"date"`
	// likes for the Post
  // required: true
	Likes       []bson.ObjectId `json:"likes"`
	// Comment on the Post
  // required: true
	Comments    Comments        `json:"comments"`
	// Image for the Post
  // required: true
	Image    		string          `json:"image"`
	// Image size for the Post
  // required: true
	ImageSize		bson.M					`json:"imageSize"`
}

// Posts is an array of Post
// swagger:model
type Posts []Post


// AddPost will add the given post to the database
func AddPost(post Post) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	db.Insert(post)
	var result Post
	db.Find(bson.M{"title": post.Title, "date": post.Date}).One(&result)
	AddPostToAssociation(result.Association, result.ID)
	return result
}

// UpdatePost will update the post linked to the given ID,
// with the field of the given post, in the database
func UpdatePost(id bson.ObjectId, post Post) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"title"				:	post.Title,
		"description"	:	post.Description,
		"image"				:	post.Image,
		"imageSize"		:	post.ImageSize,
	}}
	db.Update(postID, change)
	var result Post
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DeletePost will delete the given post from the database
func DeletePost(post Post) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	db.RemoveId(post.ID)
	var result Post
	db.FindId(post.ID).One(result)
	DeleteNotificationsForPost(post.ID)
	RemovePostFromAssociation(post.Association, post.ID)
	for _, userId := range post.Likes{
		DislikePost(userId, post.ID)
	}
	return result
}

// GetPost will return an Post object from the given ID
func GetPost(id bson.ObjectId) Post {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	var result Post
	db.FindId(id).One(&result)
	return result
}

// GetLastestPosts will return an array of the last N Posts
func GetLastestPosts(number int) Posts {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	var result Posts
	db.Find(bson.M{}).Sort("-date").Limit(number).All(&result)
	return result
}

// LikePostWithUser will add the user to the list of
// user that liked the post (cf. Likes field)
func LikePostWithUser(id bson.ObjectId, userID bson.ObjectId) (Post, User) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"likes": userID,
	}}
	db.Update(postID, change)
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	user := LikePost(userID, post.ID)
	return post, user
}

// DislikePostWithUser will remove the user to the list of
// users that liked the post (cf. Likes field)
func DislikePostWithUser(id bson.ObjectId, userID bson.ObjectId) (Post, User) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("post")
	postID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"likes": userID,
	}}
	db.Update(postID, change)
	var post Post
	db.Find(bson.M{"_id": id}).One(&post)
	user := DislikePost(userID, post.ID)
	return post, user
}
