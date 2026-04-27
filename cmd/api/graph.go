package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *application) graphImageHandler(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "jobID"))

	if err != nil {
		http.Error(w, "unable to parse uuid", http.StatusBadRequest)
		return
	}
	bytes, err := app.graphSvc.CreateGraph(r.Context(), jobID, "png")
	if err != nil {
		http.Error(w, "unable to create graph", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
