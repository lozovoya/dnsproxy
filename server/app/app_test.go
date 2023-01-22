package app

import (
	"dnproxier/server/cache"
	"dnproxier/server/pool"
	"reflect"
	"testing"
)

func TestApp_parseMessage(t *testing.T) {
	data := []byte{0, 25, 29, 194, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 118, 97, 111, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	type fields struct {
		pool       pool.ConnectionPoolInterface
		dnsCache   cache.DNSCacheInterface
		listenPort string
	}
	type args struct {
		buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   uint16
		wantErr bool
	}{
		{
			name: "parse correct dns message",
			args: args{
				buffer: data,
			},
			want:    "vao.com.",
			want1:   7618,
			wantErr: false,
		},
		{
			name: "parse wrong dns message",
			args: args{
				buffer: []byte{1, 2, 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				pool:       tt.fields.pool,
				dnsCache:   tt.fields.dnsCache,
				listenPort: tt.fields.listenPort,
			}
			got, got1, err := a.parseMessage(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestApp_setResponseID(t *testing.T) {
	data := []byte{0, 25, 29, 194, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 118, 97, 111, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	response := []byte{0, 25, 0, 111, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 118, 97, 111, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	type fields struct {
		pool       pool.ConnectionPoolInterface
		dnsCache   cache.DNSCacheInterface
		listenPort string
	}
	type args struct {
		buffer []byte
		id     uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "set for correct dns message",
			args: args{
				buffer: data,
				id:     111,
			},
			want:    response,
			wantErr: false,
		},
		{
			name: "set for wrong dns message",
			args: args{
				buffer: []byte{1, 2, 3},
				id:     111,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				pool:       tt.fields.pool,
				dnsCache:   tt.fields.dnsCache,
				listenPort: tt.fields.listenPort,
			}
			got, err := a.setResponseID(tt.args.buffer, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("setResponseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setResponseID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
