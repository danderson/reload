# Live Reload for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/danderson/reload.svg)](https://pkg.go.dev/github.com/danderson/reload)

This package provides very simple live reloading for Go web
servers. It is intended purely for development, and is not specialized
to any web or client framework. The reload is literally an automated
push of the refresh button in the browser, nothing more.

To use, serve the reloader from some URL in your web server:

```go
rl := &reload.Reloader{}
mux.Handle("/.magic/live", rl)
```

Then reference it through a script tag in your HTML:

```html
<!DOCTYPE html>
<html>
  <head>
    <script src="/.magic/live"></script>
  </head>
  <body>
    <p>
      This page will refresh when the server changes, or when
      something visits /.magic/reload.
    </p>
  </body>
</html>
```

Every page that loads the script will reload itself whenever the
server restarts, or when the server invokes `Reload` on the live
handler. For example, you can use the latter to make a magic URL that
reloads all clients:

```go
mux.HandleFunc("/.magic/reload", func(w http.ResponseWriter, r *http.Request) {
  rl.Reload()
})
```

## Credit

Thank you to [Andy Dote](https://andydote.co.uk/) and his [blog post
about hot reloading in
Go](https://andydote.co.uk/2023/11/15/hot-reload-for-serverside-rendering/). I
had the exact problem his post solves, and this library is a mild
variation on the solution he presents there.
