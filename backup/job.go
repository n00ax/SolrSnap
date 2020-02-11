package backup

import (
	"bsb-pr-solr-snapshot-service/app"
	"bsb-pr-solr-snapshot-service/backup/fusion"
	"bsb-pr-solr-snapshot-service/backup/repo"
	log "github.com/sirupsen/logrus"
	"time"
)

func createJobResponse(err error, start time.Time) map[string]string{
	if err != nil {
		return map[string]string{
			"status" : "failure",
			"error" : err.Error(),
			"message" : "Something failed, probably my fault :( -Noah",
		}
	} else{
		return map[string]string{
			"status" : "success",
			"duration" : time.Now().Sub(start).String(),
			"message" : "Job complete, have a nice day :) -Noah",
		}
	}
}

func StartJob() (map[string]string, error) {
	start := time.Now()
	log.Info("Starting export from Fusion Zookeeper at ", time.Now())
	export, err := fusion.GetZKExport(app.FusionBaseUrl)
	if err != nil{
		return createJobResponse(err, start), err
	}
	err = repo.Update(app.GitRemote, app.GitUsername, app.GitPassword, export.Response)
	if err != nil{
		return createJobResponse(err, start), err
	}
	log.Info(export.Response.Path)
	log.Info("Finished export from Fusion Zookeeper, took ", time.Now().Sub(start))
	return createJobResponse(err, start), err
}
