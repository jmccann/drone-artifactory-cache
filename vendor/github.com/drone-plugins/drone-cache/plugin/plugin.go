package plugin

import (
	"github.com/drone-plugins/drone-cache/cache"
	"github.com/drone-plugins/drone-cache/storage"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

const (
	restoreMode = "restore"
	rebuildMode = "rebuild"
)

func Exec(c *cli.Context, s storage.Storage) error {
	p, err := newPlugin(c)

	if err != nil {
		return err
	}

	ca, err := cache.NewCache(s)

	if err != nil {
		return err
	}

	path := p.Path + p.Filename

	if p.Mode == rebuildMode {
		log.Infof("Rebuilding cache at %s", path)
		err = ca.Rebuild(p.Mount, path)

		if err == nil {
			log.Infof("Cache rebuilt")
		}

		return err
	}

	log.Infof("Restoring cache at %s", path)
	err = ca.Restore(path)

	if err == nil {
		log.Info("Cache restored")
	}

	return err
}

type plugin struct {
	Filename string
	Path     string
	Mode     string
	Mount    []string
}
