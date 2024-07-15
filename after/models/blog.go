package models

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title   string `gorm:"size:255"`
	Content string `gorm:"type:text"`
}

func BlogsAll() *[]Blog {
	var blogs []Blog
	DB.Where("deleted_at is NULL").Order("updated_at desc").Find(&blogs)
	return &blogs
}

func BlogsFind(id uint64) *Blog {
	var blog Blog
	DB.Where("id = ?", id).First(&blog)
	return &blog
}

func (b *Blog) AddToIndex() {
	document := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		b.Title, b.Content,
	}
	data, err := json.Marshal(document)
	if err != nil {
		log.Fatalf("Error marshaling the document: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      SearchIndex,
		DocumentID: strconv.Itoa(int(b.ID)),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	resp, err := req.Do(context.Background(), ESClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer resp.Body.Close()

	log.Printf("Indexed document %s to index %s\n", resp.String(), SearchIndex)
}

func BlogSearch(searchQuery string) *[]Blog {
	var buf bytes.Buffer
	/*query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": searchQuery,
			},
		},
	}*/

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchQuery,
				"fields": []string{"title", "content"},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := ESClient.Search(
		ESClient.Search.WithIndex(SearchIndex),
		ESClient.Search.WithBody(&buf),
	)

	if err != nil {
		return nil
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil
	}

	// Extract IDs from the response
	var ids []uint
	if hits, ok := r["hits"].(map[string]interface{}); ok {
		if hitsHits, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsHits {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if idStr, ok := hitMap["_id"].(string); ok {
						id, _ := strconv.Atoi(idStr)
						ids = append(ids, uint(id))
					}
				}
			}
		}
	}

	var blogs []Blog
	DB.Where("deleted_at is NULL").Where("id in ?", ids).Order("updated_at desc").Find(&blogs)

	return &blogs
}
