package main
// Import Go classes
import (
	"log"
	"net/http"	// Listens to root path
	"text/template"
	"path/filepath"
	"sync"
	"flag"
)

//	templ represents a single template 
//	sync.Once garantees the template only loads once
type templateHandler struct {
	once		sync.Once   
	filename	string
	templ 		*template.Template
}
//	ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()	// parse the flags
	r := newRoom()
	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	//	get the room going
	go r.run()



	//	Writes out the hardcoded HTML when a request is made
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head>
					<title>Chat</title>
				</head>
				<body>
					Let's chat!
				</body>
			</html>
		`))
	})*/
	// 	Starts web server using ListenAndServe
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}