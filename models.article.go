package main

import "errors"

type article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articleList = []article{
	article{ID: 1, Title: "Article 1", Content: "Article 1 body"},
	article{ID: 2, Title: "Article 2", Content: "Article 2 body"},
	article{ID: 3, Title: "Article 3", Content: "Article 3 body"},
	article{ID: 4, Title: "Article 4", Content: "Article 4 body"},
	article{ID: 5, Title: "Article 5", Content: "Article 5 body"},
	article{ID: 6, Title: "Article 6", Content: "Article 6 body"},
	article{ID: 7, Title: "Article 7", Content: "Article 7 body"},
	article{ID: 8, Title: "Article 8", Content: "Article 8 body"},
	article{ID: 9, Title: "Article 9", Content: "Article 9 body"},
	article{ID: 10, Title: "Article 10", Content: "Article 10 body"},
	article{ID: 11, Title: "Article 11", Content: "Article 11 body"},
	article{ID: 12, Title: "Article 12", Content: "Article 12 body"},
	article{ID: 13, Title: "Article 13", Content: "Article 13 body"},
	article{ID: 14, Title: "Article 14", Content: "Article 14 body"},
}

// Return a list of all the articles
func getAllArticles() []article {
	return articleList
}

func getArticleByID(id int) (*article, error) {
	for _, a := range articleList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Article not found")
}
