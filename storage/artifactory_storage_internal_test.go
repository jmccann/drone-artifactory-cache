package storage

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestArtifactoryStorageInternal(t *testing.T) {
	g := Goblin(t)

	g.Describe("createUploadFlags", func() {
		g.It("Should create proper flags", func() {
			flags := createUploadFlags(artiStor)

			g.Assert(flags.ArtDetails.Url).Equal("http://company.com")
			g.Assert(flags.DryRun).Equal(false)
			g.Assert(flags.Threads).Equal(3)
		})
	})

  g.Describe("parseArtifactoryUrl", func() {
		g.It("Should parse the artifactory url without subdirs", func() {
			repo, path, err := parseArtifactoryUrl("my-repo-key/filename.tar")

			g.Assert(err == nil).IsTrue("should not error")
			g.Assert(repo).Equal("my-repo-key")
			g.Assert(path).Equal("filename.tar")
		})

		g.It("Should parse the artifactory url with subdirs", func() {
			repo, path, err := parseArtifactoryUrl("my-repo-key/project/path/filename.tar")

			g.Assert(err == nil).IsTrue("should not error")
			g.Assert(repo).Equal("my-repo-key")
			g.Assert(path).Equal("project/path/filename.tar")
		})
	})
}

var (
	intOpts = &ArtifactoryOptions{
		Url: "http://company.com",
		Username: "johndoe",
		Password: "supersecret",
		DryRun: false,
	}

	artiStor = &artifactoryStorage{
		opts: intOpts,
	}
)
