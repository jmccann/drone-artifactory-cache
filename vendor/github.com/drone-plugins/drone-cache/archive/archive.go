package archive

import (
	"fmt"
	"io"
	"strings"
)

// Archive is an interface for packing and unpacking archive formats.
type Archive interface {
	// Pack writes an archive containing the source
	Pack(srcs []string, w io.Writer) error

	// Unpack reads the archive and restores it to the destination
	Unpack(dst string, r io.Reader) error
}

// FromFilename determines the archive format to use based on the name.
func FromFilename(name string) (Archive, error) {
	if strings.HasSuffix(name, ".tar") {
		return NewTarArchive(), nil
	}

	return nil, fmt.Errorf("Unknown file format for archive %s", name)
}
