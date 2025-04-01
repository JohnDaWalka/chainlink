package datastore

import (
	"sync"
)

type ContractMetadataStore[M Cloneable[M]] interface {
	Store[ContractMetadataKey, ContractMetadata[M]]
}

type MutableContractMetadataStore[M Cloneable[M]] interface {
	MutableStore[ContractMetadataKey, ContractMetadata[M]]
}

var _ ContractMetadataStore[DefaultMetadata] = &MemoryContractMetadataStore[DefaultMetadata]{}
var _ MutableContractMetadataStore[DefaultMetadata] = &MemoryContractMetadataStore[DefaultMetadata]{}

type MemoryContractMetadataStore[M Cloneable[M]] struct {
	mu      sync.RWMutex
	records []ContractMetadata[M]
}

func NewMemoryContractMetadataStore[M Cloneable[M]]() *MemoryContractMetadataStore[M] {
	return &MemoryContractMetadataStore[M]{records: []ContractMetadata[M]{}}
}

func (s *MemoryContractMetadataStore[M]) indexOf(key ContractMetadataKey) int {
	for i, record := range s.records {
		if record.Key().Equals(key) {
			return i
		}
	}
	return -1
}

func (s *MemoryContractMetadataStore[M]) Get(key ContractMetadataKey) (ContractMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ContractMetadata[M]{}, ErrContractMetadataNotFound
	}
	return s.records[idx].Clone(), nil
}

func (s *MemoryContractMetadataStore[M]) Fetch() ([]ContractMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := []ContractMetadata[M]{}
	for _, record := range s.records {
		records = append(records, record.Clone())
	}
	return records, nil
}

func (s *MemoryContractMetadataStore[M]) Filter(filters ...FilterFunc[ContractMetadataKey, ContractMetadata[M]]) []ContractMetadata[M] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := append([]ContractMetadata[M]{}, s.records...)
	for _, filter := range filters {
		records = filter(records)
	}
	return records
}

func (s *MemoryContractMetadataStore[M]) Add(record ContractMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		return ErrContractMetadataExists
	}
	s.records = append(s.records, record)
	return nil
}

func (s *MemoryContractMetadataStore[M]) AddOrUpdate(record ContractMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		s.records = append(s.records, record)
		return nil
	}
	s.records[idx] = record
	return nil
}

func (s *MemoryContractMetadataStore[M]) Update(record ContractMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		return ErrContractMetadataNotFound
	}
	s.records[idx] = record
	return nil
}

func (s *MemoryContractMetadataStore[M]) Delete(key ContractMetadataKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ErrContractMetadataNotFound
	}
	s.records = append(s.records[:idx], s.records[idx+1:]...)
	return nil
}
