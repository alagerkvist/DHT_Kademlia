package kademlia

import (
	"io/ioutil"
	"os"
	"fmt"
	"encoding/base64"
	"time"
	"log"
)

/** FileManager
* filesStored: represent all the files stored in context of Kademlia
*				string: file name
*				FileInfo: see below
*/
type FileManager struct{
	filesStored map[string]*FileInfo
}


/** FileInfo: all the information about a file
* fileName: the file name
* lastOriginalRefreshedStored: In case of the node is the node that store this file, this time remember the last refresh
								that has to be perform each 24h
* lastTimeRefreshed: Time of the last received Store for this file
* initialStore: The time that this file has been stored the first time
* expirationTime: The living time of this file on this node
* originalStore: is this node that perform the first store for this file ?
* immutable: take or not into account the expiration time
*/
type FileInfo struct {
	fileName string
	lastOriginalRefreshedStored time.Time
	lastTimeRefreshed time.Time
	initialStore time.Time
	expirationTime float64
	originalStore bool
	immutable bool
}

//Where all the files will be stored
const filesDirectory = "kademlia/Files/"

/** checkAndStore
* PARAM: fileManager
*		 fileName: the name of the file
*		 data: the data into the file
* If the file does not exist, create it and add it the fileManager structure
*/
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
	} else{
		if fileManager.filesStored[fileName] == nil{
			file := FileInfo{fileName, time.Now().Local(), time.Now().Local(), time.Now().Local(), 24.0, false, false}
			fileManager.filesStored[fileName] = &file
		}
	}
}


/** checkAndStore
* PARAM: fileManager
*		 fileName: the name of the file
*		 data: the data into the file
* If the file does not exist, create it and add it the fileManager structure
*/
func (fileManager *FileManager) updateTime(fileName string){
	fileInfo := fileManager.filesStored[fileName]
	if fileInfo != nil{
		fileInfo.lastTimeRefreshed = time.Now().Local()
	}
}

/** checkIfFileExist
* PARAM: fileManager
*		 fileName: file to check if it exists
*
* OUTPUT: if the file exist or not
*/
func (f *FileManager) checkIfFileExist(fileName string) bool{
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

/** readData
* PARAM: fileManager
*		 fileName: the name to get the data
* OUTPUT: the data of the file in bytes
*/
func (f *FileManager) readData(fileName string) []byte{
	data, _ := ioutil.ReadFile(fileName)
	return data
}


/** RemoveFile
* PARAM: fileManager
*		 fileName: the name of the file to remove
*		 data: the data into the file
* Remove the file
*/
func (fileManager *FileManager) RemoveFile(filename string){
	delete(fileManager.filesStored, filename)
	os.Remove(filesDirectory + filename)
}

/** ListFiles
* List all the files in this node
*/
func ListFiles(){
	files, err := ioutil.ReadDir("./kademlia/Files")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

/** PinFile
* PARAM: fileManager
*		 fileName: the name of the file to pin/unpin
*		 pin: true = set to immutable, false = we can remove it
*/
func (fileManager *FileManager) PinFile(fileName string, pin bool){
	file, ok := fileManager.filesStored[fileName]
	if !ok {
		fmt.Println("This file does not exist")
	} else {
		file.immutable = pin
	}
}


func (fileInfo *FileInfo) Print(){
	fmt.Println("file name:" + fileInfo.fileName)
	fmt.Println("last time from original store refresh ", fileInfo.lastOriginalRefreshedStored)
	fmt.Println("last time received store for it ", fileInfo.lastTimeRefreshed)
	fmt.Println("initial store ", fileInfo.initialStore)
	fmt.Println("expiration time", fileInfo.expirationTime)
	fmt.Println("Initial store here ?", fileInfo.originalStore)
	fmt.Println("Is it immutabla ?", fileInfo.immutable)
}
