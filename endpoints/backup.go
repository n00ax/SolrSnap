package endpoints

import (
	"bsb-pr-solr-snapshot-service/backup"
	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)


func BackupFusionZookeeper(context* gin.Context){

	response, err := backup.StartJob()
	if err == nil {
		context.JSON(200, response)
	} else {
		log.Error(err.Error())
		context.JSON(500, response)
	}

}