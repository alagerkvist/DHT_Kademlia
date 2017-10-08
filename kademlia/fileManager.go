package kademlia

import (
	"io/ioutil"
	"os"
	"fmt"
	"encoding/base64"
	"time"
	"log"
)

type FileManager struct{
	encoder *base64.Encoding
	filesStored map[string]*FileInfo
}

type FileInfo struct {
	fileName string
	lastOriginalRefreshedStored time.Time
	lastTimeRefreshed time.Time
	initialStore time.Time
	expirationTime float64
	originalStore bool
	immutable bool
}

const filesDirectory = "kademlia/Files/"


func (fileManager *FileManager) CheckAndStore(fileName string, data string) {
	_, err := ioutil.ReadFile(filesDirectory + fileName)


	if os.IsNotExist(err){
		d1, err := base64.StdEncoding.DecodeString(data)
		if err != nil{
			fmt.Println("Error while decoding file")
		}

		f, err := os.Create(filesDirectory + fileName)
		if err != nil{
			fmt.Println(err)
		}

		_, err = f.Write(d1)
		file := FileInfo{fileName, time.Now().Local(), time.Now().Local(), time.Now().Local(), 24.0, false, false}
		fileManager.filesStored[fileName] = &file
		if err != nil{
			fmt.Println("\n !!! Error while writing in the file: ")
			fmt.Println(err)
		}
		defer f.Close()
	}
}


func (fileManager *FileManager) updateTime(fileName string){
	fileInfo := fileManager.filesStored[fileName]
	fileInfo.lastTimeRefreshed = time.Now().Local()
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

func (fileManager *FileManager) RemoveFile(filename string){
	delete(fileManager.filesStored, filename)
	os.Remove(filesDirectory + filename)
}


func (kademlia *Kademlia) CheckFiles(){

	for{
		time.Sleep(1 * time.Minute)
		fmt.Println("Check for files")
		for k, file := range kademlia.network.FileManager.filesStored {

				//Refreshing each file that has not been refresh from one hour
			if time.Since(file.lastTimeRefreshed).Hours() >= 1 {
				kademlia.Store(k)

				//Refreshing files owned, each 24h
			} else if file.originalStore && time.Since(file.lastOriginalRefreshedStored).Hours() >= 24{
				file.lastOriginalRefreshedStored = time.Now().Local()
				kademlia.Store(k)

				//Delete expirated files
			} else if !file.immutable && time.Since(file.initialStore).Hours() >= file.expirationTime{
				kademlia.network.FileManager.RemoveFile(k)
			}
		}
	}
}

func ListFiles(){
	files, err := ioutil.ReadDir("./kademlia/Files")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func (fileManager *FileManager) pinFile(fileName string, pin bool){
	file, ok := fileManager.filesStored[fileName]
	if !ok {
		fmt.Println("This file does not exist")
	} else {
		file.immutable = pin
	}
}