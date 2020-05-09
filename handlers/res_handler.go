package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"project/foodDeliveringApp/apicontext"
	"project/foodDeliveringApp/config"
	"project/foodDeliveringApp/core"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// AddRestaurantHandler : handler for uploading restaurant info
func AddRestaurantHandler(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	if ctx.APIContext.RoleID != apicontext.AdminRole {
		log.Println("unauthorized user access, err:", errors.New("user access forbidden"))
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorUnauthorized, errors.New("user access forbidden"))
		return
	}

	//to-do: validate that the userID is correct

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read request body, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorNoContent, err)
		return
	}

	var restaurant RestaurantInfo
	if err := json.Unmarshal(body, &restaurant); err != nil {
		log.Println("restaurant json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDecodingPayload, err)
		return
	}

	//to-do: validate req fields in front-end

	err = restaurant.ValidateAddRestaurantPayload()
	if err != nil {
		log.Println("restaurant payload validation failed, err:", err)
		if strings.Contains(err.Error(), "already exists") {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDocumentExists, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		}

		return
	}

	//to-do: set default values
	if restaurant.OpenTime.IsZero() {
		restaurant.OpenTime = time.Time{}
	} else if restaurant.CloseTime.IsZero() {
		restaurant.CloseTime = time.Time{}
	} else if len(restaurant.Status) == 0 {
		restaurant.Status = "not-active"
	}

	err = restaurant.SaveRestaurant()
	if err != nil {
		log.Println("restaurant info uploading failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorInsertingMongoDoc, err)
		return
	}

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("restaurant registration successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "restaurant registered successfully", nil)
}

// UpdateRestaurantName : handler to update restaurant name
func UpdateRestaurantName(w http.ResponseWriter, r *http.Request) { //to-do: not resolved yet

	ctx := apicontext.UpgradeCtx(r.Context())

	if ctx.APIContext.RoleID != apicontext.AdminRole {
		log.Println("unauthorized user access, err:", errors.New("user access forbidden"))
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorUnauthorized, errors.New("user access forbidden"))
		return
	}

	var updateRes RestaurantInfo

	if err := json.NewDecoder(r.Body).Decode(&updateRes); err != nil {
		log.Println("updateRestaurant json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDecodingPayload, err)
		return
	}

	err := updateRes.UpdateRestaurant()
	if err != nil {
		log.Println("restaurant updation failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDocumentNotFound, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		}

		return
	}

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("restaurant updation successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "restaurant updated successfully", nil)
}

// DeleteRestaurant : handler to update restaurant info
func DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeCtx(r.Context())

	if ctx.APIContext.RoleID != apicontext.AdminRole {
		log.Println("unauthorized user access, err:", errors.New("user access forbidden"))
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorUnauthorized, errors.New("user access forbidden"))
		return
	}

	res := RestaurantInfo{}
	res.ResID = ctx.RestaurantID
	res.BrID = ctx.BranchID

	err := res.DeleteRestaurant()
	if err != nil {
		log.Println("restaurant deletion failed, err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDocumentNotFound, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDeletingMongoDoc, err)
		}

		return
	}

	header := w.Header()
	SetLoginHeader(ctx, &header)

	log.Println("restaurant deleted successful")
	core.HTTPResponse(ctx, w, http.StatusOK, "restaurant deleted successfully", nil)
}

// GetAllRestaurant : handler to fetch all restaurant names
func GetAllRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeCtx(r.Context())

	restaurants, err := FetchAllRestaurants()
	if err != nil {
		log.Println("failed to fetch all restaurants")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		return
	}

	if len(restaurants) == 0 {
		log.Println("no restaurants found")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorNoContent, err)
		return
	}

	log.Println("restaurants fetched successfully")
	core.HTTPResponse(ctx, w, http.StatusOK, "all restaurants fetched successfully", restaurants)
}

