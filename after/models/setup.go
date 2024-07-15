package models

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ESClient *elasticsearch.Client

const SearchIndex = "blogs"

func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("gin_elasticsearch:tmp_pwd@tcp(127.0.0.1:3306)/gin_elasticsearch?charset=utf8&parseTime=true"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func DBMigrate() {
	DB.AutoMigrate(&Blog{})
}

func ESClientConnection() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	ESClient = client
}

func ESCreateIndexIfNotExist() {
	_, err := esapi.IndicesExistsRequest{
		Index: []string{SearchIndex},
	}.Do(context.Background(), ESClient)

	if err != nil {
		ESClient.Indices.Create(SearchIndex)
	}
}
