package fileservice

import (
	"fmt"
	"time"
)

type FileServiceLogger struct {
	svc FileService
}

func CreateFileServiceLogger(svc FileService) FileService {
	return &FileServiceLogger{
		svc: svc,
	}
}

func (svc *FileServiceLogger) Get(id int) (path string, err error) {
	defer func(start time.Time) {
		fmt.Printf("Get(%d) started=%v took=%v path=%s err=%v\n", id, start, time.Since(start), path, err)
	}(time.Now())

	return svc.svc.Get(id)
}

func (svc *FileServiceLogger) Set(path string) (id int, err error) {
	defer func(start time.Time) {
		fmt.Printf("Set(%s) started=%v took=%v id=%d err=%v\n", path, start, time.Since(start), id, err)
	}(time.Now())

	return svc.svc.Set(path)
}

func (svc *FileServiceLogger) Delete(id int) (err error) {
	defer func(start time.Time) {
		fmt.Printf("Delete(%d) started=%v took=%v err=%v\n", id, start, time.Since(start), err)
	}(time.Now())

	return svc.svc.Delete(id)
}
