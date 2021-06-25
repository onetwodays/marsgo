package main

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"fmt"
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
func main(){
	_, err :=ReadConfig()
	if err !=nil{
		fmt.Println(err)
		os.Exit(1)
	}
}
