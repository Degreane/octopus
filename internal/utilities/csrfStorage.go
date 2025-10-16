package utilities

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/degreane/octopus/internal/middleware"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

type MemoryStore struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

func (s *MemoryStore) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key], nil
}

func (s *MemoryStore) Set(key string, val []byte, exp time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
	return nil
}

func (s *MemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *MemoryStore) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string][]byte)
	return nil
}

func (s *MemoryStore) Close() error {
	return nil
}
func GetCsrfToken(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.CsrfStore.Get(c)
		if err != nil {
			log.Printf("Error getting CSRF session: %v", err)
			L.Push(lua.LNil)
			return 1
		}

		token := sess.Get("csrf_token")
		switch v := token.(type) {
		case string:
			L.Push(lua.LString(v))
		case nil:
			L.Push(lua.LNil)
		default:
			L.Push(lua.LString(fmt.Sprintf("%v", v)))
		}
		return 1
	}
}
