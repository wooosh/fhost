package main

import (
    "fmt"
    "net/http"
    "io"
    "encoding/base64"
    "os"
    "io/ioutil"
    "time"
    "sync"

    "github.com/wooosh/base62"
)

var auth, address string
var idNext int64
var idMux sync.Mutex
var authEnabled bool

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    // Add colon to auth string so username is not required
    if os.Getenv("AUTH") != "" {
        authEnabled = true
        auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(":" + os.Getenv("AUTH")))
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "80"
    }

    if os.Getenv("WEBPATH") != "" {
        address = os.Getenv("WEBPATH") + "/files/"
    }

    files, err := ioutil.ReadDir("./files")
    check(err)
    if len(files) == 0 {
        idNext = time.Now().Unix()
    } else {
        // Set id to last created file
        idNext, err = base62.Decode(files[len(files)-1].Name())
        check(err)
    }

    http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))
    http.HandleFunc("/", Uploader)
    check(http.ListenAndServe(":" + port, nil))
}

func Uploader(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
        return
    }
    if r.Header.Get("Authorization") != auth {
        http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
        return
    }

    idMux.Lock()
    fileid := time.Now().Unix()
    // Check if a fileid has already been used
    if fileid <= idNext {
        fileid = idNext + 1
    }
    idNext = fileid
    idMux.Unlock()

    filename := base62.Encode(fileid)
    file, _  := os.Create("./files/" + filename)
    formfile, _, _ := r.FormFile("file")
    defer file.Close()
    defer formfile.Close()
    io.Copy(file, formfile)

    // Send file address back
    fmt.Fprintln(w, address + filename)
}
