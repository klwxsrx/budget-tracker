package persistence

type Synchronization interface {
	CriticalSection(name string, f func() error) error
}
