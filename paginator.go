package pagination

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Result struct {
	PaginationData pagination    `json:"pagination_data"`
	Items          []interface{} `json:"items"`
}
type pagination struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	ItemPerPage int `json:"item_per_page"`
	TotalItem   int `json:"total_item"`
	TotalPages  int `json:"total_pages"`
}

type hits struct {
	Hits   []*hits     `json:"hits,omitempty"`
	Source interface{} `json:"_source,omitempty"`
}
type response struct {
	Hits hits `json:"hits,omitempty"`
}

func HttpWriter(w http.ResponseWriter, limit, currentPage int, r esapi.Response) error {
	w.Header().Set("Content-Type", "application/json")
	res, err := Resolve(limit, currentPage, r.Body)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(res)
}

func Resolve(limit, currentPage int, res io.ReadCloser) (*Result, error) {
	var err error
	defer func() {
		_ = res.Close()
	}()
	s := &response{}
	if err = json.NewDecoder(res).Decode(s); err != nil {
		return nil, err
	}

	items := make([]interface{}, 0)
	for i := range s.Hits.Hits {
		items = append(items, s.Hits.Hits[i].Source)
	}
	tp := len(items)
	if tp < 1 {
		return nil, fmt.Errorf("no result")
	}
	r := &Result{
		PaginationData: pagination{
			ItemPerPage: limit,
			CurrentPage: currentPage,
			TotalItem:   tp,
			LastPage:    int(math.Ceil(float64(tp) / float64(limit))),
			TotalPages:  int(math.Ceil(float64(tp) / float64(limit))),
		},
	}

	if currentPage > r.PaginationData.TotalPages || len(items) == 0 {
		return r, nil
	}
	if currentPage*limit > len(items) {
		r.Items = items[currentPage*(limit)-limit:]
		return r, nil
	}
	r.Items = items[(currentPage*limit)-limit : currentPage*limit]
	return r, nil
}
