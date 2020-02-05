package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

type Book struct {
	ID     int    `json:"ID"`
	Title  string `json:"Title"`
	Author string `json:"Author"`
	ISBN   string `json:"Isbn"`
}

var (
	sh        = shell.NewShell("http://127.0.0.1:5001")
	db        = make(map[string]Book)
	masterCID = ""
)

func newBook(ID int, Title, Author, ISBN string) Book {
	return Book{
		ID,
		Title,
		Author,
		ISBN,
	}
}

//print menu
func showMenu() {
	fmt.Println("1 add book")
	fmt.Println("2 show all book")
	fmt.Println("3 show ipld address")
	fmt.Println("4 exit")
}

//wait for user menu input
func userInput() Book {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("please add new book info")
	fmt.Print("enter book name:")
	scanner.Scan()
	title := scanner.Text()
	fmt.Print("enter book author:")
	scanner.Scan()
	author := scanner.Text()
	fmt.Print("enter book isbn:")
	scanner.Scan()
	isbn := scanner.Text()
	b := newBook(len(db), title, author, isbn)
	return b
}

//add book
func addBook() (cid string) {
	book := userInput()
	db[book.ISBN] = book
	obj, err := json.Marshal(db)
	if err != nil {
		fmt.Println("error")
		return
	}
	masterCID, err = sh.DagPut(obj, "json", "cbor")
	return
}

//show book
func getBooks() (res map[string]Book, err error) {
	err = sh.DagGet(masterCID, &res)
	return
}

//exit os.Exit(0)
func main() {
	if masterCID == "" {
		fmt.Println("book store is empty")
		addBook()
	}
	for {
		db, err := getBooks()
		if err != nil {
			print("err")
			return
		}
		showMenu()
		var option int
		fmt.Scanln(&option)
		switch option {
		case 1:
			addBook()
		case 2:
			fmt.Println(db)
		case 3:
			fmt.Printf("http://127.0.0.1:8080/ipns/explore.ipld.io/#/explore/%s\n", masterCID)
		case 4:
			os.Exit(0)
		}
	}
}
