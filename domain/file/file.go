package file

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

const maxFileSize = 8 << 20 // 8MiB

var supportedFileTypeRegex = regexp.MustCompile(`.(jpg|jpeg|JPG|png|PNG|gif|GIF|ico|ICO)$`)

type File struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	FileName     string `json:"fileName"`
	FileServerID string `json:"-"`
	FilePath     string `json:"-"`
	AccessType   int8   `json:"-"` // FileAccessType
	RelatedID    int64  `json:"-"` // Depend on access type
}

func UploadPublicFile(file *multipart.FileHeader) (File, error) {
	// only support img now
	if !supportedFileTypeRegex.MatchString(file.Filename) {
		return File{}, util.NewErrorWithForbidden("unsupported file type")
	}

	record := File{
		FileName:   file.Filename,
		AccessType: int8(FileTypePublic),
	}
	err := uploadFile(file, &record)
	return record, err
}

func UploadTodoFile(userID, todoID int64, file *multipart.FileHeader) (File, error) {
	_, err := todo.OwnTodo(userID, todoID)
	if err != nil {
		return File{}, err
	}

	record := File{
		FileName:   file.Filename,
		AccessType: int8(FileTypeTodo),
		RelatedID:  todoID,
	}
	if err := uploadFile(file, &record); err != nil {
		return File{}, err
	}

	if err := TodoFileRepository.Save(todoID, record.ID); err != nil {
		return File{}, fmt.Errorf("fails to upload todo file: %w", err)
	}

	return record, nil
}

func uploadFile(file *multipart.FileHeader, record *File) error {
	write := func(err error) error {
		return fmt.Errorf("fails to upload file: %w", err)
	}

	if file.Size > maxFileSize {
		return util.NewError(errors.ErrRequestEntityTooLarge, "file too large")
	}

	if err := FileRepository.Save(record); err != nil {
		return write(err)
	}

	// TODO[pref]: avoid duplicate save, remove :id in template?
	record.FileServerID = config.Server.ID
	record.FilePath = applyFilePathTemplate(record)
	if err := util.SaveFile(file, record.FilePath); err != nil {
		return write(err)
	}

	if err := FileRepository.Save(record); err != nil {
		return write(err)
	}

	return nil
}

func GetFile(fileID int64) (*File, error) {
	file, err := FileRepository.Find(fileID)
	return file, fmt.Errorf("fails to get file: %w", err)
}

// Get file path, auto
func OwnFile(userID, fileID int64) (*File, error) {
	file, err := GetFile(fileID)
	if err != nil {
		return nil, err
	}

	switch FileAccessType(file.AccessType) {
	case FileTypePublic:
		break

	case FileTypeTodo:
		if _, err := todo.OwnTodo(userID, file.RelatedID); err != nil {
			return nil, util.NewErrorWithForbidden("unable to get non-owned file: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid file access type: %v", file.AccessType)
	}

	return file, nil
}

func applyFilePathTemplate(file *File) string {
	template := config.Server.FilePathTemplate
	template = strings.ReplaceAll(template, ":id", strconv.FormatInt(file.ID, 10))
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	template = strings.ReplaceAll(template, ":name", file.FileName)
	template = strings.ReplaceAll(template, ":path", file.FilePath)
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	return template
}

func GetFilePath(file *File) string {
	// TODO[feat]: If exist multi servers, how to get file? maybe we need redirect
	return file.FilePath
}
