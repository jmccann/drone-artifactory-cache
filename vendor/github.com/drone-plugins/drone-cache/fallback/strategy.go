package fallback

// Strategy is an interface that is used to determine where the cache should
// look if the current branch is not cached.
type Strategy interface {
	// Branch retrieves the fallback branch.
	Branch() (string, error)
}
