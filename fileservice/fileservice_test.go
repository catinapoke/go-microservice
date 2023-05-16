package fileservice

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	svc := LocalFileService{
		storage_path: "./bin",
	}

	abs, err := filepath.Abs(svc.getPath(0))
	if err == nil {
		fmt.Println("Absolute:", abs)
	} else {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	svc := CreateDefaultFileService()
	_, err := svc.Get(0)

	if err != nil {
		t.Fail()
	}
}

func createRandomFile() (string, error) {
	random_path := "./" + randStringRunes(10)
	// Create temp empty file in current directory
	err := os.WriteFile(random_path, []byte(random_path), 0666)

	return random_path, err
}

func TestSet(t *testing.T) {
	svc := CreateDefaultFileService()

	path, err := createRandomFile()

	if err != nil {
		t.Fail()
	}

	// Set this file
	id, err := svc.Set(path)

	if err != nil || id < 1 {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	svc := CreateDefaultFileService()
	path, err := createRandomFile()

	if err != nil {
		t.Fail()
	}

	// Set this file
	id, err := svc.Set(path)

	if err != nil {
		t.Fail()
	}

	svc.Delete(id)

	if err != nil {
		t.Fail()
	}
}

func launchBackgroundFileService() {
	apiServer := CreateAPIServer(CreateDefaultFileService())
	go func() {
		log.Println(apiServer.Start(":3001"))
	}()

	time.Sleep(100 * time.Millisecond)
}

func createController() FileServiceController {
	return CreateController("127.0.0.1", "3001")
}

func TestControllerSet(t *testing.T) {
	// Specify the URL to which you want to send the file
	controller := createController()
	launchBackgroundFileService()

	id, err := controller.Set("go.mod")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Set with id %d\n", id)
}

func TestControllerSetUrl(t *testing.T) {
	url := "https://go.dev/images/go-logo-white.svg"

	// Specify the URL to which you want to send the file
	controller := createController()
	launchBackgroundFileService()

	id, err := controller.Set(url)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Set with id %d\n", id)
}

func TestControllerGet(t *testing.T) {
	controller := createController()
	launchBackgroundFileService()

	path, err := controller.Get(0)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestControllerSetDelete(t *testing.T) {
	controller := createController()
	launchBackgroundFileService()

	id, err := controller.Set("go.mod")

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("TestControllerSetDelete - Set with id %d\n", id)

	err = controller.Delete(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestControllerSetGet(t *testing.T) {
	controller := createController()
	launchBackgroundFileService()

	id, err := controller.Set("go.mod")

	if err != nil {
		t.Fatal(err)
	}

	path, err := controller.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	// Compare files

	old_file, err := os.Open("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	defer old_file.Close()

	new_file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer new_file.Close()

	const chunk = 64000
	for {
		chunk_old := make([]byte, chunk)
		_, err_old := old_file.Read(chunk_old)

		chunk_new := make([]byte, chunk)
		_, err_new := new_file.Read(chunk_new)

		if err_old != nil || err_new != nil {
			if err_old == io.EOF && err_new == io.EOF {
				break
			} else if err_old == io.EOF || err_new == io.EOF {
				t.Fatal("Data is not equal by length!")
			} else {
				t.Fatal(err_old, err_new)
			}
		}

		if !bytes.Equal(chunk_old, chunk_new) {
			t.Fatal("Data is not equal!")
		}
	}
}

func TestControllerSetDeleteGet(t *testing.T) {
	controller := createController()
	launchBackgroundFileService()

	id, err := controller.Set("go.mod")

	if err != nil {
		t.Fatal(err)
	}

	err = controller.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	path, err := controller.Get(id)
	if err == nil {
		os.Remove(path)
		t.Fatal("Successfully got file, but expected error")
	}
}
