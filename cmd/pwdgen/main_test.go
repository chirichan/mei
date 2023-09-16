package main

import "testing"

func TestFullPassword(t *testing.T) {
	type args struct {
		level  int
		length int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{level: 4, length: 32},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fullPassword(tt.args.level, tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("fullPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.args.length {
				t.Errorf("fullPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
