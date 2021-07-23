package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Blacklists handles requests for the /api/blacklists endpoint
func (as *Server) Blacklists(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		bs, err := models.GetBlacklists(ctx.Get(r, "user_id").(int64))
		if err != nil {
			log.Error(err)
		}
		JSONResponse(w, bs, http.StatusOK)
	//POST: Create a new blacklist and return it as JSON
	case r.Method == "POST":
		b := models.Blacklist{}
		// Put the request into a blacklist
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}
		// Check to make sure the name is unique
		_, err = models.GetPageByName(b.Name, ctx.Get(r, "user_id").(int64))
		if err != gorm.ErrRecordNotFound {
			JSONResponse(w, models.Response{Success: false, Message: "Blacklist name already in use"}, http.StatusConflict)
			log.Error(err)
			return
		}
		b.ModifiedDate = time.Now().UTC()
		b.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PostBlacklist(&b)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, b, http.StatusCreated)
	}
}

// Blacklist contains functions to handle the GET'ing, DELETE'ing and PUT'ing
// of a Blacklist object
func (as *Server) Blacklist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	b, err := models.GetBlacklist(id, ctx.Get(r, "user_id").(int64))
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Blacklist not found"}, http.StatusNotFound)
		return
	}
	switch {
	case r.Method == "GET":
		JSONResponse(w, b, http.StatusOK)
	case r.Method == "DELETE":
		err = models.DeleteBlacklist(id, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error deleting blacklist"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Blacklist Deleted Successfully"}, http.StatusOK)
	case r.Method == "PUT":
		b = models.Blacklist{}
		err = json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			log.Error(err)
		}
		if b.Id != id {
			JSONResponse(w, models.Response{Success: false, Message: "/:id and /:blacklist_id mismatch"}, http.StatusBadRequest)
			return
		}
		b.ModifiedDate = time.Now().UTC()
		b.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PutBlacklist(&b)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error updating blacklist: " + err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, b, http.StatusOK)
	}
}
