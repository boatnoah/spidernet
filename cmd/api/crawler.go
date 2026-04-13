package main

import (
	"encoding/json"
	"net/http"

	"github.com/boatnoah/spidernet/internal/queue"
	"github.com/boatnoah/spidernet/internal/store"
	"github.com/google/uuid"
)

type JobPayload struct {
	StartURL string
	Depth    int
}

type Response struct {
	JobID uuid.UUID `json:"job_id"`
}

func (app *application) submitJobHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody JobPayload

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(&requestBody)

	if err != nil {
		http.Error(w, "Bad request body", http.StatusBadRequest)
		return
	}

	jobID, err := app.store.CrawlJobs.CreateJob(
		r.Context(),
		store.CrawlJobPayload{
			StartUrl: requestBody.StartURL,
			MaxDepth: requestBody.Depth,
			Status:   "running",
		})

	if err != nil {
		http.Error(w, "Unable to provide response", http.StatusInternalServerError)
		return

	}

	var pageTask queue.PageTask

	pageTask.CrawlJobID = jobID.ID
	pageTask.Depth = requestBody.Depth

	err = app.queue.Add(r.Context(), &pageTask)
	if err != nil {
		http.Error(w, "unable to start job", http.StatusInternalServerError)
		return
	}

	var response Response

	response.JobID = jobID.ID

	responseJson, err := json.MarshalIndent(response, "", " ")

	if err != nil {
		http.Error(w, "Unable to provide response", http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)

}
