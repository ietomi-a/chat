package main
import ( 
	"fmt"
	"sync"
	"text/template"
	"path/filepath"
	"net/http"
	"flag"
	"log"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP( w http.ResponseWriter, r *http.Request ){
	t.once.Do( func(){
		t.templ = 
			template.Must( 
				template.ParseFiles( filepath.Join("templates", t.filename) ) ) 
	} )
	data := map[string]interface{}{ "Host": r.Host, }
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w,data)
	//t.templ.Execute(w,r)
}

func main(){
	fmt.Print( "hello,world\n")
	var addr = flag.String("addr", ":8080", "亜ぷりけーしょん address")
	flag.Parse()
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		google.New( 
			"606208371926-dkhlooofse3hakqj2fkfrqlpv1vpuu9v.apps.googleusercontent.com",
			"wGcZPH5OnR1Jw4dPtOe1PqH_", 
			"http://localhost:8080/auth/callback/google", 
		),
	)

	r := newRoom(UseAuthAvatar)
	//r := newRoom()

	// http module への登録.
	http.Handle( "/chat", MustAuth(&templateHandler{filename: "chat.html"}) )
	//http.Handle( "/chat", &templateHandler{filename: "chat.html"} )
	http.Handle( "/login", &templateHandler{filename: "login.html"} )
	http.HandleFunc( "/auth/", loginHandler)
	http.Handle( "/room", r )

	http.HandleFunc( "/logout", func(w http.ResponseWriter, r *http.Request ){
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader( http.StatusTemporaryRedirect )
	})

	go r.run()
	log.Println("web start port", *addr )
	if err := http.ListenAndServe( *addr, nil); err != nil {
		fmt.Print("error!")
	} 
	fmt.Print("listen ok")
}
