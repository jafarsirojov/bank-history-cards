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
	if authentication.Id == 0 {
		log.Print("admin")
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
		return
	}
	log.Print("all history cards user")
	models, err := m.cardsSvc.ShowOperationsLogByOwnerId(authentication.Id)
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
	log.Print("transfer history by id card")
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't history by id card")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("can't strconv atoi: %d", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if authentication.Id == 0 {
		log.Print("admin")
		models, err := m.cardsSvc.AdminShowTransferLogByIdCadr(id)
		if err != nil {
			log.Printf("can't history %d", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print(models)
		err = rest.WriteJSONBody(writer, models)
		if err != nil {
			log.Printf("can't write json get all history %d", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	models, err := m.cardsSvc.UserShowTransferLogByIdCard(id,authentication.Id)
	if err != nil {
		log.Printf("can't history is not owner %d", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Printf("can't write json get all history %d", err)
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
	log.Print("post add history transfer")
	model := history.ModelOperationsLog{}

	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't READ json POST model: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(model)
	if model.Id != 0 {
		log.Printf("id card not 0!")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	m.cardsSvc.AddNewHistory(model)
}
