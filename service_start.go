package main

import (
	"bsb-pr-solr-snapshot-service/app"
	"bsb-pr-solr-snapshot-service/backup/blob"
	"bsb-pr-solr-snapshot-service/backup/fusion"
	"bsb-pr-solr-snapshot-service/backup/repo"
	"bsb-pr-solr-snapshot-service/endpoints"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		server()
	} else {
		cli()
	}
}

func cli(){
	fmt.Printf("%s CLI Mode, commit %s\n", app.ProductName, app.GitCommit)
	fmt.Printf("Usage ./%s action [options], available actions = %s\n", os.Args[0], app.AvailableActions)
	switch os.Args[1]{
	case "backupFusion":
		if len(os.Args) == 6 {
			fusionURL := os.Args[2]
			gitRemote := os.Args[3]
			gitUsername := os.Args[4]
			gitPassword := os.Args[5]
			export, err := fusion.GetZKExport(fusionURL)
			if err != nil{
				log.Error("Failed to get ZK export from Fusion ", err.Error())
				return
			}
			err = repo.Update(gitRemote, gitUsername,gitPassword, export.Response)
			if err != nil{
				log.Error("Failed to process and update repo error ", err.Error())
				return
			}
		} else{
			fmt.Printf("[backupFusion] - Backs up fusion to git repo, syntax [Fusion URL] [git-remote] [git username] [git password] \n")
		}
		break
	case "createBlob":
		if len(os.Args) == 6 {
			repoPath := os.Args[2]
			generateUser := os.Args[4]
			description := os.Args[5]
			image, err := blob.Generate(repoPath, blob.FusionRestExport, generateUser, description)
			if err != nil {
				log.Error("Blob generation error ", err.Error())
				return
			}
			err = image.Save(os.Args[3], 0777)
			if err != nil {
				log.Error("Failed to save generated blob error ", err.Error())
			}
		} else{
			fmt.Printf("[createBlob] - Creates IMG blob of repo, syntax [repo-path] [blob out path[ [user] [description]\n")
		}
	case "loadBlob":
		if len(os.Args) == 4 {
			blobPath := os.Args[2]
			image, err := blob.Read(blobPath)
			if err != nil {
				log.Error("Failed to parse blob error ", err.Error())
			}
			repoPath := os.Args[3]
			err = repo.Create(repoPath, image.ZKRoot)
			if err != nil {
				log.Error("Repo store failed error ", err.Error())
			}
		} else {
			fmt.Printf("[loadBlob] - Loads IMG blob to repo, syntax [blob path] [repo path]\n")
		}
	}
}

func server(){
	log.Info(app.ProductName)
	router := gin.Default()

	//Register endpoints with
	router.GET("/backup/start", endpoints.BackupFusionZookeeper)

	err := router.Run()
	if err != nil {
		log.Fatal("Failed to start scp!")
	}
}
