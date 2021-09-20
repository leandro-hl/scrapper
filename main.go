package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	_ "strings"
	"time"
)

const emailRegex = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

func main() {
	//get urls from file
	data, err := ioutil.ReadFile("urls.txt")

	if err != nil {
		fmt.Println(err)
	}

	urls := strings.Split(string(data), "\r\n")

	re := regexp.MustCompile(emailRegex)

	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.ModeAppend)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	//we don't do it with go routines so we don't kill the server and not get banned.
	for _, url := range urls {
		res, err := http.Get(url)

		if err != nil {
			fmt.Println(err)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			fmt.Println(err)
			continue
		}

		mail := re.FindString(string(body))

		if mail != "" {
			fmt.Println(mail)
			file.Write([]byte(mail + "\r\n"))
		}
		time.Sleep(2 * time.Second)
	}
}
