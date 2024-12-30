package client

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	endPoint := "http://localhost:8080"
	data := url.Values{}

	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)
	longURL, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	longURL = strings.TrimSuffix(longURL, "\n")
	data.Set("url", longURL)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, endPoint, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
