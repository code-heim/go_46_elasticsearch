package controllers

import (
	"fmt"
	"gin_elasticsearch/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BlogsIndex(c *gin.Context) {
	var blogs *[]models.Blog
	query := c.Query("query")
	if query != "" {
		blogs = models.BlogSearch(query)
	} else {
		blogs = models.BlogsAll()
	}

	c.HTML(
		http.StatusOK,
		"blogs/index.tpl",
		gin.H{
			"blogs": blogs,
		},
	)
}

func BlogsShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	blog := models.BlogsFind(id)
	c.HTML(
		http.StatusOK,
		"blogs/show.tpl",
		gin.H{
			"blog": blog,
		},
	)
}

func BlogsBuildSearchIndex(c *gin.Context) {
	blogs := models.BlogsAll()
	for _, blog := range *blogs {
		blog.AddToIndex()
	}
}
