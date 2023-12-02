package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/zserge/lorca"
)

const htmlFilePath = "./res/index.html"
const iconFilePath = "./icon/dark/favicon.ico"

const (
	refreshMenuID = iota + 1
	restartMenuID
	exitMenuID
)

func main() {
	ui, err := lorca.New("", "", 800, 600)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	go func() {
		httpDir := filepath.Dir(htmlFilePath)
		httpPort := 808
		err := httpServer(httpDir, httpPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		htmlURL := fmt.Sprintf("http://localhost:808/%s", filepath.Base(htmlFilePath))
		err := ui.Load(htmlURL)
		if err != nil {
			log.Fatal(err)
		}
	}()

	systray.Run(onReady, onExit)
}

func onReady() {
	icon, err := os.Open(iconFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer icon.Close()

	iconBytes, err := io.ReadAll(icon)
	if err != nil {
		log.Fatal(err)
	}

	systray.SetIcon(iconBytes)

	refreshMenu := systray.AddMenuItem("刷新网页", "刷新 Lorca 网页")
	restartMenu := systray.AddMenuItem("重启程序", "重启 Go 程序")
	exitMenu := systray.AddMenuItem("退出程序", "退出 Go 程序")

	go func() {
		for {
			select {
			case <-refreshMenu.ClickedCh:
				fmt.Println("刷新 Lorca 网页")
				// 这里执行刷新操作
			case <-restartMenu.ClickedCh:
				fmt.Println("重启 Go 程序")
				// 这里执行重启操作
			case <-exitMenu.ClickedCh:
				fmt.Println("退出 Go 程序")
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// 清理工作
}

func httpServer(dir string, port int) error {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}
