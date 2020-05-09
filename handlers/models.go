package handlers

import "time"

var SupportedImageExt = map[string]bool{
	"png":  true,
	"jpg":  true,
	"jpeg": true,
}

// UserInfo forms the basic structure of an user
type UserInfo struct {
	UserID     string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string   `json:"name" bson:"name"`
	Email      string   `json:"email" bson:"email"`
	Password   string   `json:"pswd" bson:"pswd"`
	ProfilePic string   `json:"profilePic" bson:"profilePic"` //S3 url
	Phone      []string `json:"phoneNo" bson:"phoneNo"`
	Address    string   `json:"address" bson:"address"`
}

// RestaurantInfo forms the struct to hold restaurant info
type RestaurantInfo struct {
	ResID         string    `json:"_id,omitempty" bson:"_id,omitempty"`
	BrID          string    `json:"branch_id" bson:"branch_id"`
	ResName       string    `json:"res_name" bson:"res_name"`
	BrName        string    `json:"branch_name" bson:"branch_name"`
	OpenTime      time.Time `json:"opentime" bson:"opentime"`
	CloseTime     time.Time `json:"closetime" bson:"closetime"`
	Status        string    `json:"status" bson:"status"`
	Location      string    `json:"location" bson:"location"`
	ContactNumber []string  `json:"contactNo" bson:"contactNo"`
	Cusines       []Cusine  `json:"cusines" bson:"cusines"`

	//to-do: add few more fields like userAgent, remoteAddr, timestamp, userID
}

// Cusine forms the struct to hold cusine details
type Cusine struct {
	CusName string `json:"cusine_name" bson:"cusine_name"`
	Items   []Item `json:"cusine_items" bson:"cusine_items"`
	//to-do: cusine type veg or non-veg
}

// Item forms the struct to hold per food item details
type Item struct {
	ItemName        string `json:"item_name" bson:"item_name"`
	ItemDescription string `json:"item_description" bson:"item_description"`
	ItemPrice       int    `json:"item_price" bson:"item_price"`
	ItemImage       string `json:"item_image" bson:"item_image"`
	IsAvailable     bool   `json:"item_isAvailable" bson:"item_isAvailable"`
	Offer           int    `json:"item_offer" bson:"item_offer"`
}