// AddCusineHandler : handler to add cusine for a specific restaurant
func AddCusineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeCtx(r.Context())

	if ctx.APIContext.RoleID != apicontext.AdminRole {
		log.Println("unauthorized user access, err:", errors.New("user access forbidden"))
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorUnauthorized, errors.New("user access forbidden"))
		return
	}

	var cusReq Cusine
	if err := json.NewDecoder(r.Body).Decode(&cusReq); err != nil {
		log.Println("addCusine request decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDecodingPayload, err)
		return
	}

	err := cusReq.AddCusine(ctx.RestaurantID, ctx.BranchID)
	if err != nil {
		log.Println("failed to add cusine to restaurantID:", ctx.RestaurantID, "branchID:", ctx.BranchID, "err:", err)
		if strings.Contains(err.Error(), "already exists") {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDocumentExists, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		}

		return
	}

	log.Println("cusine added successfully to restaurantID:", ctx.RestaurantID, "branchID:", ctx.BranchID)
	core.HTTPResponse(ctx, w, http.StatusOK, "cusine added successfully", nil)
}

// UpdateCusineHandler : handler to update the cusine of a particular restaurant
func UpdateCusineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeCtx(r.Context())

	if ctx.APIContext.RoleID != apicontext.AdminRole {
		log.Println("unauthorized user access, err:", errors.New("user access forbidden"))
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorUnauthorized, errors.New("user access forbidden"))
		return
	}

	//to-do: validate that the cusineName is correct

	var cusine Cusine

	if err := json.NewDecoder(r.Body).Decode(&cusine); err != nil {
		log.Println("cusine json decode failed, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDecodingPayload, err)
		return
	}

	err := cusine.UpdateCusine(ctx.RestaurantID, ctx.BranchID)
	if err != nil {
		log.Println("failed to update cusine to restaurantID:", ctx.RestaurantID, "branchID:", ctx.BranchID, "err:", err)
		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDocumentNotFound, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		}

		return
	}

	log.Println("cusine updated successfully to restaurantID:", ctx.RestaurantID, "branchID:", ctx.BranchID)
	core.HTTPResponse(ctx, w, http.StatusOK, "cusine updated successfully", nil)
}

// DeleteCusineHandler : handler to delete a cusine from a particular res
func DeleteCusineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeCtx(r.Context())

	if ctx.APIContext.RoleID != apicontext.AdminRole {
		log.Println("unauthorized user access, err:", errors.New("user access forbidden"))
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorUnauthorized, errors.New("user access forbidden"))
		return
	}

	cusineName := mux.Vars(r)["cusineName"]
	if len(cusineName) == 0 {
		log.Println("cusineName param missing from URL")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorWithQueryParams, errors.New("cusineName missing from request URL"))
		return
	}

	cusine := Cusine{CusName: cusineName}
	err := cusine.DeleteCusine(ctx.RestaurantID, ctx.BranchID)
	if err != nil {
		log.Println("failed to delete cusine from restaurantID:", ctx.RestaurantID, "branchID:", ctx.BranchID, "err:", err)

		if strings.Contains(err.Error(), "not found") {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorDocumentNotFound, err)
		} else {
			core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		}

		return
	}

	log.Println("cusine deleted successfully from restaurantID:", ctx.RestaurantID, "branchID:", ctx.BranchID)
	core.HTTPResponse(ctx, w, http.StatusOK, "cusine deleted successfully", nil)
}

// GetAllCusinePerBranch : handler to fetch all cusine types of a particular branch
func GetAllCusinePerBranch(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	cusineNames, err := FetchCusineNamesPerBranch(ctx.APIContext.RestaurantID, ctx.APIContext.BranchID)
	if err != nil {
		log.Println("failed to fetch all cusine names, err:", err)
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		return
	}

	if len(cusineNames) == 0 {
		log.Println("no cusines found")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorNoContent, err)
		return
	}

	log.Println("cusine names fetched successfully")
	core.HTTPResponse(ctx, w, http.StatusOK, "all cusine names fetched successfully", cusineNames)
}

// GetFoodItemPerCusinePerBranch : handler to fetch all food items from one cusine type of a particular branch
func GetFoodItemPerCusinePerBranch(w http.ResponseWriter, r *http.Request) {

	ctx := apicontext.UpgradeCtx(r.Context())

	//to-do: validate that the cusineName is correct

	cusineName := mux.Vars(r)["cusineName"]
	if len(cusineName) == 0 {
		log.Println("cusineName param missing from URL")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorWithQueryParams, errors.New("cusineName missing from request URL"))
		return
	}

	foodItems, err := FetchFoodItemPerCusPerBr(ctx.RestaurantID, ctx.BranchID, cusineName)
	if err != nil {
		log.Println("failed to fetch all food items")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorQueryingDB, err)
		return
	}

	if len(foodItems) == 0 {
		log.Println("no fooditems found")
		core.HTTPErrorResponse(ctx, w, config.RestaurantModule, core.ErrorNoContent, err)
		return
	}

	log.Println("food items fetched successfully")
	core.HTTPResponse(ctx, w, http.StatusOK, "all cusine names fetched successfully", foodItems)
}

