echo 'Serving static files on port :8080...'
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
