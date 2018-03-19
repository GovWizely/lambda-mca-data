package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestGetMcaDataItem(t *testing.T) {
	feedItem := gofeed.Item{Title: "Some title", Description: "Some desc", Link: "http://example.com/1",
		Published: "Mon, 26 Feb 2018 21:19:31 EST", GUID: "unique",
		Categories: []string{"type/spn", "deadline/2018/Mar/23", "country/us - United States", "CPV/79212000",
			"CPV/66600000", "CPV/79212400"}}
	mcaDataItem := getMcaDataItem(&feedItem)
	if mcaDataItem.CountryCode != "US" {
		t.Errorf("Country Code was incorrect, got: %s, expected: %s.", mcaDataItem.CountryCode, "US")
	}
	if mcaDataItem.CountryName != "United States" {
		t.Errorf("Country Name was incorrect, got: %s, expected: %s.", mcaDataItem.CountryName, "United States")
	}
	categories := "CPV/66600000,CPV/79212000,CPV/79212400,deadline/2018/Mar/23,type/spn"
	if mcaDataItem.Categories != categories {
		t.Errorf("Categories was incorrect, got: %s, expected: %s.", mcaDataItem.Categories, categories)
	}
}
