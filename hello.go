package hello

import ( 
	"net/http"
	"html/template"
	"appengine"
	"time"
	"appengine/datastore"
	"appengine/user"
)

type Greeting struct {
	Author string
	Content string
	Date time.Time
}

func init() {
	http.HandleFunc("/",root)
	http.HandleFunc("/sign",sign)
}

func root(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
        greetings := make([]Greeting, 0, 10)
        if _, err := q.GetAll(c, &greetings); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        if err := guestbookTemplate.Execute(w, greetings); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

const guestbookForm = `
<html>
  <body>
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`

func sign(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        g := Greeting{
                Content: r.FormValue("content"),
                Date:    time.Now(),
        }
        if u := user.Current(c); u != nil {
                g.Author = u.String()
        }
        key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
        _, err := datastore.Put(c, key, &g)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHTML))

const signTemplateHTML = `
<html>
  <body>
    <p>You wrote:</p>
    <pre>{{.}}</pre>
  </body>
</html>
`
func guestbookKey(c appengine.Context) *datastore.Key {
        return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))

// func handler(w http.ResponseWriter, r *http.Request) {
// 	c := appengine.NewContext(r)
// 	u := user.Current(c)
// 	if u == nil {
// 		url, err := user.LoginURL(c, r.URL.String())
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Location", url)
// 		w.WriteHeader(http.StatusFound)
// 		return
// 	}
// 	fmt.Fprintf(w,"Hello, %v!",u)
// }