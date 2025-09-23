package indexer

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/shivamhw/tele-index/internal/models"
	"github.com/shivamhw/tele-index/pkg/log"
)

type Indexer struct {
	idx       bleve.Index
	indexPath string
	opts      *IndexerOpts
	lock      *sync.RWMutex
}

type IndexerOpts struct {
	IndexSchema string
	IndexPath   string
}


func NewIndexer(opts *IndexerOpts) (*Indexer, error) {
	indx, err := OpenIdx(opts.IndexPath)
	if err == nil {
		return &Indexer{
			idx:       indx,
			indexPath: opts.IndexPath,
			opts:      opts,
			lock:      &sync.RWMutex{},
		}, nil
	}
	log.InfoF("opening index failed with", "err", err)
	mapping, err := createMapping(opts.IndexSchema)
	if err != nil {
		return nil, err
	}
	indx, err = CreateIdx(opts.IndexPath, mapping)
	if err != nil {
		return nil, fmt.Errorf("create index failed, err: %w", err)
	}
	return &Indexer{
		idx:       indx,
		indexPath: opts.IndexPath,
		opts:      opts,
	}, nil
}

// create mapping
func createMapping(file string) (*mapping.IndexMappingImpl, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("create mapping failed, err: %w", err)
	}
	indexMapping := bleve.NewIndexMapping()
	if err := json.Unmarshal(data, indexMapping); err != nil {
		return nil, fmt.Errorf("index unmarshel failed with err: %w", err)
	}
	return indexMapping, nil

}

// create database
func CreateIdx(indexPath string, mapping *mapping.IndexMappingImpl) (bleve.Index, error) {
	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		return nil, err
	}
	return index, nil
}

// load database
func OpenIdx(indexPath string) (bleve.Index, error) {
	index, err := bleve.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("err opening idx: err %w", err)
	}
	return index, nil
}

// migrate json

// save item
func (i *Indexer) Index(it models.Item) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	if it.FileName == "" {
		return fmt.Errorf("failed indexing %v, err: %w", it, err)
	}
	i.idx.Index(it.FileName, it)
	return nil
}

func (i *Indexer) GetTotalDocs() (uint64, error) {
	return i.idx.DocCount()
}

func (i *Indexer) BulkImport(items []*models.Item) error {
	batch := i.idx.NewBatch()
	for c, v := range items {
		if v.FileName == "" {
			fmt.Printf("skipping indexing of %v\n", v)
			continue
		}
		batch.Index(v.FileName, v)
		if c%100 == 0 {
			if err := i.idx.Batch(batch); err != nil {
				return fmt.Errorf("failed indexing batch, err: %w", err)
			}
			batch = i.idx.NewBatch()
		}
		if batch.Size() > 0 {
			if err := i.idx.Batch(batch); err != nil {
				return fmt.Errorf("failed indexing batch, err: %w", err)
			}
		}
	}
	return nil
}

// search item
func (i *Indexer) Search(query string, size int, from int, mode string) (results []*models.Item, err error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	q := bleve.NewMatchQuery(query)
	q.SetField("tokens")
	req := bleve.NewSearchRequestOptions(q, size, from, true)
	req.Explain = true
	// request specific fields
	req.Fields = []string{"*"}
	res, err := i.idx.Search(req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	for _, hit := range res.Hits {
		// fmt.Printf("%d. id=%s score=%.3f\n", i+1, hit.ID, hit.Score)
		r := &models.Item{
			Id:       int64(hit.Fields["id"].(float64)),
			FileName: hit.Fields["file_name"].(string),
			ChatId:   int64(hit.Fields["chat_id"].(float64)),
			From:     int64(hit.Fields["from"].(float64)),
			Size:     int64(hit.Fields["size"].(float64)),
			Tokens:   hit.Fields["tokens"].(string),
		}
		results = append(results, r)
	}
	return
}

func (i *Indexer) Close() error {
	return i.idx.Close()
}

func (i *Indexer) ListAll() {
	q := bleve.NewMatchAllQuery()
	req := bleve.NewSearchRequest(q)
	req.Size = 1000            // adjust size as needed
	req.Fields = []string{"*"} // fetch all stored fields

	res, err := i.idx.Search(req)
	if err != nil {
		fmt.Printf("search failed: %v", err)
	}

	fmt.Printf("Total documents: %d\n", res.Total)
	for i, hit := range res.Hits {
		fmt.Printf("%d. ID=%s\n", i+1, hit.ID)
		for field, value := range hit.Fields {
			fmt.Printf("   %s: %v\n", field, value)
		}
	}
}
