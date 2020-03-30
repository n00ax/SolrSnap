package repo

import (
	"encoding/json"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"solr-snapshot-service/config"
	"time"
)

type CommitMessage struct {
	Service     string `json:"service"`
	CommitHash  string `json:"commitHash"`
	StartMethod string `json:"startMethod"`
}

func commit(dirPath string, repository *git.Repository) error {
	workTree, err := repository.Worktree()
	if err != nil {
		return fmt.Errorf("couldn't get the repository worktree, error %s", dirPath)
	}
	globPattern := fmt.Sprintf("%s/", dirPath)
	err = workTree.AddGlob("*")
	if err != nil {
		return fmt.Errorf("couldn't add files with pattern %s error %s", globPattern, err.Error())
	}
	commitMessage, err := generateCommitMessage()
	if err != nil {
		return err
	}
	_, err = workTree.Commit(commitMessage, &git.CommitOptions{
		All: false,
		Author: &object.Signature{
			Name:  "Zookeeper",
			Email: "none",
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  config.ProductName,
			Email: config.CommitterEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("couldn't commit objects error %s", err.Error())
	}
	return nil
}

func push(auth transport.AuthMethod, repository *git.Repository) error {
	err := repository.Push(&git.PushOptions{Auth: auth, RemoteName: "origin"})
	if err != nil {
		return fmt.Errorf("couldn't push to remote error %s", err.Error())
	}
	return nil
}

func generateCommitMessage() (string, error) {
	json, err := json.Marshal(CommitMessage{
		Service:    config.ProductName,
		CommitHash: config.GitCommit,
		//TODO replace
		StartMethod: "service",
	})
	if err != nil {
		return "", fmt.Errorf("could not marshall commit message json error %s", err.Error())
	}
	return string(json), nil
}
