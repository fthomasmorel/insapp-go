package main

import (
	"net/http"

	"github.com/freehaha/token-auth"
	"github.com/freehaha/token-auth/memory"
	"github.com/gorilla/mux"
)

// Route type is used to define a route of the API
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes type is an array of Route
type Routes []Route

// NewRouter is the constructeur of the Router
// It will create every routes from the routes variable just above
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range publicRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	for _, route := range associationRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthAssociationUser.HandleFunc(route.HandlerFunc))
	}

	for _, route := range superRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthSuperUser.HandleFunc(route.HandlerFunc))
	}

	for _, route := range userRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(tokenAuthUser.HandleFunc(route.HandlerFunc))
	}

	return router
}

var tokenAuthAssociationUser = tauth.NewTokenAuth(nil, nil, memStoreAssociationUser, nil)
var tokenAuthSuperUser = tauth.NewTokenAuth(nil, nil, memStoreSuperUser, nil)
var tokenAuthUser = tauth.NewTokenAuth(nil, nil, memStoreUser, nil)

var memStoreAssociationUser = memstore.New("associationUser")
var memStoreSuperUser = memstore.New("superUser")
var memStoreUser = memstore.New("user")

var publicRoutes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"Credit", "GET", "/credit", Credit},
	Route{"Legal", "GET", "/legal", Legal},
	Route{"LogAssociation", "POST", "/login/association", LogAssociationController},
	Route{"LogUser", "POST", "/login/user", LogUserController},
	Route{"SignUser", "POST", "/signin/user/{ticket}", SignInUserController},
}

var superRoutes = Routes{
	Route{"GetUsers", "GET", "/user", GetAllUserController},
	Route{"AddAssociation", "POST", "/association", AddAssociationController},
	Route{"DeleteAssociation", "DELETE", "/association/{id}", DeleteAssociationController},
	Route{"CreateUserForAssociation", "POST", "/association/{id}/user", CreateUserForAssociationController},
	Route{"GetMyAssociations", "GET", "/association/{id}/myassociations", GetMyAssociationController},
}

var associationRoutes = Routes{
	//ASSOCIATIONS
	Route{"UpdateAssociation", "PUT", "/association/{id}", UpdateAssociationController},

	//EVENTS
	Route{"AddEvent", "POST", "/event", AddEventController},
	Route{"UpdateEvent", "PUT", "/event/{id}", UpdateEventController},
	Route{"DeleteEvent", "DELETE", "/event/{id}", DeleteEventController},

	//POSTS
	Route{"AddPost", "POST", "/post", AddPostController},
	Route{"UpdatePost", "PUT", "/post/{id}", UpdatePostController},
	Route{"DeletePost", "DELETE", "/post/{id}", DeletePostController},

	//Image
	//DEPENDENCIES : https://github.com/fengsp/color-thief-py
	Route{"UploadNewImage", "POST", "/image", UploadNewImageController},
	Route{"UploadImage", "POST", "/image/{name}", UploadImageController},
}

var userRoutes = Routes{
	//ASSOCIATIONS
	// swagger:route GET /association GetAllAssociation
	    // Get all the associations.
	    //     Responses:
	    //       default: Unknown
	    //       200: Associations
			//	     403: 403
	    //       406: 406
	Route{"GetAllAssociation", "GET", "/association", GetAllAssociationsController},
	// swagger:route GET /association/{id} GetAssociation
	    // Get the associations associated to the given id.
	    //     Responses:
	    //       default: Unknown
	    //       200: Association
			//	     403: 403
	    //       406: 406
	Route{"GetAssociation", "GET", "/association/{id}", GetAssociationController},


	//EVENTS
	// swagger:route GET /event GetFutureEvents
	    // Get the upcoming events
	    //     Responses:
	    //       default: Unknown
	    //       200: Events
			//	     403: 403
	    //       406: 406
	Route{"GetFutureEvents", "GET", "/event", GetFutureEventsController},
	// swagger:route GET /event/{id} GetEvent
	    // Get the event matching the given id
	    //     Responses:
	    //       default: Unknown
	    //       200: Event
			//	     403: 403
	    //       406: 406
	Route{"GetEvent", "GET", "/event/{id}", GetEventController},
	// swagger:route POST /event/{id}/participant/{userID} AddParticipant
	    // Add the participant userID to this event id
	    //     Responses:
	    //       default: Unknown
							//TODO: trouver autre chose : struct?
	    //       200: EventUser
			//	     403: 403
	    //       406: 406
	Route{"AddParticipant", "POST", "/event/{id}/participant/{userID}", AddParticipantController},
	// swagger:route DELETE /event/{id}/participant/{userID} RemoveParticipant
	    // Remove the participant userID to this event id
	    //     Responses:
	    //       default: Unknown
					  	//TODO: trouver autre chose : struct?
	    //       200: EventUser
			//	     403: 403
	    //       406: 406
	Route{"RemoveParticipant", "DELETE", "/event/{id}/participant/{userID}", RemoveParticipantController},


	//POSTS
	// swagger:route GET /post/{id} GetPost
	    // Get the post matching the given id
	    //     Responses:
	    //       default: Unknown
	    //       200: Event
			//	     403: 403
	    //       406: 406
	Route{"GetPost", "GET", "/post/{id}", GetPostController},
	Route{"GetLastestPost", "GET", "/post", GetLastestPostsController},
	Route{"LikePost", "POST", "/post/{id}/like/{userID}", LikePostController},
	Route{"DislikePost", "DELETE", "/post/{id}/like/{userID}", DislikePostController},
	Route{"CommentPost", "POST", "/post/{id}/comment", CommentPostController},
	Route{"UncommentPost", "DELETE", "/post/{id}/comment/{commentID}", UncommentPostController},
	Route{"ReportComment", "PUT", "/report/{id}/comment/{commentID}", ReportCommentController},


	// USER
	// swagger:route GET /user/{id} GetUser
	    // Return the user associated with the given id in the URL
	    //     Responses:
	    //       default: Unknown
	    //       200: User
			//	     403: 403
	    //       406: 406
	Route{"GetUser", "GET", "/user/{id}", GetUserController},

	// swagger:route PUT /user/{id} UpdateUser
	    // Update the user from the JSON body and return the modified user
	    //     Responses:
	    //       default: Unknown
	    //       200: User
			//	     403: 403
	    //       406: 406
	Route{"UpdateUser", "PUT", "/user/{id}", UpdateUserController},

	// swagger:route DELETE /user/{id} DeleteUser
	    // Delete the given user and answer an empty user if success
	    //     Responses:
	    //       default: Unknown
	    //       200: 200
			//	     403: 403
	    //       406: 406
	Route{"DeleteUser", "DELETE", "/user/{id}", DeleteUserController},

	// swagger:route PUT /search/users/{username} SearchUser
	    // Return the user(s) matching the entry (username)
	    //     Responses:
	    //       default: Unknown
	    //       200: 200
			//	     403: 403
	    //       406: 406
	Route{"SearchUser", "GET", "/search/users/{username}", SearchUserController},

	// swagger:route PUT /report/user/{id} ReportUser
	    // Report a user. Return empty JSON
	    //     Responses:
	    //       default: Unknown
	    //       200: Empty
			//	     403: 403
	    //       406: 406
	Route{"ReportUser", "PUT", "/report/user/{id}", ReportUserController},

	//NOTIFICATION
	Route{"Notification", "POST", "/notification", UpdateNotificationUserController},
	Route{"Notification", "GET", "/notification/{userID}", GetNotificationController},
	Route{"Notification", "DELETE", "/notification/{userID}/{id}", DeleteNotificationController},
}
