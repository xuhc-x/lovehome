package models

import (
	"fmt"
	"github.com/tedcy/fdfs_client"
)

func UpLoadFile(filename string) {
	client,err:=fdfs_client.NewClientWithConfig("conf/client.conf")
	defer client.Destory()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Print("111000")
		return
	}
	fileId, errUpload := client.UploadByFilename(filename)
	if errUpload != nil {
		fmt.Println(errUpload.Error())
		fmt.Print("000111")
		return
	}
	fmt.Println(fileId)

/*	if err := client.DownloadToFile(fileId, "tempFile", 0, 0); err != nil {
		fmt.Println(err.Error())
		return
	}
	if buffer, err := client.DownloadToBuffer(fileId, 0, 19); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(buffer))
	}
	if err := client.DeleteFile(fileId); err != nil {
		fmt.Println(err.Error())
		return
	}
*/
}