# httputils

## Native Go usage
   http.FileServer(http.FileSystem(http.Dir("MY_DIR")))

## Package usage
   http.FileServer(httputils.FileSystem(http.Dir("MY_DIR")))
