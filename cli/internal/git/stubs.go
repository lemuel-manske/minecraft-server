package git

type Stub struct {
	CommitErr error
	ListErr   error
	PruneErr  error
	Tags      []string

	CommitTag string
	PruneKeep int
}

func (s *Stub) CommitBackup(tag string) error {
	s.CommitTag = tag
	return s.CommitErr
}

func (s *Stub) ListTags() ([]string, error) {
	return s.Tags, s.ListErr
}

func (s *Stub) Prune(keep int) error {
	s.PruneKeep = keep
	return s.PruneErr
}
