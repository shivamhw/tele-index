package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

type payload struct {
	File   string
	ChatId int64
	ID     int64
	Size   int64
	Tokens string
}

var payloads []payload

func loadPayloads(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var p payload
		err := json.Unmarshal([]byte(line), &p)
		if err != nil {
			log.Fatal(err)
			return
		}
		p.Tokens = getTokens(p.File)
		payloads = append(payloads, p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func migrateJson(indexPath string, exportFile string) {

	loadPayloads(exportFile)
	fmt.Print("Loaded payloads: ", len(payloads))
	// Open or create index
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		f, err := os.Open("index.json")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		indexMapping := bleve.NewIndexMapping()
		if err := json.NewDecoder(f).Decode(indexMapping); err != nil {
			fmt.Print(err)
			log.Fatal("failed to decode mapping:", err)
		}

		// create index with mapping
		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer index.Close()
	// Use batch for bulk insert
	batch := index.NewBatch()
	for i, p := range payloads {
		docID := p.File
		if docID == "" {
			log.Printf("docID is empty for payload: %v", p)
			continue
		}
		err := batch.Index(docID, p)
		if err != nil {
			log.Fatal(err)
		}

		// Commit every 1000 docs (example, adjust for your size)
		if (i+1)%1000 == 0 {
			if err := index.Batch(batch); err != nil {
				log.Fatal(err)
			}
			batch = index.NewBatch()
		}
	}

	// Commit remaining docs
	if batch.Size() > 0 {
		if err := index.Batch(batch); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Bulk indexing completed!")
}

func getTokens(fileName string) string {
	// This function should extract tokens from the fileName or item.
	// For simplicity, we return an empty string here.
	// Implement the actual logic as needed.
	fileName = strings.ReplaceAll(fileName, "_", " ")
	fileName = strings.ReplaceAll(fileName, "-", " ")
	fileName = strings.ReplaceAll(fileName, ".", " ")
	fileName = strings.ReplaceAll(fileName, ",", " ")
	fileName = strings.ReplaceAll(fileName, "@", "")
	fileName = strings.ReplaceAll(fileName, "#", "")
	fileName = strings.ReplaceAll(fileName, "(", " ")
	fileName = strings.ReplaceAll(fileName, ")", " ")
	fileName = strings.ReplaceAll(fileName, "[", " ")
	fileName = strings.ReplaceAll(fileName, "]", " ")

	return fileName
}
