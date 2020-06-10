package pagination

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	var t1 interface{} = map[string]interface{}{"t": "1"}
	var t2 interface{} = map[string]interface{}{"t": "2"}
	var t3 interface{} = map[string]interface{}{"t": "3"}
	var t4 interface{} = map[string]interface{}{"t": "4"}
	var t5 interface{} = map[string]interface{}{"t": "5"}
	type args struct {
		limit       int
		currentPage int
		res         io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    *Result
		wantErr bool
	}{
		{
			name: "first page",
			args: args{
				limit:       3,
				currentPage: 1,
				res:         ioutil.NopCloser(strings.NewReader(sample)),
			}, want: &Result{
				PaginationData: pagination{
					CurrentPage: 1,
					LastPage:    2,
					ItemPerPage: 3,
					TotalItem:   5,
					TotalPages:  2,
				},
				Items: []interface{}{
					t1,
					t2,
					t3,
				},
			}, wantErr: false},
		{
			name: "second page",
			args: args{
				limit:       3,
				currentPage: 2,
				res:         ioutil.NopCloser(strings.NewReader(sample)),
			}, want: &Result{
				PaginationData: pagination{
					CurrentPage: 2,
					LastPage:    2,
					ItemPerPage: 3,
					TotalItem:   5,
					TotalPages:  2,
				},
				Items: []interface{}{
					t4,
					t5,
				},
			}, wantErr: false},
		{
			name: "third page",
			args: args{
				limit:       3,
				currentPage: 3,
				res:         ioutil.NopCloser(strings.NewReader(sample)),
			}, want: &Result{
				PaginationData: pagination{
					CurrentPage: 3,
					LastPage:    2,
					ItemPerPage: 3,
					TotalItem:   5,
					TotalPages:  2,
				},
				Items: nil,
			}, wantErr: false},
		{
			name: "no item",
			args: args{
				limit:       3,
				currentPage: 1,
				res:         ioutil.NopCloser(strings.NewReader(sample2)),
			}, want: nil, wantErr: true},
		{
			name: "less then request",
			args: args{
				limit:       10,
				currentPage: 1,
				res:         ioutil.NopCloser(strings.NewReader(sample)),
			}, want: &Result{
				PaginationData: pagination{
					CurrentPage: 1,
					LastPage:    1,
					ItemPerPage: 10,
					TotalItem:   5,
					TotalPages:  1,
				},
				Items: []interface{}{
					t1,
					t2,
					t3,
					t4,
					t5,
				},
			}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := Resolve(tt.args.limit, tt.args.currentPage, tt.args.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() = %#v\n, want %#v", got, tt.want)
			}
		})
	}
}
