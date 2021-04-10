package regia

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	multipartByReader = &multipart.Form{
		Value: make(map[string][]string),
		File:  make(map[string][]*multipart.FileHeader),
	}
	multipartReaderError = errors.New("http: multipart handled by MultipartReader")
)

type FileStorage interface {
	Save(filer *File, path string) error
}

type FileSystemStorage struct{}

func (f *FileSystemStorage) Save(filer *File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return filer.Copy(out)
}

type Files map[string][]*multipart.FileHeader

func (f Files) FileHeaders(key string) ([]*multipart.FileHeader, error) {
	if fhs := f[key]; len(fhs) != 0 {
		return fhs, nil
	}
	return nil, http.ErrMissingFile
}

func (f Files) Get(key string) (*File, error) {
	fhs, err := f.FileHeaders(key)
	if err != nil {
		return nil, err
	}
	return &File{fhs[0]}, nil
}

func (f Files) GetAll(key string) ([]*File, error) {
	fhs, err := f.FileHeaders(key)
	if err != nil {
		return nil, err
	}
	fs := make([]*File, len(fhs))
	for _, f := range fhs {
		fe := &File{f}
		fs = append(fs, fe)
	}
	return fs, nil
}

type File struct{ *multipart.FileHeader }

func (f *File) ContentType() (string, error) {
	file, err := f.FileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	return GetFileContentType(file)
}

func (f *File) Copy(dst io.Writer) error {
	file, err := f.FileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(dst, file)
	return err
}

// Get File ContentType
// os.File impl multipart.File
func GetFileContentType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}
