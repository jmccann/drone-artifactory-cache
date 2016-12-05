package main

import (
	"github.com/jmccann/drone-artifactory-cache/storage"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/drone-plugins/drone-cache/plugin"
	"github.com/urfave/cli"
)

var artifactoryFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "url",
		Usage:  "artifactory url",
		EnvVar: "PLUGIN_URL,ARTIFACTORY_CACHE_URL",
	},
	cli.StringFlag{
		Name:   "username",
		Usage:  "artifactory username",
		EnvVar: "PLUGIN_USERNAME,ARTIFACTORY_CACHE_USERNAME",
	},
	cli.StringFlag{
		Name:   "password",
		Usage:  "artifactory password",
		EnvVar: "PLUGIN_PASSWORD,ARTIFACTORY_CACHE_PASSWORD",
	},
}

func artifactoryOptions(c *cli.Context) (*storage.ArtifactoryOptions, error) {
	url := c.String("url")

	if len(url) == 0 {
		return nil, fmt.Errorf("No url provided")
	}

	// Get the access credentials
	username := c.String("username")
	password := c.String("password")

	if len(username) == 0 {
		return nil, fmt.Errorf("No username provided")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("No password provided")
	}

	return &storage.ArtifactoryOptions{
		Url: url,
		Username: username,
		Password: password,
		DryRun:   false,
	}, nil
}

func artifactoryPlugin(c *cli.Context) error {
	opts, err := artifactoryOptions(c)

	if err != nil {
		return err
	}

	log.Infof("Using %s as the cache", opts.Url)

	s, err := storage.NewArtifactoryStorage(opts)

	if err != nil {
		return err
	}

	return plugin.Exec(c, s)
}
