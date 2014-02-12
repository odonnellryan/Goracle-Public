package main

import (
	"testing"
)

func TestUpdateAllMongoDockerHostsInCollection(t *testing.T) {
	err := UpdateAllMongoDockerHostsInCollection()
	if err != nil {
		t.Errorf("Returned error: %s", err)
	}
}

