package tb

import "net/url"

type ServerPool interface {
	NextIndex() int
	GetNextPeer() *Backend
	AddBackend(backend *Backend)
	MarkBackendStatus(backendUrl *url.URL, alive bool)
	ChechHealth()
}
