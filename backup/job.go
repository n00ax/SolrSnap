package backup

import (
	log "github.com/sirupsen/logrus"
	"solr-snapshot-service/backup/fusion"
	"solr-snapshot-service/backup/repo"
	"solr-snapshot-service/config"
	"time"
)

func createJobResponse(err error, start time.Time) map[string]string {
	if err != nil {
		return map[string]string{
			"status":  "failure",
			"error":   err.Error(),
			"message": "Something failed, probably my fault :( -Noah",
		}
	} else {
		return map[string]string{
			"status":   "success",
			"duration": time.Now().Sub(start).String(),
			"message":  "Job complete, have a nice day :) -Noah",
		}
	}
}

func StartJob() (map[string]string, error) {
	start := time.Now()
	log.Info("Starting export from Fusion Zookeeper at ", time.Now())
	export, err := fusion.GetZKExport(config.FusionBaseUrl)
	if err != nil {
		return createJobResponse(err, start), err
	}
	err = repo.Update(config.GitRemote, config.GitUsername, config.GitPassword, export.Response)
	if err != nil {
		return createJobResponse(err, start), err
	}
	log.Info(export.Response.Path)
	log.Info("Finished export from Fusion Zookeeper, took ", time.Now().Sub(start))
	return createJobResponse(err, start), err
}
