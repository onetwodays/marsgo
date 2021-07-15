package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFile(path string)([]byte, error){
	f, err := os.Open(path)
	if err !=nil{
		return nil, errors.Errorf("%s","123")
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err !=nil{
		return nil, errors.Wrap(err,"read failed")
	}
	return buf,nil
}
func ReadConfig()([]byte, error){
	home := os.Getenv("HOME")
	config, err :=ReadFile(filepath.Join(home,".settings.xml"))
	return config, errors.WithMessage(err,"could not read config")
}
func main1(){
	_, err :=ReadConfig()
	if err !=nil{
		fmt.Println(err)
		os.Exit(1)
	}
}

func main()  {
	// 省略部分代码
	data := "test string"
	// md5.Sum() return a byte array
	h := md5.Sum([]byte(data))
	fmt.Println("h=",h)

	// with "%x" format byte array into hex string
	hexStr := fmt.Sprintf("%x", h)
	fmt.Println(hexStr)

	hx,err:=hex.DecodeString(hexStr)
	if err!=nil{
		fmt.Println(err)
	}else {
		fmt.Println("hx=",hx)
	}


}
