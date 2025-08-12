package main

import (
	"os"
	"testing"
)

func TestAddWebseed(t *testing.T) {
	torrent := &Torrent{}
	addWebseed(torrent, "http://example.com/seed1")
	if len(torrent.UrlList) != 1 || torrent.UrlList[0] != "http://example.com/seed1" {
		t.Errorf("expected webseed to be added")
	}
	// Adding duplicate
	addWebseed(torrent, "http://example.com/seed1")
	if len(torrent.UrlList) != 1 {
		t.Errorf("duplicate webseed should not be added")
	}
}

func TestRemoveWebseed(t *testing.T) {
	torrent := &Torrent{UrlList: []string{"http://example.com/seed1", "http://example.com/seed2"}}
	removeWebseed(torrent, "http://example.com/seed1")
	if len(torrent.UrlList) != 1 || torrent.UrlList[0] != "http://example.com/seed2" {
		t.Errorf("expected webseed to be removed")
	}
	// Remove non-existent
	removeWebseed(torrent, "http://notfound.com/")
	if len(torrent.UrlList) != 1 {
		t.Errorf("removing non-existent webseed should not change list")
	}
}

func TestListWebseeds(t *testing.T) {
	// Just check it doesn't panic
	torrent := &Torrent{UrlList: []string{"http://example.com/seed1"}}
	listWebseeds(torrent)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
