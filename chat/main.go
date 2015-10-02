package main
// Import Go classes
import (
	"log"
	"net/http"	// Listens to root path
	"text/template"
	"path/filepath"
	"sync"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	//"../trace"
	//"os"
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
	//	set up gomniauth
	gomniauth.SetSecurityKey("this is my own crazy phrase not")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
		google.New("AIzaSyAHdC_P8iM2SU3D5BEh5747tGb4Sr5xxj8", "3aoJeQ8Ub3l2Gfvz-wNIOXUo", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	//	output to the os.Stdout standard output pipe (prints output to the terminal)
	//r.tracer = trace.New(os.Stdout)
	// root
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	//	Directory used to store css and js files.
	//http.Handle("/assets/",
	//	http.StripPrefix("/assets",
	//		http.FileServer(http.Dir("/path/to/assets/"))))


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