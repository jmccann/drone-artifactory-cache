package storage

import (
	"io"
	"testing"
	. "github.com/franela/goblin"

	"github.com/jmccann/drone-artifactory-cache/storage/fixtures"
	"net/http/httptest"
)

func TestArtifactoryStorage(t *testing.T) {
	g := Goblin(t)
	server := httptest.NewServer(fixtures.Handler())

	g.Describe("ArtifactoryStorage", func() {
		g.After(func() {
			server.Close()
		})

		g.It("Should create Storage object without errors", func() {
			_, err := NewArtifactoryStorage(opts)

			g.Assert(err == nil).IsTrue("should create storage object")
		})

		g.It("Should upload a file", func() {
			// Use fake server url
			opts.Url = server.URL

			// Create new storage object
			s, err := NewArtifactoryStorage(opts)

			g.Assert(err == nil).IsTrue("Failed to create storage object")

			// Act like 'cache'
			reader, writer := io.Pipe()

			cw := make(chan error)

			go func() {
				defer writer.Close()

				io.WriteString(writer, "hello")

				cw <- nil
			}()

			cr := make(chan error)

			go func() {
				defer reader.Close()

				// Upload content
				cr <- s.Put("the-repo-key/project/filename.tar", reader)
			}()

			werr := <-cw
			rerr := <-cr

			g.Assert(rerr == nil).IsTrue("Failed to read the content to upload")
			g.Assert(werr == nil).IsTrue("Failed to upload the content")
		})
	})
}

var (
	opts = &ArtifactoryOptions{
		Url: "http://company.com",
		Username: "johndoe",
		Password: "supersecret",
		DryRun: true,
	}
)
