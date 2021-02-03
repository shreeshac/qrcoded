package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(rootHandler))
	http.Handle("/api", http.HandlerFunc(qrCodeGetHandler))
	http.Handle("/api/file", http.HandlerFunc(qrCodeFileGetHandler))
	fmt.Println("Listening to port", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func qrCodeGetHandler(w http.ResponseWriter, req *http.Request) {
	qrtext := req.FormValue("qrtext")
	qrsize := req.FormValue("qrsize")
	if qrsize == "" {
		qrsize = "128"
	}
	qrsizeStr, err := strconv.Atoi(qrsize)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	if png, err := qrcode.Encode(qrtext, qrcode.Medium, qrsizeStr); err != nil {
		fmt.Println("error", err)
	} else {
		w.Write(png)
	}

}

func qrCodeFileGetHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		return
	}
	qrtext := req.FormValue("qrtext")
	qrsize := req.FormValue("qrsize")
	qrFilePath := req.FormValue("path")
	if qrsize == "" {
		qrsize = "128"
	}
	if qrFilePath == "" {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprintf(w, "{ \"Message\":\"param `path` is required.\"}")

		return
	}
	qrsizeInt, err := strconv.Atoi(qrsize)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	if err := qrcode.WriteFile(qrtext, qrcode.Medium, qrsizeInt, qrFilePath); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprintf(w, "{ \"Message\":\"%v.\"}", err)
		fmt.Println("error", err)
	}

}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("qrtext"))

}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="/api?qrtext={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name="qrtext" value="" title="Text to QR Encode"><input type=submit
value="Show QR" >
</form>
</body>
</html>`
