package storage

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	. "github.com/drone-plugins/drone-cache/storage"

	log "github.com/Sirupsen/logrus"
	"github.com/jfrogdev/jfrog-cli-go/artifactory/commands"
	"github.com/jfrogdev/jfrog-cli-go/artifactory/utils"
	"github.com/jfrogdev/jfrog-cli-go/utils/config"
	"github.com/jfrogdev/jfrog-cli-go/utils/ioutils"
)

// ArtifactoryOptions contains configuration for the artifactory connection.
type ArtifactoryOptions struct {
	Url      string
	Username string
	Password string
	DryRun   bool
}

type artifactoryStorage struct {
	opts   *ArtifactoryOptions
}

// NewS3Storage creates an implementation of Storage with S3 as the backend.
func NewArtifactoryStorage(opts *ArtifactoryOptions) (Storage, error) {
	return &artifactoryStorage{
		opts:   opts,
	}, nil
}

func (s *artifactoryStorage) Get(p string, dst io.Writer) error {
	// Get repo key and download path
	repo, path, err := parseArtifactoryUrl(p)

	if err != nil {
		return err
	}

	// Check repo exists
	err = isRepoExist(s, repo)
	if err != nil {
		return err
	}

	// Create tmpfile to write to
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		return err
	}

	// Download file
	flags := createDownloadFlags(s)
	spec := createDownloadSpec(repo, path, tmpfile)
	err = commands.Download(spec, flags)

	logFileSize(tmpfile)

	if err != nil {
		return err
	}

	_, err = io.Copy(dst, tmpfile)

	if err != nil {
		return err
	}

	return nil
}

func (s *artifactoryStorage) Put(p string, src io.Reader) error {
	// Get repo key and upload path
	repo, path, err := parseArtifactoryUrl(p)

	if err != nil {
		return err
	}

	// Check repo exists
	err = isRepoExist(s, repo)
	if err != nil {
		return err
	}

	// Create temp file from src
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name()) // clean up

	contents, err := ioutil.ReadAll(src)

	if err != nil {
		return err
	}

	if _, err := tmpfile.Write(contents); err != nil {
		return err
	}

	logFileSize(tmpfile)

	// Close file
	if err := tmpfile.Close(); err != nil {
		return err
	}

	// Upload file
	flags := createUploadFlags(s)
	spec := utils.CreateSpec(tmpfile.Name(), fmt.Sprintf("/%s/%s", repo, path), "", false, true, false)
	uploaded, failed, err := commands.Upload(spec, flags)

	if err != nil {
		return err
	}

	if failed > 0 {
		return errors.New("Upload failed")
	}

	if uploaded == 0 {
		return errors.New("No files uploaded")
	}

	return nil
}

func createDownloadFlags(s *artifactoryStorage) *commands.DownloadFlags {
	flags := new(commands.DownloadFlags)
	url := s.opts.Url

	if strings.HasSuffix(s.opts.Url, "/") == false {
		url = s.opts.Url + "/"
	}

	flags.ArtDetails = new(config.ArtifactoryDetails)
	flags.ArtDetails.Url = url
	flags.ArtDetails.User = s.opts.Username
	flags.ArtDetails.Password = s.opts.Password
	flags.DryRun = s.opts.DryRun
	flags.Threads = 3

	return flags
}

func createDownloadSpec(repo string, path string, f *os.File) *utils.SpecFiles {
	pattern := fmt.Sprintf("%s/%s", repo, path)
	target := f.Name()

	return utils.CreateSpec(pattern, target, "", false, true, false)
}

func createUploadFlags(s *artifactoryStorage) *commands.UploadFlags {
	flags := new(commands.UploadFlags)
	flags.ArtDetails = new(config.ArtifactoryDetails)
	flags.ArtDetails.Url = s.opts.Url
	flags.ArtDetails.User = s.opts.Username
	flags.ArtDetails.Password = s.opts.Password
	flags.DryRun = s.opts.DryRun
	flags.Threads = 3

	return flags
}

func isRepoExist(s *artifactoryStorage, repokey string) error {
	log.Infof("Checking repo %s exists", fmt.Sprintf("%s/%s", s.opts.Url, repokey))
	resp, _, _,err := ioutils.SendGet(fmt.Sprintf("%s/%s", s.opts.Url, repokey), true, ioutils.HttpClientDetails{User: s.opts.Username, Password: s.opts.Password})

	if err != nil {
		return err
	}

	if resp.StatusCode != 400 {
		return nil
	}

	return errors.New(fmt.Sprintf("Repo %s does not exist", fmt.Sprintf("%s/%s", s.opts.Url, repokey)))
}

func logFileSize(f *os.File) error {
	// Find file size
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	log.Infof("The archive file is %d bytes", fi.Size())

	return nil
}

func parseArtifactoryUrl(p string) (string, string, error) {
	u, err := url.Parse(p)

	if err != nil {
		return "", "", err
	}

	path := u.Path

	// Remove initial forward slash
	e := strings.Split(path, "/")

	return e[0], strings.Join(e[1:len(e)], "/"), nil
}
