package redmine

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nixys/nxs-go-redmine/mimereader"
)

/* Get */

// AttachmentObject struct used for attachments get operations
type AttachmentObject struct {
	ID          int    `json:"id"`
	FileName    string `json:"filename"`
	FileSize    string `json:"filesize"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
	ContentURL  string `json:"content_url"`
	Author      IDName `json:"author"`
	CreatedOn   string `json:"created_on"`
}

/* Upload */

// AttachmentUploadObject struct used for attachments upload operations
type AttachmentUploadObject struct {
	ID          int    `json:"id,omitempty"`
	Token       string `json:"token"`
	Filename    string `json:"filename"`     // This field fills in AttachmentUpload() function, not by Redmine. User can redefine this value manually
	ContentType string `json:"content_type"` // This field fills in AttachmentUpload() function, not by Redmine. User can redefine this value manually
}

/* Internal types */

type attachmentSingleResult struct {
	Attachment AttachmentObject `json:"attachment"`
}

type attachmentUploadResult struct {
	Upload AttachmentUploadObject `json:"upload"`
}

// AttachmentSingleGet gets single attachment info
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Attachments#GET
func (r *Context) AttachmentSingleGet(id int) (AttachmentObject, int, error) {

	var a attachmentSingleResult

	ur := url.URL{
		Path: "/attachments/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.get(&a, ur, http.StatusOK)

	return a.Attachment, status, err
}

// AttachmentUpload uploads file
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_api#Attaching-files
func (r *Context) AttachmentUpload(filePath string) (AttachmentUploadObject, int, error) {

	var a attachmentUploadResult

	ur := url.URL{
		Path: "/uploads.json",
	}

	f, err := os.Open(filePath)
	if err != nil {
		return a.Upload, 0, err
	}
	defer f.Close()

	mr := mimereader.New(f)

	status, err := r.uploadFile(mr, &a, ur, http.StatusCreated)
	if err != nil {
		return a.Upload, status, err
	}

	a.Upload.ContentType = mr.DetectContentType()
	a.Upload.Filename = filepath.Base(filePath)

	return a.Upload, status, nil
}

// AttachmentUploadStream uploads file as a stream.
func (r *Context) AttachmentUploadStream(f io.Reader, fileName string) (AttachmentUploadObject, int, error) {

	var a attachmentUploadResult

	ur := url.URL{
		Path: "/uploads.json",
	}

	mr := mimereader.New(f)

	status, err := r.uploadFile(mr, &a, ur, http.StatusCreated)
	if err != nil {
		return a.Upload, status, err
	}

	a.Upload.ContentType = mr.DetectContentType()
	a.Upload.Filename = filepath.Base(fileName)

	return a.Upload, status, nil
}
