package store

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

var log_key = "STORE"

type Store struct {
	dbEnv           *lmdb.Env
	clientToChannel map[string]chan<- []byte
	keyToClients    map[string][]string
	currentClientId int16
	m               sync.Mutex
}

func NewStore(dbPath string) *Store {
	env, err := lmdb.NewEnv()
	if err != nil {
		log.Fatalf("[%s] Error while creating lmdb environment: %v", log_key, err)
	}

	if err := env.SetMapSize(1 << 30); err != nil { // 1GB
		log.Fatalf("[%s] Error setting map size: %v", log_key, err)
	}
	if err := env.SetMaxDBs(1); err != nil {
		log.Fatalf("[%s] Error setting max dbs: %v", log_key, err)
	}

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		log.Fatalf("[%s] Failed to create db path: %v", log_key, err)
	}

	if err := env.Open(dbPath, 0, 0664); err != nil {
		log.Fatalf("[%s] Error opening LMDB env: %v", log_key, err)
	}

	return &Store{
		dbEnv:           env,
		clientToChannel: make(map[string]chan<- []byte),
		keyToClients:    make(map[string][]string),
		currentClientId: 1000,
	}
}

func (s *Store) GetUniqueClient() string {
	s.m.Lock()
	defer s.m.Unlock()

	s.currentClientId = s.currentClientId + 1

	return fmt.Sprintf("%d", s.currentClientId)
}

func (s *Store) RegisterUserChannel(clientId string, ch chan<- []byte) {
	s.m.Lock()
	defer s.m.Unlock()

	s.clientToChannel[clientId] = ch
}

func (s *Store) Subscribe(clientId string, key string) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.clientToChannel[clientId]; !ok {
		return
	}

	s.keyToClients[key] = append(s.keyToClients[key], clientId)
}

func (s *Store) Set(key []byte, value []byte) error {
	txn, err := s.dbEnv.BeginTxn(nil, 0)
	if err != nil {
		return fmt.Errorf("begin txn failed: %w", err)
	}
	defer txn.Abort()

	dbi, err := txn.OpenRoot(0)
	if err != nil {
		return fmt.Errorf("open db failed: %w", err)
	}

	if err := txn.Put(dbi, key, value, 0); err != nil {
		return fmt.Errorf("put failed: %w", err)
	}

	if err := txn.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	s.m.Lock()
	defer s.m.Unlock()
	clients, ok := s.keyToClients[string(key)]
	if !ok {
		return nil
	}

	for _, client := range clients {
		if ch, ok := s.clientToChannel[client]; ok {
			ch <- value
		}
	}

	return nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	txn, err := s.dbEnv.BeginTxn(nil, lmdb.Readonly)
	if err != nil {
		return []byte{}, fmt.Errorf("begin read txn failed: %w", err)
	}
	defer txn.Abort()

	dbi, err := txn.OpenRoot(0)
	if err != nil {
		return []byte{}, fmt.Errorf("open db failed: %w", err)
	}

	val, err := txn.Get(dbi, []byte(key))
	if err == lmdb.NotFound {
		return []byte{}, fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return []byte{}, fmt.Errorf("get failed: %w", err)
	}

	return val, nil
}

func (s *Store) Close() {
	if s.dbEnv != nil {
		s.dbEnv.Close()
	}
}
