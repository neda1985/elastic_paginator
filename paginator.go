package pagination

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"math"
	"net/http"
)

func paginate(writer http.ResponseWriter, limit int, currentPage int, r esapi.Response) {
	var (
		m map[string]interface{}
	)
	if r.IsError() {
		var e errorResponse
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				r.Status(),
				e.Error["type"],
				e.Error["reason"],
			)
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	var items []interface{}
	resultExtras := make(map[string]interface{})
	result := make(map[string]interface{})
	var totalPages = int(math.Ceil(float64(int(m["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))) / float64(limit)))
	var totalItems = int(m["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, hit := range m["hits"].(map[string]interface{})["hits"].([]interface{}) {
		items = append(items, hit.(map[string]interface{})["_source"])

	}
	if currentPage > totalPages {
		items = nil
	} else {
		if currentPage == 1 {
			if len(items) < limit {
				items = items[0:len(items)]
			} else {
				items = items[0:limit]
			}
		} else {
			//check if last page
			if totalItems > currentPage*limit {
				items = items[(currentPage-1)*limit-1 : currentPage*limit-1]
			} else {
				items = items[(currentPage-1)*limit-1:]
			}

		}
	}
	result["items"] = items
	resultExtras["current_page"] = currentPage
	resultExtras["last_page"] = totalPages
	resultExtras["items_per_page"] = limit
	resultExtras["items_count"] = len(items)
	resultExtras["total_items"] = totalItems
	resultExtras["total_pages"] = totalPages
	result["pagination_data"] = resultExtras
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Write(js)
}

type errorResponse struct {
	Error map[string]interface{}
}
