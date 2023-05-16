package fileservice

import (
	"fmt"
	"os"
	"strconv"
)

type FileService interface {
	Get(id int) (string, error)
	Set(path string) (int, error)
	Delete(id int) error
}

type LocalFileService struct {
	storage_path string
}

func CreateFileService(path string) FileService {
	return &LocalFileService{path}
}

func CreateDefaultFileService() FileService {
	return CreateFileService("./bin")
}

func (svc *LocalFileService) getPath(id int) string {
	return fmt.Sprintf("%s/%d", svc.storage_path, id)
}

func (svc *LocalFileService) fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (svc *LocalFileService) getFreeId() int {
	items, _ := os.ReadDir(svc.storage_path)
	max := -1

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		number, err := strconv.Atoi(item.Name())
		if err != nil {
			continue
		}

		if number > max {
			max = number
		}
	}

	return max + 1
}

func (svc *LocalFileService) Get(id int) (string, error) {
	path := svc.getPath(id)

	if !svc.fileExists(path) {
		return "", fmt.Errorf("Get() error: %w", os.ErrNotExist)
	}

	return path, nil
}

func (svc *LocalFileService) Set(path string) (int, error) {
	id := svc.getFreeId()
	local_path := svc.getPath(id)
	err := os.Rename(path, local_path)

	if err != nil {
		return -1, fmt.Errorf("Set() error: %w", err)
	}

	return id, nil
}

func (svc *LocalFileService) Delete(id int) error {
	path := svc.getPath(id)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("Delete() error: %w", err)
	}

	return nil
}
