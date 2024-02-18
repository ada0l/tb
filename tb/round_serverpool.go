package tb

import (
	"log"
	"net/url"
	"sync/atomic"
)

type RoundServerPool struct {
	backends []*Backend
	current  uint64
}

func (s *RoundServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (s *RoundServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (s *RoundServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

func (s *RoundServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

func (s *RoundServerPool) ChechHealth() {
	for _, b := range s.backends {
		status := "up"
		alive, _ := b.CheckHealth()
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}
