package util

import (
	c "ProcessEngine/src/constant"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func SaveUniqueFile(fileName string, fileData []byte) (error, string) {
	if fileName == "" {
		return errors.New("filename cannot be null"), ""
	}
	// filename "filename.suffix"
	cur, _ := os.Getwd()
	fName := GetSp() + GenerateUnq() + GetSp() + fileName
	fPath := cur + GetSp() + c.CONFIG_RES_PATH + fName
	err1 := os.MkdirAll(filepath.Dir(fPath), 0755)
	if err1 != nil {
		log.Println("file err1:", err1.Error())
		return err1, fName
	}
	f, err2 := os.Create(fPath)
	if err2 != nil {
		log.Println("file err2:", err2.Error())
		return err2, fName
	}
	defer f.Close()
	err3 := ioutil.WriteFile(fPath, fileData, 0666)
	if err3 != nil {
		log.Println("file err3:", err3.Error())
		return err3, fName
	}
	return err3, fName
}
func GenerateUnq() string {
	t := sha1.New()
	io.WriteString(t, strconv.FormatInt(time.Now().Unix(), 10))
	return fmt.Sprintf("%x", t.Sum(nil))
}
