package inmem

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestInMemCache_AddToCache(t *testing.T) {
	type fields struct {
		records map[string][]byte
		mu      sync.RWMutex
	}
	type args struct {
		ctx      context.Context
		url      string
		response []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add record",
			fields: fields{
				records: make(map[string][]byte),
				mu:      sync.RWMutex{},
			},
			args: args{
				ctx:      context.Background(),
				url:      "test.com.",
				response: []byte("some test data"),
			},
			wantErr: false,
		},
		{
			name: "add empty record",
			fields: fields{
				records: make(map[string][]byte),
				mu:      sync.RWMutex{},
			},
			args: args{
				ctx:      context.Background(),
				url:      "test.com.",
				response: []byte(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemCache{
				records: tt.fields.records,
				mu:      tt.fields.mu,
			}
			if err := i.AddToCache(tt.args.ctx, tt.args.url, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("AddToCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemCache_DeleteFromCache(t *testing.T) {
	type fields struct {
		records map[string][]byte
		mu      sync.RWMutex
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete record",
			fields: fields{
				records: map[string][]byte{"test.com.": []byte("some test data")},
				mu:      sync.RWMutex{},
			},
			args: args{
				ctx: context.Background(),
				url: "test.com.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemCache{
				records: tt.fields.records,
				mu:      tt.fields.mu,
			}
			if err := i.DeleteFromCache(tt.args.ctx, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFromCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemCache_GetFromCache(t *testing.T) {
	type fields struct {
		records map[string][]byte
		mu      sync.RWMutex
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse []byte
		wantErr      bool
	}{
		{
			name: "get existing record",
			fields: fields{
				records: map[string][]byte{"test.com.": []byte("some test data")},
				mu:      sync.RWMutex{},
			},
			args: args{
				ctx: context.Background(),
				url: "test.com.",
			},
			wantResponse: []byte("some test data"),
			wantErr:      false,
		},
		{
			name: "get not existing record",
			fields: fields{
				records: map[string][]byte{"test.com.": []byte("some test data")},
				mu:      sync.RWMutex{},
			},
			args: args{
				ctx: context.Background(),
				url: "exist.com.",
			},
			wantResponse: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemCache{
				records: tt.fields.records,
				mu:      tt.fields.mu,
			}
			gotResponse, err := i.GetFromCache(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("GetFromCache() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestInMemCache_ListAllRecords(t *testing.T) {
	type fields struct {
		records map[string][]byte
		mu      sync.RWMutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantList []string
		wantErr  bool
	}{
		{
			name: "list records",
			fields: fields{
				records: map[string][]byte{
					"test.com.":  []byte("some test data"),
					"test2.com.": []byte("some test data"),
					"test3.com.": []byte("some test data"),
				},
				mu: sync.RWMutex{},
			},
			args: args{
				ctx: context.Background(),
			},
			wantList: []string{"test.com.", "test2.com.", "test3.com."},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InMemCache{
				records: tt.fields.records,
				mu:      tt.fields.mu,
			}
			gotList, err := i.ListAllRecords(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAllRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("ListAllRecords() gotList = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}
