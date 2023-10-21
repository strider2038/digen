package httphandler

import (
	"encoding/json"
	"net/http"

	"basic/app/usecase"
)

type FindEntityRequest struct {
	ID int `json:"id"`
}

type FindEntity struct {
	useCase *usecase.FindEntity
}

func NewFindEntity(useCase *usecase.FindEntity) *FindEntity {
	return &FindEntity{useCase: useCase}
}

func (h *FindEntity) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var r FindEntityRequest
	if err := json.NewDecoder(request.Body).Decode(&r); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := h.useCase.Handle(request.Context(), r.ID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(entity)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}
