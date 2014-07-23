package main

import (
	"net/http"
	"strconv"

	"fmt"
	"github.com/emicklei/go-restful"
)

type App struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type AppRequest struct {
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

func (a AppRequest) Validate() error {
	if a.Label == nil {
		return fmt.Errorf("Label must not be nil.")
	}
	if a.Description == nil {
		return fmt.Errorf("Description must not be nil.")
	}

	return nil
}

type AppResource struct {
	// normally one would use DAO (data access object)
	apps map[string]App
}

func (a AppResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/apps").
		Doc("Manage Apps").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{app-id}").To(a.findApp).
		// docs
		Doc("get an app").
		Operation("findApp").
		Param(ws.PathParameter("app-id", "identifier of the app").DataType("string")).
		Writes(App{})) // on the response

	ws.Route(ws.PUT("/{app-id}").To(a.updateApp).
		// docs
		Doc("update an app").
		Operation("updateApp").
		Param(ws.PathParameter("app-id", "identifier of the app").DataType("string")).
		Reads(AppRequest{})) // from the request

	ws.Route(ws.PATCH("/{app-id}").To(a.updateApp).
		// docs
		Doc("patch an app with partial request").
		Operation("updateApp").
		Param(ws.PathParameter("app-id", "identifier of the app").DataType("string")).
		Reads(AppRequest{})) // from the request

	ws.Route(ws.POST("").To(a.createApp).
		// docs
		Doc("create an app").
		Operation("createApp").
		Reads(AppRequest{})) // from the request

	ws.Route(ws.DELETE("/{app-id}").To(a.removeApp).
		// docs
		Doc("delete an app").
		Operation("removeApp").
		Param(ws.PathParameter("app-id", "identifier of the app").DataType("string")))

	container.Add(ws)
}

// GET http://localhost:8080/apps/1
//
func (a AppResource) findApp(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("app-id")
	app := a.apps[id]
	if len(app.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "App could not be found.")
		return
	}
	response.WriteEntity(app)
}

// POST http://localhost:8080/apps
func (a *AppResource) createApp(request *restful.Request, response *restful.Response) {
	app := new(App)
	err := request.ReadEntity(app)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	app.Id = strconv.Itoa(len(a.apps) + 1) // simple id generation
	a.apps[app.Id] = *app
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(app)
}

// PUT http://localhost:8080/apps/1
func (a *AppResource) updateApp(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("app-id")
	app := a.apps[id]
	if len(app.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "App could not be found.")
		return
	}
	appReq := new(AppRequest)
	err := request.ReadEntity(&appReq)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	if request.Request.Method == "PUT" {
		err = appReq.Validate()
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, err.Error())
			return
		}
	}

	app.Id = id
	if appReq.Label != nil {
		app.Label = *appReq.Label
	}
	if appReq.Description != nil {
		app.Description = *appReq.Description
	}

	a.apps[id] = app
	response.WriteEntity(app)
}

// DELETE http://localhost:8080/apps/1
//
func (a *AppResource) removeApp(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("app-id")
	delete(a.apps, id)
}
