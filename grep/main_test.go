package main

import (
	"io"
	"reflect"
	"testing"
)

func Test_search(t *testing.T) {
	type args struct {
		r io.Reader
		c config
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := search(tt.args.r, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("search() = %v, want %v", got, tt.want)
			}
		})
	}
}
