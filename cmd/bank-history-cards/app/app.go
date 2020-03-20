package app

import (
	"bank-history-cards/pkg/core/auth"
	"bank-history-cards/pkg/core/history"
	"github.com/jafarsirojov/mux/pkg/mux"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
)

type MainServer struct {
	exactMux *mux.ExactMux
	cardsSvc *history.Service
}

func NewMainServer(exactMux *mux.ExactMux, cardsSvc *history.Service) *MainServer {
	return &MainServer{exactMux: exactMux, cardsSvc: cardsSvc}
}

func (m *MainServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.exactMux.ServeHTTP(writer, request)
}

func (m *MainServer) HandleGetAllShowOperationsLog(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	if authentication.Id != 0 {
		log.Printf("can't authentication is not admin, this id user = %d", authentication.Id)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print("all history")
	models, err := m.cardsSvc.All()
	if err != nil {
		log.Print("can't get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Print("can't write json get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandleGetShowOperationsLogById(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	log.Print("user by id")
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	models, err := m.cardsSvc.ShowOperationsLogById(id)
	if err != nil {
		log.Print("can't get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Print("can't write json get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandlePostAddHistory(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	log.Print("post user")
	model := history.ModelOperationsLog{}

	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't READ json POST model: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(model)
	m.cardsSvc.AddNewHistory(model)

}

func (m *MainServer) HandleGetShowOperationsLogByOwnerId(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	log.Print("user by id")
	id := authentication.Id
	models, err := m.cardsSvc.ShowOperationsLogByOwnerId(id)
	if err != nil {
		log.Print("can't get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Print("can't write json get all history")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
