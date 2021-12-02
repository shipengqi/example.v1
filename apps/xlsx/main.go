package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/xuri/excelize/v2"
)

func main() {
	var xlsxPath, sheet string
	flag.StringVar(&xlsxPath, "xlsx", "", "xlsx file path")
	flag.StringVar(&sheet, "sheet", "", "sheet")
	flag.Parse()
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
		return
	}

	for k, row := range rows {
		if k == 0 {
			continue
		}
		dirname := row[0]
		var title string
		if len(row) >= 3 {
			title = row[2]
		}
		var content string
		if len(row) >= 5 {
			content = row[4]
		}
		var imgs string
		if len(row) >= 6 {
			imgs = row[5]
		}

		if len(title) == 0 {
			continue
		}
		var startStr string
		titles := strings.Split(title, " ")
		if len(titles) == 1 {
			startStr = "empty"
		} else {
			startStr = titles[0]
			startStr = strings.TrimSpace(startStr)
			startStr = strings.TrimSuffix(startStr, ":")
			if len(startStr) == 0 {
				startStr = "unknown"
			}
		}

		if exist := PathExists(dirname); !exist {
			if err = MkDirAll(dirname); err != nil {
				log.Fatal(err)
			}
		}
		if err != nil {
			log.Fatal(err)
		}
		filename := filepath.Join(dirname, title)
		log.Printf("write %s\n", filename)
		err = ioutil.WriteFile(filename, []byte(content), 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
		if len(imgs) > 0 {

			urls := strings.Split(imgs, "\n")
			for _, url := range urls {
				err = download(dirname, url, startStr)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}

		log.Println()
	}
}

func download(dirname, url, num string) error {
	url = strings.TrimSuffix(url, "_x000d_")
	if len(strings.TrimSpace(url)) == 0 {
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	baseUrl := filepath.Base(url)
	filename := filepath.Join(dirname, fmt.Sprintf("%s-%s", num, baseUrl))
	log.Printf("download %s\n", filename)
	return ioutil.WriteFile(filename, data, 0644)
}

func MkDirAll(fpath string) error {
	return os.MkdirAll(fpath, os.ModePerm)
}

// PathExists whether the path exists.
func PathExists(fpath string) bool {
	if fpath == "" {
		return false
	}
	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
