## Editing page
 
```go
func editHandler(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/edit/"):]
    p, err := load(title)
    if err != nil {
        p = &Page{Title: title}
    }
    fmt.Fprintf(w, "<h1>Editing %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
        "<textarea name=\"body\">%s</textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>",
        p.Title, p.Title, p.Body)
}
```
直接在裡面寫 html 太 hardcored 了，

其實可以 import go 內建的 `html/template` package來解決這件事情

```html
<h1>Editing {{.Title}}</h1>
<form action="/save/{{.Title}}" method="post">
    <div>
        <textarea name="body" cols="30" rows="10">
            {{printf "%s" .Body}}
        </textarea>
    </div>
    <div>
        <input type="submit" value="Save">
    </div>
</form>
```

```go
func editHandler(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/edit/"):]
    p, err := load(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("edit.html")
    t.Execute(w, p)
}
```

`template.ParseFiles`會回傳一個`*template.Template`

edit.html 裡面的`{{.Title}}`會從 Execute 的第二個參數去找
（也就是 Page struct 的 Title）

而 view 也能夠做到一樣的事情

```go
func viewHandler(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/view/"):]
    p, err := load(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("view.html")
    t.Execute(w, p)
}
```

注意到這裏兩邊會出現重複的 code，

可將重複的部分抽出