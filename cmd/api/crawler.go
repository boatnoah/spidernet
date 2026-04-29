package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/boatnoah/spidernet/internal/queue"
	"github.com/boatnoah/spidernet/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type JobPayload struct {
	StartURL string `json:"start_url"`
	Depth    int    `json:"depth"`
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
		log.Print(err)
		http.Error(w, "Unable to provide response", http.StatusInternalServerError)
		return
	}

	pageTask := queue.CreatePageTask(jobID.ID, requestBody.StartURL, 0)

	err = app.queue.Add(r.Context(), pageTask)
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

func (app *application) getJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "jobID"))

	if err != nil {
		http.Error(w, "unable to parse job id", http.StatusBadRequest)
		return
	}

	crawlJob, err := app.store.CrawlJobs.GetJobById(r.Context(), jobID)

	if err != nil {
		http.Error(w, "unable to find job", http.StatusBadRequest)
		return
	}

	jsonB, err := json.MarshalIndent(crawlJob, "", " ")

	if err != nil {
		http.Error(w, "unable to parse json", http.StatusInternalServerError)
		return
	}

	w.Write(jsonB)

}
