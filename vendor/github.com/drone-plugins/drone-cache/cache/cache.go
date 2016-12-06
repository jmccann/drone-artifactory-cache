package cache

import (
	"io"

	"github.com/drone-plugins/drone-cache/archive"
	"github.com/drone-plugins/drone-cache/storage"

	log "github.com/Sirupsen/logrus"
)

type Build struct {
	Owner  string
	Repo   string
	Branch string
}

type Cache struct {
	s storage.Storage
}

func NewCache(s storage.Storage) (Cache, error) {
	return Cache{
		s: s,
	}, nil
}

func (c Cache) Rebuild(srcs []string, dst string) error {
	a, err := archive.FromFilename(dst)

	if err != nil {
		return err
	}

	return rebuildCache(srcs, dst, c.s, a)
}

func (c Cache) Restore(src string) error {
	a, err := archive.FromFilename(src)

	if err != nil {
		return err
	}

	err = restoreCache(src, c.s, a)

	// Cache plugin should print an error but it should not return it
	// this is so the build continues even if the cache cant be restored
	if err != nil {
		log.Warnf("Cache could not be restored %s", err)
	}

	return nil
}

func restoreCache(src string, s storage.Storage, a archive.Archive) error {
	reader, writer := io.Pipe()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		err := s.Get(src, writer)

		if err != nil {
			cw <- err
			return
		}
	}()

	return a.Unpack("", reader)
}

func rebuildCache(srcs []string, dst string, s storage.Storage, a archive.Archive) error {
	log.Infof("Rebuilding cache at %s to %s", srcs, dst)

	reader, writer := io.Pipe()
	defer reader.Close()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		err := a.Pack(srcs, writer)

		if err != nil {
			cw <- err
			return
		}
	}()

	return s.Put(dst, reader)
}
