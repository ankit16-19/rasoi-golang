package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"

	. "github.com/ankit16-19/rasoi/dao"
	. "github.com/ankit16-19/rasoi/models"
)

var mdao = MenuDAO{}

// GetAllMenu :
func GetAllMenu(w http.ResponseWriter, r *http.Request) {
	menu, err := mdao.FindAll()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, menu)
}

// CreateMenu : add new Menu
func CreateMenu(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var menu Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	menu.ID = bson.NewObjectId()

	// set Date for every entry
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	dates := WholeWeekDates(time.Now().AddDate(0, 0, 7))

	for i := range days {
		reflect.ValueOf(&menu.MessUP).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
		reflect.ValueOf(&menu.MessDown).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
	}

	if err := mdao.Insert(menu); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, menu)
}

// UpdateMenu :
func UpdateMenu(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var menu Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// set Date for every entry
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	dates := WholeWeekDates(time.Now().AddDate(0, 0, 7))

	for i := range days {
		reflect.ValueOf(&menu.MessUP).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
		reflect.ValueOf(&menu.MessDown).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
	}
	if err := mdao.Update(menu); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// UpdateMenuDateIfWeekChange :
func UpdateMenuDateIfWeekChange() error {
	menus, err := mdao.FindAll()
	if err != nil {
		return err
	}
	// compare last sunday date with curretn date
	if time.Now().AddDate(0, 0, 7).After(menus[0].MessUP.Sun.Date) {
		fmt.Print("inside if")
		// set Date for every entry
		days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
		dates := WholeWeekDates(time.Now().AddDate(0, 0, 7))

		for i := range days {
			reflect.ValueOf(&menus[0].MessUP).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
			reflect.ValueOf(&menus[0].MessDown).Elem().FieldByName(days[i]).FieldByName("Date").Set(reflect.ValueOf(dates[i]))
		}
		if err := mdao.Update(menus[0]); err != nil {
			return err
		}
	}
	fmt.Print("UpdateMenuDateIfWeekChange function runnin....... exiting")
	return nil
}

func init() {
	mdao.Collection = "menu"
}
