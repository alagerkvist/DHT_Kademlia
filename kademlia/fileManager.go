package kademlia

import (
	"io/ioutil"
	"os"
	"fmt"
	"encoding/base64"
)

type FileManager struct{
	encoder *base64.Encoding
}

const filesDirectory = "kademlia/Files/"

func (f *FileManager) checkAndStore(fileName string, data string) {
	_, err := ioutil.ReadFile(filesDirectory + fileName)


	if os.IsNotExist(err){
		d1, err := f.encoder.DecodeString(data)
		if err != nil{
			fmt.Println("Error while decoding file")
		}

		f, err := os.Create(filesDirectory + fileName)
		if err != nil{
			fmt.Println(err)
		}

		fmt.Println(d1)
		_, err = f.Write(d1)
		if err != nil{
			fmt.Println("\n !!! Error while writing in the file : ")
			fmt.Println(err)
		}
		defer f.Close()
	}
}

func (f *FileManager) checkIfFileExist(fileName string) bool{
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *FileManager) readData(fileName string) []byte{
	data, _ := ioutil.ReadFile(fileName)
	return data
}
