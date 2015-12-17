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

## Handle the non-existing page

如果 /view/nonexist，

還是會跑到 view 頁面，

但實際上這個 page並沒有被存下來。

應該要將它導到 /edit/nonexist，

創造一個屬於他的頁面才對。

```go
func viewHandler(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/view/"):]
    p, err := load(title)
    if err != nil {
        http.Redirect(w, req, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view.html", p)
}
```

這裏需要改寫 viewHanlder，加入http.Redirect

## Saving page

```go
func saveHandler(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/save/"):]
    body := req.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, req, "/view/"+title, http.StatusFound)
}
```

body 是從 edit.html 裡面的 form 去拿的。

## Error Handling

前面的例子大部份省略錯誤處理

但實際上我們並不希望錯誤發生的時候卻什麼也不知道

而且這樣程式也會相當不 robust

讓我們按照 golang 中的習慣寫法來改寫`renderTemplate`

```go
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        // the error msg should be plain text
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

http.Error會幫我們處理。
> 查完文件後
> 第二個參數一定要放 plain text，第三個則是http package中的constant: 500

##Template Cacheing

前面每次 render template時的會重新 parse一次

這樣做其實是沒有效率的

最好是能只parse一次，接著有需要被用到時

就從parse過的 templates 中挑要的那個就行了（呼叫的是*Template）

> 實作方法就是我們拿到 *Template以後，再用[ExecuteTemplate](https://golang.org/pkg/html/template/#Template.ExecuteTemplate)來 render


##Validation

MustCompile is distinct from Compile in that it will panic if the expression compilation fails, while Compile returns an error as a second parameter.



