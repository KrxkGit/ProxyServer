package main

import (
	"encoding/json"
	"fmt"
	"github.com/KrxkGit/ProxyServer/Proxy"
	"log"
	"net"
	"os"
)

type Setting struct {
	Address string `json:"Address"`
	Port    string `json:"Port"`
}

func main() {
	setting := new(Setting)
	readSetting(setting)
	fmt.Println("Welcome to Use ProxyServer Developed by Krxk.")
	fmt.Printf("Listening at %s:%s", setting.Address, setting.Port)

	listener, err := net.Listen("tcp", setting.Address+":"+setting.Port)
	if err != nil {
		log.Println(err.Error())
	}
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		go Proxy.MainProxyHandler(client)
	}
}

func readSetting(setting *Setting) error {
	filePath := "setting.json"
	var err error
	if fi, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		defer file.Close()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		buf := make([]byte, fi.Size())
		n, err := file.Read(buf)
		if err != nil || int64(n) != fi.Size() {
			log.Println(err.Error())
			return err
		}
		json.Unmarshal(buf, setting)
		return nil
	}
	return err
}
