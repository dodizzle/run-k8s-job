package main

import (
	"testing"
)

func TestBuildK8sConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		input   ActionInput
		wantErr bool
	}{}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {

			_, err := BuildK8sConfig(tt.input)

			if didErr := err != nil; didErr != tt.wantErr {
				t.Errorf("'wantErr' was %t, but err value was: '%v'", tt.wantErr, err)
			}
		})
	}
}
