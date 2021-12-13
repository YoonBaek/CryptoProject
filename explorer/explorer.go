package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/YoonBaek/CryptoProject/blockchain"
)

const (
	templateDir string = "explorer/template/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("template/pages/home.gohtml"))
	data := homeData{PageTitle: "SuperSexy", Blocks: nil}
	err := templates.ExecuteTemplate(rw, "home", data)
	if err != nil {
		log.Fatal(err)
	}
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := templates.ExecuteTemplate(rw, "add", nil)
		if err != nil {
			log.Fatal(err)
		}
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.BlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(portNum int) {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d\n", portNum)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portNum), nil))
}
