package thesaurus

import (
	"os"
	"testing"
)

func TestBigHuge(t *testing.T) {
	apiKey := os.Getenv("BHT_APIKEY")
	testobj := &BigHuge{APIKey: apiKey}

	_, err := testobj.Synonyms("git")
	if err != nil {
		t.Fatal("error", err)
	}
}
