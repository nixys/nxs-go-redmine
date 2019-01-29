package redmine

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
	ID          string `json:"id,omitempty"`
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

	u := url.URL{
		Path: "/attachments/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.get(&a, u.String(), 200)

	return a.Attachment, status, err
}

// AttachmentUpload gets multiple issues info
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_api#Attaching-files
func (r *Context) AttachmentUpload(filePath string) (AttachmentUploadObject, int, error) {

	var a attachmentUploadResult

	u := url.URL{
		Path: "/uploads.json",
	}

	status, err := r.uploadFile(filePath, &a, u.String(), 201)
	if err != nil {
		return a.Upload, status, err
	}

	if err := attachmentGetFileProp(filePath, &a.Upload); err != nil {
		return a.Upload, status, err
	}

	return a.Upload, status, nil
}

// attachmentGetFileProp sets the file properties, such as filename and content type
func attachmentGetFileProp(filePath string, f *AttachmentUploadObject) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return err
	}

	f.ContentType = http.DetectContentType(buffer[:n])
	f.Filename = filepath.Base(filePath)

	return nil
}
