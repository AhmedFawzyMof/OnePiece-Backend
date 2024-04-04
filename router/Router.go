package router

import (
	"Onepiece/middleware"
	"net/http"
	"strings"
)

func (router *Trie) Router(res http.ResponseWriter, req *http.Request) {
	resultNode, isFound := router.Search(req.URL.Path)
	var token string = req.Header.Get("Authorization")

	if !isFound {
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}

	middleware.SetCors(res, resultNode.method)

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusOK)
		return
	}

	if token != "" && strings.Contains(token, "Token ") {
		token = strings.Replace(token, "Token ", "", 1)
		email, err := middleware.VerifyToken(token)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		resultNode.params["authEmail"] = email
	}

	if req.Method == resultNode.method {
		resultNode.handler(res, req, resultNode.params)
		return
	}
}
