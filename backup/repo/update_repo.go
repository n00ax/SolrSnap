package repo

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io/ioutil"
	"os"
	"path/filepath"
	"solr-snapshot-service/backup/blob"
)
import log "github.com/sirupsen/logrus"

const DirSaveMode = 0777
const FileSaveMode = 0777

func Update(gitRemote string, gitUsername string, gitPassword string, payload blob.Readable) error {
	log.Info("Save job started to Git remote = ", gitRemote, " total objects = ", 1)
	// Create temp file, and then store data there
	dirPath, err := ioutil.TempDir("/tmp", "temp-repo")
	if err != nil {
		return fmt.Errorf("unable to create temporary directory %s error %s", dirPath, err.Error())
	}
	log.Info("Cloning repo to temp directory" + dirPath)
	gitAuth := http.BasicAuth{
		Username: gitUsername,
		Password: gitPassword,
	}
	repo, err := cloneRemote(gitRemote, &gitAuth, dirPath)
	if err != nil {
		return err
	}
	log.Info("Processing Zookeeper objects into directory ", dirPath)
	err = processNode(dirPath, payload)
	if err != nil {
		return err
	}
	log.Info("Pushing new data to Git")
	err = commit(dirPath, repo)
	if err != nil {
		return err
	}
	err = push(&gitAuth, repo)
	if err != nil {
		return err
	}
	return nil
}
func Create(repoPath string, payload blob.Readable) error {
	log.Info("Create job started to local repo = %s", payload)
	err := os.Mkdir(repoPath, DirSaveMode)
	if err != nil {
		return fmt.Errorf("failed to create repo directory error %s", err.Error())
	}
	repo, err := initEmptyLocal(repoPath)
	log.Info("Processing BLOB Zookeeper objects into directory")
	err = processNode(repoPath, payload)
	if err != nil {
		return err
	}
	err = commit(repoPath, repo)
	if err != nil {
		return err
	}
	return nil
}
func processNode(dirPath string, node blob.Readable) error {
	newDirPath := fmt.Sprintf("%s%s", dirPath, node.GetPath())
	err := os.MkdirAll(newDirPath, DirSaveMode)
	if err != nil {
		return fmt.Errorf("failed to create directory %s with zookeeper path %s error %s", newDirPath, node.GetPath(), err.Error())
	}
	// create file named directory if data field exists
	if node.GetData() != "" {
		_, fileName := filepath.Split(node.GetPath())
		newDirDataPath := fmt.Sprintf("%s/%s", newDirPath, fileName)
		err = ioutil.WriteFile(newDirDataPath, []byte(node.GetData()), FileSaveMode)
		if err != nil {
			return fmt.Errorf("could not create node data file %s with zookeeper path %s error %s", newDirDataPath, node.GetPath(), err.Error())
		}
	}
	for _, child := range node.GetChildren() {
		err := processNode(dirPath, child)
		if err != nil {
			return err
		}
	}
	return nil
}
