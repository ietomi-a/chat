package main
import ( 
	"strings"
	"fmt"
	"log"
	"net/http" 
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP( w http.ResponseWriter, r *http.Request ){
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 認証されていない場合.
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}else if err != nil {
		panic( err.Error() )
	}else{
		h.next.ServeHTTP(w,r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler( w http.ResponseWriter, r *http.Request ){
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("TODO: ろぐいん処理", provider)
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証ぷろばいだーの取得に失敗しました。:", provider, "-", err )
		}
		loginUrl, err := provider.GetBeginAuthURL(nil,nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にえらーがしょうじました。:", provider, "-", err )
		}
		w.Header().Set("Location", loginUrl )
		w.WriteHeader( http.StatusTemporaryRedirect )

	case "callback":
		fmt.Print("callback okok\n")
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証ぷろばいだーの取得に失敗しました。:", provider, "-", err )
		}

		creds, err := provider.CompleteAuth( objx.MustFromURLQuery(r.URL.RawQuery) )
		if err != nil {
		 	log.Fatalln("認証を完了できませんでした:", provider, "-", err )
		}
		user, err := provider.GetUser(creds)
		if err != nil {
		 	log.Fatalln("user の取得に失敗しました:", provider, "-", err )
		}
		authCookieValue := objx.New(map[string]interface{}{ "name": user.Name(), } ).MustBase64()
		http.SetCookie(w, &http.Cookie{ Name: "auth", Value: authCookieValue, Path: "/"} )
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader( http.StatusTemporaryRedirect )

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "あくしょん%sには非対応です", action )
	}
}
