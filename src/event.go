package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Event represents the model of an event
// swagger:model Event
type Event struct {
	// id for the event
  // required: true
	ID           	bson.ObjectId   `bson:"_id,omitempty"`
	// name for the event
  // required: true
	Name         	string          `json:"name"`
	// association for the event
  // required: true
	Association  	bson.ObjectId   `json:"association" bson:"association"`
	// description for the event
  // required: true
	Description  	string          `json:"description"`
	// participants for the event
  // required: true
	Participants 	[]bson.ObjectId `json:"participants" bson:"participants,omitempty"`
	// status for the event
  // required: true
	Status       	string          `json:"status"`
	// palette for the event
  // required: true
	Palette			 	[][]int				 	`json:"palette"`
	// selectedcolor for the event
  // required: true
	SelectedColor int						 	`json:"selectedcolor"`
	// Start date of the event
  // required: true
	DateStart    	time.Time       `json:"dateStart"`
	// End date of the event
  // required: true
	DateEnd      	time.Time       `json:"dateEnd"`
	// Image for the event
  // required: true
	Image     	 	string          `json:"image"`
	// Background color of the event
  // required: true
	BgColor      	string          `json:"bgColor"`
	// Foreground color of the event
  // required: true
	FgColor      	string          `json:"fgColor"`
}

// Events is an array of Event
type Events []Event

// GetEvent returns an Event object from the given ID
func GetEvent(id bson.ObjectId) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Event
	db.FindId(id).One(&result)
	return result
}

// GetFutureEvents returns an array of Event objects
// that will happen after "NOW"
func GetFutureEvents() Events {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	var result Events
	var now = time.Now()
	db.Find(bson.M{"dateend": bson.M{"$gt": now}}).All(&result)
	return result
}

// AddEvent will add the Event event to the database
func AddEvent(event Event) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	db.Insert(event)
	var result Event
	db.Find(bson.M{"name": event.Name, "datestart": event.DateStart}).One(&result)
	AddEventToAssociation(result.Association, result.ID)
	return result
}

// UpdateEvent will update the Event event in the database
func UpdateEvent(id bson.ObjectId, event Event) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$set": bson.M{
		"name"					:	event.Name,
		"description"		:	event.Description,
		"status"				:	event.Status,
		"image"					:	event.Image,
		"palette"				:	event.Palette,
		"selectedcolor"	:	event.SelectedColor,
		"datestart"			:	event.DateStart,
		"dateend"				:	event.DateEnd,
		"bgcolor"				:	event.BgColor,
		"fgcolor"				: event.FgColor,
	}}
	db.Update(eventID, change)
	var result Event
	db.Find(bson.M{"_id": id}).One(&result)
	return result
}

// DeleteEvent will delete the given Event
func DeleteEvent(event Event) Event {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	db.Remove(event)
	DeleteNotificationsForEvent(event.ID)
	RemoveEventFromAssociation(event.Association, event.ID)
	for _, userId := range event.Participants{
		RemoveEventFromUser(userId, event.ID)
	}
	var result Event
	db.Find(event.ID).One(result)
	return result
}

// AddParticipant add the given userID to the given eventID as a participant
func AddParticipant(id bson.ObjectId, userID bson.ObjectId) (Event, User) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$addToSet": bson.M{
		"participants": userID,
	}}
	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)
	user := AddEventToUser(userID, event.ID)
	return event, user
}

// RemoveParticipant remove the given userID from the given eventID as a participant
func RemoveParticipant(id bson.ObjectId, userID bson.ObjectId) (Event, User) {
	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("event")
	eventID := bson.M{"_id": id}
	change := bson.M{"$pull": bson.M{
		"participants": userID,
	}}
	db.Update(eventID, change)
	var event Event
	db.Find(bson.M{"_id": id}).One(&event)
	user := RemoveEventFromUser(userID, event.ID)
	return event, user
}
