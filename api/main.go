package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/zoekt"
	"github.com/google/zoekt/query"
)

// DefaultAPIPath is the rpc path used by zoekt-webserver
const DefaultAPIPath = "/api"

type Server struct {
	Searcher zoekt.Searcher
}

func NewMux(searcher zoekt.Searcher) *http.ServeMux {
	s := &Server{
		Searcher: searcher,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Search)
	return mux
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryStr, ok := r.URL.Query()["q"]
	if !ok || len(queryStr[0]) < 1 {
		http.Error(w, `Url Param 'q' is missing`, http.StatusBadRequest)
		return
	}
	q, _ := query.Parse(queryStr[0])
	// if err != nil {
	// 	return err
	// }
	sOpts := &zoekt.SearchOptions{
		MaxWallTime: 10 * time.Second,
	}

	sOpts.SetDefaults()
	sOpts.ShardMaxMatchCount = 100
	sOpts.ShardMaxImportantMatch = 100
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := s.Searcher.Search(ctx, q, sOpts)
	if err != nil {
		// return err
	}
	jsonData, _ := json.Marshal(resp)
	w.Write(jsonData)

}
