package repository

type Loader interface {
	LoadRepository(string, RepositoryType) (*Repository, error)
}
