package pkg

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"os"
	"strings"
)

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("open err:", err)
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func ReadDir(path string, fileInfoFunc func(path string, info os.FileInfo)) {
	//以只读的方式打开目录
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		fmt.Println(err.Error())
	}
	//延迟关闭目录
	defer f.Close()
	fileInfo, _ := f.Readdir(-1)
	for _, info := range fileInfo {
		fileInfoFunc(path, info)
	}
}

func ToTxt(str, fileName string) {
	fileName = "resource/" + fileName
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_SYNC|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666) //创建文件
	f.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(str + "\n")
}

func GetSub(res, start, end string) []string {
	result := make([]string, 0)
	ss := strings.Split(res, start)
	for i, v := range ss {
		if i > 0 {
			value := strings.Split(v, end)[0]
			result = append(result, value)
		}
	}
	if len(result) == 0 {
		ToTxt(res+"--------"+start+"--------"+end, "error.log")
	}
	result = append(result, "")
	return result
}

func WriteExcel() {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save xlsx file by the given path.

	/*if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}*/
	f.Save()
}
