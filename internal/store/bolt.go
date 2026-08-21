package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"engineering-document-vault/internal/domain"
	"go.etcd.io/bbolt"
)

var (
	documentsBucket = []byte("ProjectDocument")
	reviewsBucket   = []byte("ReviewRecord")
	revisionsBucket = []byte("DocumentRevision")
	archivesBucket  = []byte("ArchiveRecord")
)

type BoltStore struct{ db *bbolt.DB }

func Open(path string) (*BoltStore, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	store := &BoltStore{db: db}
	if err = store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *BoltStore) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{documentsBucket, reviewsBucket, revisionsBucket, archivesBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BoltStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func putJSON(tx *bbolt.Tx, bucket []byte, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put([]byte(key), data)
}

func getJSON(bucket *bbolt.Bucket, key string, out any) error {
	value := bucket.Get([]byte(key))
	if value == nil {
		return errors.New("record not found")
	}
	return json.Unmarshal(value, out)
}

func (s *BoltStore) SaveDocument(document domain.ProjectDocument) error {
	if err := document.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, documentsBucket, document.ID, document) })
}

func (s *BoltStore) LoadDocument(id string) (domain.ProjectDocument, error) {
	var document domain.ProjectDocument
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx.Bucket(documentsBucket), id, &document) })
	return document, err
}

func (s *BoltStore) ListDocuments() ([]domain.ProjectDocument, error) {
	documents := []domain.ProjectDocument{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(documentsBucket).ForEach(func(_, value []byte) error {
			var document domain.ProjectDocument
			if err := json.Unmarshal(value, &document); err != nil {
				return err
			}
			documents = append(documents, document)
			return nil
		})
	})
	domain.SortDocuments(documents)
	return documents, err
}

func (s *BoltStore) SaveReview(review domain.ReviewRecord) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, reviewsBucket, review.ID, review) })
}

func (s *BoltStore) ListReviews(documentID string) ([]domain.ReviewRecord, error) {
	items := []domain.ReviewRecord{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(reviewsBucket).ForEach(func(_, value []byte) error {
			var item domain.ReviewRecord
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			if item.DocumentID == documentID {
				items = append(items, item)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].Sequence < items[j].Sequence })
	return items, err
}

func (s *BoltStore) SaveRevision(revision domain.DocumentRevision) error {
	key := fmt.Sprintf("%s:%d", revision.DocumentID, revision.Version)
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, revisionsBucket, key, revision) })
}

func (s *BoltStore) ListRevisions(documentID string) ([]domain.DocumentRevision, error) {
	items := []domain.DocumentRevision{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(revisionsBucket).ForEach(func(_, value []byte) error {
			var item domain.DocumentRevision
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			if item.DocumentID == documentID {
				items = append(items, item)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].Version < items[j].Version })
	return items, err
}

func (s *BoltStore) SaveArchive(record domain.ArchiveRecord) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, archivesBucket, record.ID, record) })
}

func (s *BoltStore) LoadArchive(documentID string) (domain.ArchiveRecord, error) {
	var result domain.ArchiveRecord
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(archivesBucket).ForEach(func(_, value []byte) error {
			var item domain.ArchiveRecord
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			if item.DocumentID == documentID {
				result = item
				return nil
			}
			return nil
		})
	})
	if result.ID == "" && err == nil {
		err = errors.New("archive record not found")
	}
	return result, err
}

func TempPath(prefix string) string {
	file, err := os.CreateTemp("", prefix+"-*.db")
	if err != nil {
		return ""
	}
	path := file.Name()
	_ = file.Close()
	_ = os.Remove(path)
	return path
}
