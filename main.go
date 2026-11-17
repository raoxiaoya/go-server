package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/raoxiaoya/go-server/internal"
)

func main() {
	serverConfigList, err := internal.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for _, config := range serverConfigList {
		c := config
		InitProxy(c)
		go StartServer(c)
	}

	StartProxy()

	for {

	}
}

var needNginxProxy bool

func InitProxy(config internal.Config) {
	if config.ProxyPort == "" {
		return
	}
	if _, err := strconv.Atoi(config.ProxyPort); err != nil {
		return
	}
	arr := strings.Split(config.CmdStr, " ")
	if len(arr) < 3 {
		return
	}
	arr = strings.Split(arr[2], ":")
	if len(arr) < 2 {
		return
	}
	port := arr[1]
	if _, err := strconv.Atoi(port); err != nil {
		return
	}
	
	dir := "D:/phpStudy/PHPTutorial/nginx-1.28.2/conf/conf.d"
	confFile := dir + "/" + config.Name + "-proxy.conf"
	file, err := os.OpenFile(confFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	// 写入配置
	_, err = file.WriteString(fmt.Sprintf(`
server {
	listen       %s;
	server_name  0.0.0.0;

	location / {
		proxy_pass http://127.0.0.1:%s;
		
		add_header Access-Control-Allow-Headers *;
		add_header Access-Control-Allow-Origin *;
		add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
	}
}
`, config.ProxyPort, port))
	if err != nil {
		log.Println(err)
		return
	}
	
	needNginxProxy = true
	log.Println("init proxy conf file success:", confFile)
}

func StartServer(config internal.Config) {
	cmdStr := strings.Split(config.CmdStr, " ")
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = config.Dir
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}
}

func StartProxy() {
	if !needNginxProxy {
		return
	}

	// 强制kill全部的nginx进程
	cmdStr1 := strings.Split("taskkill /f /t /im nginx.exe", " ")
	cmd1 := exec.Command(cmdStr1[0], cmdStr1[1:]...)
	cmd1.Run()

	// 启动nginx.exe
	// 此处不能直接执行 start nginx 命令，会报错，因为 start 不是独立的可执行文件，而是 cmd.exe 的内置命令
	cmd2 := exec.Command("cmd", "/C", "start nginx")
	cmd2.Dir = "D:/phpStudy/PHPTutorial/nginx-1.28.2"
	if err := cmd2.Start(); err != nil {
		log.Println(err)
		return
	}

	log.Println("start nginx proxy success")
}