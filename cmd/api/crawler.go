package main

import (
	"encoding/json"
	"net/http"

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

func (app *application) crawlJobHandler(w http.ResponseWriter, r *http.Request) {

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

	var response Response

	response.jobID = jobID.ID

	responseJson, err := json.MarshalIndent(response, "", " ")

	if err != nil {
		http.Error(w, "Unable to provide response", http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)

}
