package repo

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

func cloneRemote(gitRemote string, auth transport.AuthMethod, dirPath string) (*git.Repository, error) {
	repo, err := git.PlainClone(dirPath, false, &git.CloneOptions{
		URL:  gitRemote,
		Auth: auth,
	})
	if err != nil {
		if err == transport.ErrEmptyRemoteRepository {
			logrus.Info("Remote repository empty, initializing...")
			repo, err := initEmptyLocal(dirPath)
			if err != nil {
				return nil, err
			}
			_, err = repo.CreateRemote(&config.RemoteConfig{
				Name: "origin",
				URLs: []string{gitRemote},
			})
			if err != nil {
				return nil, fmt.Errorf("could not add remote %s error %s", gitRemote, err.Error())
			}
			return repo, nil
		} else {
			return nil, fmt.Errorf("could not clone to %s from %s error %s", dirPath, gitRemote, err.Error())
		}
	}
	return repo, nil
}

func initEmptyLocal(dirPath string) (*git.Repository, error) {
	repo, err := git.PlainInit(dirPath, false)
	if err != nil {
		return nil, fmt.Errorf("could not init repo inside %s", dirPath)
	}
	return repo, nil
}
