package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

func main() {
	// 定义命令行参数
	dirPath := flag.String("dir", ".", "指定目录路径")
	days := flag.Int("days", 30, "指定天数")

	// 解析命令行参数
	flag.Parse()

	// 获取当前时间
	now := time.Now()

	// 计算截止日期
	endDate := now.AddDate(0, 0, -*days)

	// 遍历指定目录下的文件
	files, err := ioutil.ReadDir(*dirPath)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return
	}

	for _, file := range files {
		// 检查文件修改时间
		if file.ModTime().Before(endDate) {
			// 获取文件路径
			filePath := filepath.Join(*dirPath, file.Name())

			// 检测文件类型
			mime, err := mimetype.DetectFile(filePath)
			if err != nil {
				fmt.Println("检测文件类型失败:", err)
				continue
			}

			// 创建目标文件夹路径
			destDir := filepath.Join(*dirPath, mime.String())

			// 如果目标文件夹不存在，创建它
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				os.Mkdir(destDir, os.ModeDir)
			}

			// 移动文件到目标文件夹
			newPath := filepath.Join(destDir, file.Name())
			err = os.Rename(filePath, newPath)
			if err != nil {
				fmt.Println("移动文件失败:", err)
			}
		}
	}
}


//使用方法：./your_program -dir /path/to/your/directory -days 30
