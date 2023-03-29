# Snippet Box  

A snippet box created following Alex Edwards' *Let's Go*.  

## Instructions

### Running on dev 
```bash
go run ./cmd/web # from project root
```
Or simply, 
```bash
make serve
```

## Notes  
* Handlers: similar to controllers in MVC. Responsible for executing application logic and writing HTTP response headers and bodies.  
* Servemux: (router) stroes a mapping between the URL patterns and corresponding handlers. Usually, one servemux containing all routes per application.  
* Server: Go can establish a web server and listen for incoming requests natively without external third-party servers like Nginx or Apache.  

### Why not use http.HandleFunc? 

```go
func main() {
    http.HandleFunc("/", home) 
    http.HandleFunc("/snippet/view", snippetView) http.HandleFunc("/snippet/create", snippetCreate)
    log.Println("Starting server on :4000") err := http.ListenAndServe(":4000", nil) log.Fatal(err)
}
```

This shortens code slightly, but what it is doing is simply defining a NewServeMux() as a global variable and storing all handlers in said mux behind the hood. However, *because it defines a new servemux as global*, any package can access it and register a route-- including any third-party packages that the application imports. This could expose malicious handlers to the web. Thus, using http.HandleFunc is not advisable.  

### Servemux features and quirks  
* routes that end with "/" catch all their subtrees (equivalent to "/*"). Routes that do not end with "/" need to be matched exactly.
* Go will always match the handler with the longest path when there are two handlers that match a route. E.g., if there is "tree/" and "tree/subtree" mapped to different handlers, then the url "tree/subtree" will map to the handler with route "tree/subtree".  
*  if you have registered the subtree path `/foo/`, then any request to `/foo` will be redirected to `/foo/`, with a `301 Permanent Redirect`.  

### Host name matching  
```go
mux := http.NewServeMux() mux.HandleFunc("foo.example.org/baz", foobazHandler) mux.HandleFunc("bar.example.org/baz", barbazHandler) mux.HandleFunc("/baz", bazHandler)
```
It’s possible to include host names in your URL patterns. This can be useful when you want to redirect all HTTP requests to a canonical URL, or if your application is acting as the back end for multiple sites or services.  
When it comes to pattern matching, any host-specific patterns will be checked first and if there is a match the request will be dispatched to the corresponding handler. Only when there isn’t a host-specific match found will the non-host specific patterns also be checked.  

### RESTful routing  
* ServeMux is lightweight, doesn't support clean URLS with variables, regexp-based patterns, or request method.  
* Still sufficient for many applications.  

### Structuring  
Keep it simple.  

* The `cmd` directory will contain the application-specific code for the executable applications in the project. For now we’ll have just one executable application — the `web` application — which will live under the `cmd/web` directory.  
* The `internal` directory will contain the ancillary non-application-specific code used in the project. We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project. Packages in the `internal` directory in Go can *only be imported by code inside the project directory. It cannot be imported by code outside of the project.*
* The `ui` directory will contain the user-interface assets used by the web application. Specifically, the `ui/html` directory will contain HTML templates, and the `ui/static` directory will contain static files (like CSS and images).  

This structure cleanly separates Go and non-Go assets. Further, it scales nicely when adding another executable application to the project. For example, you might want to add a CLI (Command Line Interface) to automate some administrative tasks in the future. With this structure, you could create this CLI application under `cmd/cli` and it will be able to import and reuse all the code you’ve written under the internal directory.  

### Templating  
`html/template` allows recycled use of templates defined in ui/html.  
* block action. This acts like the {{template}} action, except it allows you to specify some default content if the template being invoked doesn’t exist in the current template set. In the context of a web application, this is useful when you want to provide some default content (such as a sidebar) which individual pages can override on a case-by-case basis if they need to. But — if you want — you don’t need to include any default content between the {{block}} and {{end}} actions. In that case, the invoked template acts like it’s ‘optional’. If the template exists in the template set, then it will be rendered. But if it doesn’t, then nothing will be displayed.
```html
{{define "base"}}
    <h1>An example template</h1> 
    {{block "sidebar" .}}
        <p>My default sidebar content</p> 
    {{end}}
{{end}}
```
