package serve

import (
	"net/http"
	"project/foodDeliveringApp/core"
	"project/foodDeliveringApp/handlers"
)

const (
	//Port is the port number of the current service
	Port = "8989"
	//SubRoute is the basic route of all APIs of this service
	SubRoute = "/foodApp"
)

func init() {

	//Adding ping api check for all the services for health check
	core.AddNoAuthRoutes(
		"Ping Check",
		http.MethodGet,
		"/health",
		core.Ping)
}

//StartRoutes is the starting point of the http service
func StartRoutes() {

	/*---------------------------------User Registration-------------------------------*/
	core.AddLoginRoutes(
		"User SignUp",
		http.MethodPost,
		"/register",
		nil,
		handlers.UserSignUp,
	)

	/*---------------------------------User Login-------------------------------*/
	core.AddLoginRoutes(
		"User Login",
		http.MethodPost,
		"/login",
		nil,
		handlers.UserLogin,
	)

	/*---------------------------------User Update Address-------------------------------*/
	core.AddUserRoute(
		"User Update Address",
		http.MethodPost,
		"/update/address",
		nil,
		handlers.UserUpdateAddress,
	)

	/*---------------------------------User Update Phone Number-------------------------------*/
	core.AddUserRoute(
		"User Update PhoneNumber",
		http.MethodPost,
		"/update/phone-num",
		nil,
		handlers.UserUpdatePhone,
	)

	/*---------------------------------User Update ProfilePic-------------------------------*/
	core.AddUserRoute(
		"User Update ProfilePic",
		http.MethodPost,
		"/update/profile-pic",
		nil,
		handlers.UserUpdateProfile,
	)

	/*---------------------------------User Update Password-------------------------------*/
	core.AddUserRoute(
		"User Update Password",
		http.MethodPost,
		"/update/password",
		nil,
		handlers.UserUpdatePassword,
	)

	/*---------------------------------User Logout-------------------------------*/
	core.AddUserRoute(
		"User Logout",
		http.MethodGet,
		"/logout",
		nil,
		handlers.UserLogout,
	)

	/*---------------------------------Add Restaurant-------------------------------*/
	core.AddUserRoute(
		"Register Restaurant",
		http.MethodPost,
		"/addrestaurant",
		nil,
		handlers.AddRestaurantHandler,
	)

	/*---------------------------------Update Restaurant Name-------------------------------*/
	core.AddRoute(
		"Update Restaurant Name",
		http.MethodPost,
		"/updateRestaurant/name",
		nil,
		handlers.UpdateRestaurantName,
	)

	/*---------------------------------Delete Restaurant-------------------------------*/
	core.AddRoute(
		"Delete Restaurant",
		http.MethodDelete,
		"/deleteRestaurant",
		nil,
		handlers.DeleteRestaurant,
	)

	/*---------------------------------Fetch All Restaurant-------------------------------*/
	core.AddUserRoute(
		"Fetch All Restaurant",
		http.MethodGet,
		"/getAllRestaurant",
		nil,
		handlers.GetAllRestaurant,
	)

	/*---------------------------------Add Cusine to a particular Restaurant-------------------------------*/
	core.AddRoute(
		"Add Cusine to Specific Restaurant",
		http.MethodPost,
		"/addCusine",
		nil,
		handlers.AddCusineHandler,
	)

	/*---------------------------------Update Cusine-------------------------------*/
	core.AddRoute(
		"Update Cusine",
		http.MethodPost,
		"/updateCusine",
		nil,
		handlers.UpdateCusineHandler,
	)

	/*---------------------------------Delete Cusine-------------------------------*/
	core.AddRoute(
		"Delete Cusine",
		http.MethodGet,
		"/deleteCusine/{cusineName}",
		nil,
		handlers.DeleteCusineHandler,
	)

	/*---------------------------------Fetch All Cusine Per Branch-------------------------------*/
	core.AddRoute(
		"Fetch All Cusine",
		http.MethodGet,
		"/getAllCusine",
		nil,
		handlers.GetAllCusinePerBranch,
	)

	/*---------------------------------Fetch Food Items Per Cusine Per Branch-------------------------------*/
	core.AddRoute(
		"Food Items Per Cusine Per Branch",
		http.MethodGet,
		"/foodPerCusine/{cusineName}",
		nil,
		handlers.GetFoodItemPerCusinePerBranch,
	)

	core.StartServer(Port, SubRoute)
}
