# File Store Service

File Store Service is a simple file storage service implemented in Go. It provides an HTTP server and a command-line client to store, update, delete files, and perform various operations on files stored in the server.

## Features

- Add files to the store.
- List files in the store.
- Remove a file from the store.
- Update contents of a file in the store.
- Word count of all files stored in the server.
- List the most frequent words in all files combined.

## Usage

### Server

To run the server, use the following command:

```sh
go run server.go

Client
To run the client, use the following command:

go run client.go <command>

Replace <command> with one of the following commands:

add <file1> [<file2> ...]: Add files to the server.
ls: List files in the server.
rm <file>: Remove a file from the server.
update <file>: Update the contents of a file in the server.
wc: Get the total number of words in all files stored on the server.
freq-words [--limit|-n 10] [--order=dsc|asc]: Get the most frequent words in all files stored on the server.

go run client.go add file1.txt file2.txt
go run client.go ls
go run client.go update file1.txt
go run client.go wc
go run client.go freq-words --limit 10 --order asc
