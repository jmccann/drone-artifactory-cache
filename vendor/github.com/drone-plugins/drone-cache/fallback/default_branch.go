package fallback

type defaultBranchStrategy struct {
	branch string
}

// NewDefaultBranchStrategy creates a FallbackStrategy where a default target
// branch is used when the cache is not found.
func NewDefaultBranchStrategy(branch string) Strategy {
	return &defaultBranchStrategy{
		branch: branch,
	}
}

func (d *defaultBranchStrategy) Branch() (string, error) {
	return d.branch, nil
}
