package models

import "testing"

func TestAddLocatorShare(t *testing.T) {
	err := AddLocatorShare(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetLocatorShares(t *testing.T) {
	_, err := GetLocatorShares(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetLocatorResponses(t *testing.T) {
	_, err := GetLocatorResponses(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteLocatorShare(t *testing.T) {
	err := DeleteLocatorShare(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}
