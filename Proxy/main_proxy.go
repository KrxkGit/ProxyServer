package Proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func MainProxyHandler(client net.Conn) {
	defer client.Close()

	// 获取请求头
	log.Printf("remote addr: %v\n", client.RemoteAddr())
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 分析请求头
	var method, URL, serverAddr string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &URL)
	hostPortURL, err := url.Parse(URL)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 连接 远程服务器
	if strings.Index(hostPortURL.Host, ":") == -1 {
		serverAddr = hostPortURL.Host + ":80" // 默认为 80 端口
	}
	fmt.Println("Forward to: ", serverAddr)
	server, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 分别处理 Http 与 Https
	if method == "CONNECT" {
		HandleHttps(client)
	} else {
		HandleHttp(server, b[:n])
	}
	io.Copy(client, server)
}

func HandleHttp(server net.Conn, buf []byte) {
	server.Write(buf)
}

func HandleHttps(client net.Conn) {
	fmt.Fprintf(client, "HTTP/1.1 200 Connection established\r\n\r\n")
}
