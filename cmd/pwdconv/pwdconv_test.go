package main

import (
	"testing"
)

func TestLoadDomainList(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadDomainList()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDomainList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
			for _, s := range got {
				t.Logf("%s\n", s)
			}
		})
	}
}

func TestParseHost(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				host: "",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseHost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
