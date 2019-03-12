package api

import (
	"log"
	"os"
)

func getProject() (projectID string) {
	var exists bool
	envCandidates := []string{
		"GCP_PROJECT",
		"GOOGLE_CLOUD_PROJECT",
		"GCLOUD_PROJECT",
	}
	for _, e := range envCandidates {
		projectID, exists = os.LookupEnv(e)
		if exists {
			return projectID
		}
	}
	if !exists {
		log.Fatalf("Set project ID via one of the supported env variables.")
	}
	return ""
}
