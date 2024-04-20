package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func handleAdd(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/add/")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Check if the file already exists
	if _, err := os.Stat(filepath.Base(file)); !os.IsNotExist(err) {
		http.Error(w, "File already exists", http.StatusConflict)
		return
	}

	// Create the file
	f, err := os.Create(filepath.Base(file))
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write file content
	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, "Failed to write file content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	var fileList []string
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	response := strings.Join(fileList, "\n")
	w.Write([]byte(response))
}

func handleRemove(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/rm/")
	err := os.Remove(file)
	if err != nil {
		http.Error(w, "Failed to remove file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimPrefix(r.URL.Path, "/update/")

	// Create or open the file
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write file content
	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, "Failed to write file content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleWordCount(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	var wordCount int
	for _, file := range files {
		if !file.IsDir() {
			content, err := ioutil.ReadFile(file.Name())
			if err != nil {
				http.Error(w, "Failed to read file content", http.StatusInternalServerError)
				return
			}
			wordCount += countWords(string(content))
		}
	}

	response := fmt.Sprintf("%d", wordCount)
	w.Write([]byte(response))
}

func handleFreqWords(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	order := r.URL.Query().Get("order")

	files, err := ioutil.ReadDir(".")
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	var words []string
	for _, file := range files {
		if !file.IsDir() {
			content, err := ioutil.ReadFile(file.Name())
			if err != nil {
				http.Error(w, "Failed to read file content", http.StatusInternalServerError)
				return
			}
			words = append(words, strings.Fields(string(content))...)
		}
	}

	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	sortedWords := make([]string, 0, len(wordFreq))
	for word := range wordFreq {
		sortedWords = append(sortedWords, word)
	}

	sort.Slice(sortedWords, func(i, j int) bool {
		if order == "asc" {
			return wordFreq[sortedWords[i]] < wordFreq[sortedWords[j]]
		}
		return wordFreq[sortedWords[i]] > wordFreq[sortedWords[j]]
	})

	if limit > len(sortedWords) {
		limit = len(sortedWords)
	}

	var result strings.Builder
	for i := 0; i < limit; i++ {
		result.WriteString(fmt.Sprintf("%d\t%s\n", wordFreq[sortedWords[i]], sortedWords[i]))
	}

	w.Write([]byte(result.String()))
}

func countWords(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

func setupRoutes() {
	http.HandleFunc("/add/", handleAdd)
	http.HandleFunc("/ls", handleList)
	http.HandleFunc("/rm/", handleRemove)
	http.HandleFunc("/update/", handleUpdate)
	http.HandleFunc("/wc", handleWordCount)
	http.HandleFunc("/freq-words", handleFreqWords)
}

func main() {
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
