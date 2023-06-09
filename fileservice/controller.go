package fileservice

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type FileServiceController struct {
	url string
}

func CreateController(ip, port string) FileServiceController {
	return CreateControllerByUrl(fmt.Sprintf("http://%s:%s", ip, port))
}

func CreateControllerByUrl(url string) FileServiceController {
	return FileServiceController{
		url: url,
	}
}

func (c *FileServiceController) Get(id int) (string, error) {
	request, err := http.NewRequest("GET", c.url+"/get?id="+strconv.Itoa(id), nil)
	if err != nil {
		return "", err
	}

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("got http error: " + response.Status)
	}

	// Copy the response body
	filename := "./temp/" + randStringRunes(10) + "_" + strconv.Itoa(id)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (c *FileServiceController) Set(path string) (int, error) {
	// Open the file you want to send

	var body *bytes.Buffer = nil
	var writer *multipart.Writer
	var err error

	if path[:4] != "http" {
		body, writer, err = c.createMultipartForm(path)

		if err != nil {
			return -1, fmt.Errorf("FileSerivceController Set() error: %w", err)
		}
	}

	// Create a new POST request with the multipart form as the body
	var request *http.Request
	if body != nil {
		request, err = http.NewRequest("POST", c.url+"/set", body)
		if err != nil {
			return -1, err
		}
		request.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		request, err = http.NewRequest("POST", c.url+"/set?url="+path, nil)
		if err != nil {
			return -1, err
		}
		request.Header.Set("Content-Type", "multipart/form-data")
	}

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return -1, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return -1, errors.New("got http error: " + response.Status)
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	id, err := strconv.Atoi(string(responseBody))
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (c *FileServiceController) createMultipartForm(path string) (*bytes.Buffer, *multipart.Writer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a part for the file and add it to the form
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	// Close the multipart form
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}

	return body, writer, nil
}

func (c *FileServiceController) createEmptyMultipartForm(path string) (*bytes.Buffer, *multipart.Writer, error) {
	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	return body, writer, nil
}

func (c *FileServiceController) Delete(id int) error {
	request, err := http.NewRequest("GET", c.url+"/delete?id="+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("got http error " + response.Status)
	}

	return nil
}