// UnmarshalJSON useful to keep a check on json values and also add default values
// func (br *Branch) UnmarshalJSON(data []byte) error {
// 	fmt.Println("start")
// 	var m map[string]interface{}

// 	err := json.Unmarshal(data, &m)
// 	if err != nil {
// 		fmt.Println("unmarshall", err)
// 		return err
// 	}

// 	fmt.Println(reflect.ValueOf(m))

// 	branch, ok := m["res_branch"].(map[string]interface{})
// 	if !ok {
// 		return errors.New("res_branch not found") // this code is redundant
// 	}

// 	for k, v := range branch {

// 		if k == "branch_name" {
// 			val, ok := v.(string)
// 			if !ok {
// 				fmt.Println("branch name")
// 				return errors.New("branch_name should be string")
// 			}
// 			br.BrName = val
// 		}

// 		if k == "branch_opentime" {
// 			val, ok := v.(string)
// 			if !ok {
// 				fmt.Println("branch opentime")
// 				return errors.New("branch_opentime should be in string")
// 			}

// 			if len(val) == 0 {
// 				t, err := time.Parse("15:04:05", "00:00:00")
// 				if err != nil {
// 					fmt.Println("opentime length zero")
// 					return err
// 				}
// 				br.OpenTime = t
// 			} else {
// 				t, err := time.Parse("15:04:05", val)
// 				if err != nil {
// 					fmt.Println("open time parse")
// 					return err
// 				}
// 				br.OpenTime = t
// 			}
// 		}

// 		if k == "branch_closetime" {
// 			val, ok := v.(string)
// 			if !ok {
// 				fmt.Println("branch closetime")
// 				return errors.New("branch_close should be in string")
// 			}

// 			if len(val) == 0 {
// 				t, err := time.Parse("15:04:05", "00:00:00")
// 				if err != nil {
// 					fmt.Println("closetime length zero")
// 					return err
// 				}
// 				br.OpenTime = t
// 			} else {
// 				t, err := time.Parse("15:04:05", val)
// 				if err != nil {
// 					fmt.Println("close time parse")
// 					return err
// 				}
// 				br.CloseTime = t
// 			}
// 		}

// 		if k == "branch_status" {
// 			val, ok := v.(string)
// 			if !ok {
// 				fmt.Println("branch status")
// 				return errors.New("branch_status should be string")
// 			}
// 			br.Status = val
// 		}

// 		if k == "branch_location" {
// 			val, ok := v.(string)
// 			if !ok {
// 				fmt.Println("branch location")
// 				return errors.New("branch_location should be string")
// 			}
// 			br.Location = val
// 		}

// 		if k == "branch_contactNo" {

// 			fmt.Println(reflect.TypeOf(branch[k]))
// 			val, ok := v.([]interface{})
// 			if !ok {
// 				fmt.Println("branch contact number")
// 				return errors.New("branch_contactNo should be string")
// 			}

// 			contact := make([]string, 0)

// 			for _, n := range val {
// 				cNo, ok := n.(string)
// 				if !ok {
// 					return errors.New("branch_contactNo should be array of string")
// 				}

// 				contact = append(contact, cNo)
// 			}

// 			br.ContactNumber = contact
// 		}

// 		if k == "branch_cusines" {
// 			val, ok := v.([]interface{})
// 			if !ok {
// 				fmt.Println("branch cusine")
// 				return errors.New("branch_cusines should be of an array of Cusine")
// 			}

// 			cusines := make([]Cusine, 0)
// 			for _, v := range val {
// 				cus, ok := v.(Cusine)
// 				if !ok {
// 					return errors.New("each element of branch_cusines should have Cusine struct")
// 				}

// 				cusines = append(cusines, cus)
// 			}
// 			br.Cusines = cusines
// 		}
// 	}
// 	return nil
// }
