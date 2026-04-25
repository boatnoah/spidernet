package main

import (
	"io"
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
	bytes, err := app.graphSvc.CreateGraph(r.Context(), jobID)
	if err != nil {
		http.Error(w, "unable to create graph", http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(bytes)

	if err != nil {
		http.Error(w, "unable to convert to bytes", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(b))

}
