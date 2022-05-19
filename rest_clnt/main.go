package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"time"
)

type CreateDesktop struct {
	ComputerName string    `json:"cname"`
	Ip           string    `json:"cip"`
	UserName     string    `json:"user"`
	At           time.Time `json:"at"`
}

type desktopClient struct {
}

func NewDesktopClient() *desktopClient {
	return &desktopClient{}
}

// GetParams возвращает типичный набор параметров для отправки серверу
func (dc *desktopClient) GetParams() (*bytes.Buffer, error) {
	UserName, err := user.Current()
	if err != nil {
		return nil, err
	}
	HostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	Addrs, err := net.LookupHost(HostName)
	if err != nil || len(Addrs) == 0 {
		return nil, err
	}
	js, err := json.Marshal(CreateDesktop{
		ComputerName: HostName,
		Ip:           Addrs[0], // у ЭВМ м.б. несколько адресов, по хорошему надо искать исходящий отправляя запрос
		UserName:     UserName.Username,
		At:           time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(js), nil
}

// GetUrl тут могут быть какие-то особые формирователи url
func (dc *desktopClient) GetUrl(operation string) (string, string, error) {
	method := ""
	if operation == "register" {
		method = "POST"
	} else if operation == "unregister" {
		method = "DELETE"
	} else if operation == "update" {
		method = "GET"
	} else {
		return "", "", fmt.Errorf("unable to create url for %s method", operation)
	}
	return fmt.Sprintf("http://%s:%s/%s/", os.Getenv("SERVERURL"), os.Getenv("SERVERPORT"), operation), method, nil
}

func main() {
	client := NewDesktopClient()
	operation := "register"
	if len(os.Args) >= 2 {
		operation = os.Args[1] // разрешено : register, unregister, update
	}
	url, method, err := client.GetUrl(operation)
	if err != nil {
		log.Fatal(err)
		return
	}
	bb, err := client.GetParams()
	if err != nil {
		log.Fatal(err)
		return
	}
	req, err := http.NewRequest(method, url, bb)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	cl := &http.Client{}
	res, err := cl.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(string(body))
}
