


package main

import(
	"fmt"
	"os"
	"strconv"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/requestAndResponse", requestAndResponse)
	http.HandleFunc("/getVersion", getVersion)
	http.HandleFunc("/ipAndStatus", ipAndStatus)
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(":80", nil)
	if nil != err{
		log.Fatal(err)
	}
}

func requestAndResponse(response http.ResponseWriter, request *http.Request) {
	fmt.Println("/requestAndResponse")
	headers := request.Header
	for header := range headers {
		values := headers[header]
		for index, _ := range values {
			values[index] =strings.TrimSpace(values[index])
		}
		response.Header().Set(header, strings.Join(values, ","))
	}

	fmt.Fprintln(response, "All header data: ", headers)
	io.WriteString(response, "succeed")
}


func getVersion(response http.ResponseWriter, request *http.Request) {
	fmt.Println("/getVersion")
	envStr := os.Getenv("VERSION")
	response.Header().Set("VERSION", envStr)
	io.WriteString(response, "succeed")
}

func ipAndStatus(response http.ResponseWriter, request *http.Request){
	fmt.Println("/ipAndStatus")

	form := request.RemoteAddr
	fmt.Println("Client->ip:port=", form)
	ipStr := strings.Split(form, ":")
	fmt.Println("Client->ip=", ipStr[0])

	fmt.Println("Client->response code=", strconv.Itoa(http.StatusOK))
	io.WriteString(response, "succeed")

}

func healthz(response http.ResponseWriter, request *http.Request){
	fmt.Println("/healthz")
	response.WriteHeader(200)
	io.WriteString(response, "succeed")
}
