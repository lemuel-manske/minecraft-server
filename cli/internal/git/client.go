package git

type Client interface {
	CommitBackup(tag string) error
	ListTags() ([]string, error)
	Prune(keep int) error
}
