package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/sectwo/stcoin/blockchain"
)

const (
	port        string = ":8080"
	templateDir string = "templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(w http.ResponseWriter, r *http.Request) {
	//tmpl := template.Must(template.ParseFiles("templates/home.html"))
	// tmpl := template.Must(template.ParseFiles("templates/pages/home.html"))
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(w, "home", data)
}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.html"))
	templates = template.Must(template.ParseGlob(templateDir + "partials/*.html"))
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}

// chain := blockchain.GetBlockchain()
// chain.AddBlock("Second Block")
// chain.AddBlock("Third Block")
// chain.AddBlock("Fourth Block")

// for _, block := range chain.AllBlocks() {
// 	fmt.Println("Data : ", block.Data)
// 	fmt.Println("Hash : ", block.Hash)
// 	fmt.Println("PrevHash : ", block.PrevHash)
// 	fmt.Println()
// }
