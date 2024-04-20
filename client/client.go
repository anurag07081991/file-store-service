package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	serverAddress string
)

func addFiles(files []string) {
	for _, file := range files {
		// Check if the file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("File '%s' not found\n", file)
			continue
		}

		// Open the file
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file '%s': %v\n", file, err)
			continue
		}
		defer f.Close()

		// Read file content
		content, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Printf("Error reading file '%s': %v\n", file, err)
			continue
		}

		// Send file content to the server
		resp, err := http.Post(serverAddress+"/add/"+file, "text/plain", bytes.NewReader(content))
		if err != nil {
			fmt.Printf("Error sending file '%s' to server: %v\n", file, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("File '%s' added successfully\n", file)
		} else {
			fmt.Printf("Failed to add file '%s'. Server returned: %s\n", file, resp.Status)
		}
	}
}

func listFiles() {
	resp, err := http.Get(serverAddress + "/ls")
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to list files. Server returned: %s\n", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println(string(body))
}

func removeFile(file string) {
	resp, err := http.Get(serverAddress + "/rm/" + file)
	if err != nil {
		fmt.Printf("Error removing file '%s': %v\n", file, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("File '%s' removed successfully\n", file)
	} else {
		fmt.Printf("Failed to remove file '%s'. Server returned: %s\n", file, resp.Status)
	}
}

func updateFile(file string) {
	// Check if the file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("File '%s' not found\n", file)
		return
	}

	// Open the file
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening file '%s': %v\n", file, err)
		return
	}
	defer f.Close()

	// Read file content
	content, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", file, err)
		return
	}

	// Send file content to the server
	resp, err := http.Post(serverAddress+"/update/"+file, "text/plain", bytes.NewReader(content))
	if err != nil {
		fmt.Printf("Error sending file '%s' to server: %v\n", file, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("File '%s' updated successfully\n", file)
	} else {
		fmt.Printf("Failed to update file '%s'. Server returned: %s\n", file, resp.Status)
	}
}

func wordCount() {
	resp, err := http.Get(serverAddress + "/wc")
	if err != nil {
		fmt.Printf("Error getting word count: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get word count. Server returned: %s\n", resp.Status)
		return
	}

	count, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Printf("Total number of words: %s\n", string(count))
}

func freqWords(limit int, order string) {
	url := fmt.Sprintf("%s/freq-words?limit=%d&order=%s", serverAddress, limit, order)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting frequent words: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get frequent words. Server returned: %s\n", resp.Status)
		return
	}

	words, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Println(string(words))
}

func main() {
	flag.StringVar(&serverAddress, "server", "http://localhost:8080", "Server address")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage: client <command>")
		fmt.Println("Commands:")
		fmt.Println("  add <file1> [<file2> ...]")
		fmt.Println("  ls")
		fmt.Println("  rm <file>")
		fmt.Println("  update <file>")
		fmt.Println("  wc")
		fmt.Println("  freq-words [--limit|-n 10] [--order=dsc|asc]")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "add":
		files := os.Args[2:]
		addFiles(files)
	case "ls":
		listFiles()
	case "rm":
		if len(os.Args) != 3 {
			fmt.Println("Usage: client rm <file>")
			os.Exit(1)
		}
		removeFile(os.Args[2])
	case "update":
		if len(os.Args) != 3 {
			fmt.Println("Usage: client update <file>")
			os.Exit(1)
		}
		updateFile(os.Args[2])
	case "wc":
		wordCount()
	case "freq-words":
		limit := 10
		order := "asc"

		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--limit", "-n":
				limit, _ = strconv.Atoi(os.Args[i+1])
				i++
			case "--order":
				order = os.Args[i+1]
				i++
			}
		}

		freqWords(limit, order)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
