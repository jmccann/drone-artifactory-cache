package plugin

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	filenameFlag = "filename"
	pathFlag     = "path"
	mountFlag    = "mount"
	rebuildFlag  = "rebuild"
	restoreFlag  = "restore"

	debugFlag = "debug"

	repoOwnerFlag    = "repo.owner"
	repoNameFlag     = "repo.name"
	commitBranchFlag = "commit.branch"
)

func PluginFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   filenameFlag,
			Usage:  "Filename for the cache",
			EnvVar: "PLUGIN_FILENAME",
		},
		cli.StringFlag{
			Name:   pathFlag,
			Usage:  "path",
			EnvVar: "PLUGIN_PATH",
		},
		cli.StringFlag{
			Name:   "mount",
			Usage:  "cache directories",
			EnvVar: "PLUGIN_MOUNT",
		},
		cli.BoolFlag{
			Name:   rebuildFlag,
			Usage:  "rebuild the cache directories",
			EnvVar: "PLUGIN_REBUILD",
		},
		cli.BoolFlag{
			Name:   restoreFlag,
			Usage:  "restore the cache directories",
			EnvVar: "PLUGIN_RESTORE",
		},

		cli.BoolFlag{
			Name:   debugFlag,
			Usage:  "debug plugin output",
			EnvVar: "PLUGIN_DEBUG",
		},

		// Build information

		cli.StringFlag{
			Name:   repoOwnerFlag,
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   repoNameFlag,
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   commitBranchFlag,
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
	}
}

func newPlugin(c *cli.Context) (*plugin, error) {
	if c.GlobalBool(debugFlag) {
		log.SetLevel(log.DebugLevel)
	}

	// Determine the mode for the plugin
	rebuild := c.GlobalBool(rebuildFlag)
	restore := c.GlobalBool(restoreFlag)

	if rebuild && restore {
		return nil, errors.New("Cannot rebuild and restore the cache")
	} else if !rebuild && !restore {
		return nil, errors.New("No action specified")
	}

	var mode string
	var mount string

	if rebuild {
		// Look for the mount points to rebuild
		mount = c.GlobalString(mountFlag)

		if len(mount) == 0 {
			return nil, errors.New("No mounts specified")
		}

		mode = rebuildMode
	} else {
		mode = restoreMode
		mount = ""
	}

	// Get the path to place the cache files
	path := c.GlobalString(pathFlag)

	// Defaults to <owner>/<repo>/<branch>/
	if len(path) == 0 {
		log.Info("No path specified. Creating default")

		path = fmt.Sprintf(
			"/%s/%s/%s/",
			c.GlobalString(repoOwnerFlag),
			c.GlobalString(repoNameFlag),
			c.GlobalString(commitBranchFlag),
		)
	}

	// Get the filename
	filename := c.GlobalString(filenameFlag)

	if len(filename) == 0 {
		log.Info("No filename specified. Creating default")

		filename = "archive.tar"
	}

	return &plugin{
		Filename: filename,
		Path:     path,
		Mount:    mount,
		Mode:     mode,
	}, nil
}
