
// Command-line tool to add, remove, or list webseeds in a .torrent file.
package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/jackpal/bencode-go"
)

// Torrent represents a minimal .torrent file structure for webseed editing.
type Torrent struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Comment      string     `bencode:"comment,omitempty"`
	CreatedBy    string     `bencode:"created by,omitempty"`
	CreationDate int64      `bencode:"creation date,omitempty"`
	Info         struct {
		Name        string `bencode:"name"`
		PieceLength int    `bencode:"piece length"`
		Pieces      string `bencode:"pieces"`
		Length      int    `bencode:"length,omitempty"`
		Files       []struct {
			Length int      `bencode:"length"`
			Path   []string `bencode:"path"`
		} `bencode:"files,omitempty"`
	} `bencode:"info"`
	UrlList []string `bencode:"url-list"`
}


// addWebseed adds a webseed URL if not already present.
func addWebseed(torrent *Torrent, webseedURL string) {
	for _, url := range torrent.UrlList {
		if url == webseedURL {
			fmt.Println("Webseed URL already exists in the torrent file.")
			return
		}
	}
	torrent.UrlList = append(torrent.UrlList, webseedURL)
	fmt.Println("Webseed URL added successfully.")
}


// removeWebseed removes a webseed URL if present.
func removeWebseed(torrent *Torrent, webseedURL string) {
	for i, url := range torrent.UrlList {
		if url == webseedURL {
			torrent.UrlList = append(torrent.UrlList[:i], torrent.UrlList[i+1:]...)
			fmt.Println("Webseed URL removed successfully.")
			return
		}
	}
	fmt.Println("Webseed URL not found in the torrent file.")
}


// listWebseeds prints all webseed URLs in the torrent.
func listWebseeds(torrent *Torrent) {
	if len(torrent.UrlList) == 0 {
		fmt.Println("No webseeds found in the torrent file.")
		return
	}
	fmt.Println("Webseeds in the torrent file:")
	for _, url := range torrent.UrlList {
		fmt.Println(url)
	}
}


func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-a|-r|-l] [torrent_file] [webseed_url]\n", os.Args[0])
}

func main() {
	addFlag := flag.Bool("a", false, "Add a webseed URL")
	removeFlag := flag.Bool("r", false, "Remove a webseed URL")
	listFlag := flag.Bool("l", false, "List all webseed URLs")
	flag.Usage = usage
	flag.Parse()

	if (*addFlag && *removeFlag) || (*addFlag && *listFlag) || (*removeFlag && *listFlag) || (!*addFlag && !*removeFlag && !*listFlag) {
		usage()
		os.Exit(2)
	}

	torrentFile := flag.Arg(0)
	webseedURL := flag.Arg(1)

	if torrentFile == "" || ((*addFlag || *removeFlag) && webseedURL == "") {
		usage()
		os.Exit(2)
	}

	file, err := os.Open(torrentFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var torrent Torrent
	err = bencode.Unmarshal(file, &torrent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing torrent file: %v\n", err)
		os.Exit(1)
	}

	if *listFlag {
		listWebseeds(&torrent)
		return
	}

	if *addFlag {
		addWebseed(&torrent, webseedURL)
	} else if *removeFlag {
		removeWebseed(&torrent, webseedURL)
	}

	file, err = os.Create(torrentFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	err = bencode.Marshal(file, torrent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving torrent file: %v\n", err)
		os.Exit(1)
	}
}
