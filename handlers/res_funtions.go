package handlers

import (
	"errors"
	"project/foodDeliveringApp/mongodb"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// ValidateAddRestaurantPayload verifies restaurant payload and checks duplicacy based on restaurant name and branch name
func (res RestaurantInfo) ValidateAddRestaurantPayload() error {
	if len(res.ResName) == 0 {
		return errors.New("restaurant name not found")
	} else if len(res.BrName) == 0 {
		return errors.New("branch name not found")
	} else if len(res.Location) == 0 { //to-do: location implementation will change
		return errors.New("branch location not found")
	} else if len(res.ContactNumber) == 0 {
		return errors.New("branch contact number not found")
	}

	err := mongodb.Exists(mongodb.RestaurantCollection, bson.M{"res_name": res.ResName, "branch_name": res.BrName}, nil)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	} else if err == nil {
		return errors.New("restaurant already exists")
	}

	return nil
}

// SaveRestaurant makes restaurant entry into database
func (res *RestaurantInfo) SaveRestaurant() error {

	resID, err := mongodb.Count(mongodb.RestaurantCollection, nil)
	if err != nil {
		return err
	}

	brID, err := mongodb.Count(mongodb.RestaurantCollection, bson.M{"res_name": res.ResName})
	if err != nil {
		return err
	}

	//to-do: incorrect logic
	res.ResID = strconv.Itoa(resID + 1)
	res.BrID = strconv.Itoa(brID + 1)

	err = mongodb.Insert(mongodb.RestaurantCollection, *res)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRestaurant updates restaurant details to database
func (res RestaurantInfo) UpdateRestaurant() error {

	err := mongodb.Update(mongodb.RestaurantCollection, bson.M{"_id": res.ResID, "branch_id": res.BrID}, bson.M{"$set": bson.M{"res_name": res.ResName}})

	return err
}

// DeleteRestaurant deletes restaurant details from database
func (res RestaurantInfo) DeleteRestaurant() error {

	err := mongodb.HardDelete(mongodb.RestaurantCollection, bson.M{"_id": res.ResID, "branch_id": res.BrID})
	if err != nil {
		return err
	}

	return nil
}

// FetchAllRestaurants retrieves all restaurants
func FetchAllRestaurants() (map[string][]string, error) {

	var names []RestaurantInfo
	err := mongodb.ReadAll(mongodb.RestaurantCollection, nil, bson.M{"res_name": 1, "_id": 0, "branch_name": 1}, &names)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]string)

	for _, v := range names {
		m[v.ResName] = append(m[v.ResName], v.BrName)
	}

	return m, nil
}

// AddCusine stores cusine info in db
func (cusine Cusine) AddCusine(resID, brID string) error {

	existQuery := bson.M{"_id": resID, "branch_id": brID, "cusines.cusine_name": cusine.CusName}
	err := mongodb.Exists(mongodb.RestaurantCollection, existQuery, nil)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	} else if err == nil {
		return errors.New("restaurant already exists")
	}

	err = mongodb.Update(mongodb.RestaurantCollection, bson.M{"_id": resID, "branch_id": brID}, bson.M{"$addToSet": bson.M{"cusines": cusine}})
	return err
}

// UpdateCusine updates the cusine of a specific restaurant in db
func (cusine Cusine) UpdateCusine(resID, brID string) (err error) {

	query := bson.M{"_id": resID, "branch_id": brID, "cusines.cusine_name": cusine.CusName}
	update := bson.M{"$set": bson.M{"cusines.$": cusine}}

	err = mongodb.Update(mongodb.RestaurantCollection, query, update)

	return
}

// DeleteCusine removes the cusine entry from the specific restaurant in db
func (cusine Cusine) DeleteCusine(resID, brID string) (err error) {

	query := bson.M{"_id": resID, "branch_id": brID}
	update := bson.M{"$pull": bson.M{"cusines": bson.M{"cusine_name": cusine.CusName}}}

	err = mongodb.Update(mongodb.RestaurantCollection, query, update)

	return
}

// FetchCusineNamesPerBranch retrieves all cusine names per branch
func FetchCusineNamesPerBranch(resID, brID string) (allCusine []string, err error) {

	var names struct {
		Cusines []Cusine `json:"cusines"`
	}

	err = mongodb.ReadOne(mongodb.RestaurantCollection, bson.M{"_id": resID, "branch_id": brID}, bson.M{"cusines.cusine_name": 1, "_id": 0}, &names)
	if err != nil {
		return
	}

	for _, v := range names.Cusines {
		allCusine = append(allCusine, v.CusName)
	}

	return
}

// FetchFoodItemPerCusPerBr retrieves all food item per cusine per branch
func FetchFoodItemPerCusPerBr(resID, brID, cusineName string) ([]Item, error) {

	var items struct {
		Cusines []Cusine `json:"cusines"`
	}

	err := mongodb.ReadOne(mongodb.RestaurantCollection, bson.M{"_id": resID, "branch_id": brID, "cusines.cusine_name": cusineName}, bson.M{"cusines.$": 1, "_id": 0}, &items)
	if err != nil {
		return nil, err
	}

	i := items.Cusines[0].Items

	return i, nil
}
