package pagination

import (
	"encoding/json"
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
	Total struct {
		Value int `json:"value,omitempty"`
	} `json:"total,omitempty"`
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
	tp := s.Hits.Total.Value
	r := &Result{
		PaginationData: pagination{
			ItemPerPage: limit,
			CurrentPage: currentPage,
			TotalItem:   tp,
			LastPage:    tp,
			TotalPages:  int(math.Ceil(float64(tp) / float64(limit))),
		},
	}
	items := make([]interface{}, 0)
	if currentPage > r.PaginationData.TotalPages {
		return r, nil
	}
	for i := range s.Hits.Hits {
		items = append(items, s.Hits.Hits[i].Source)
	}
	r.Items = items
	return r, nil
}
