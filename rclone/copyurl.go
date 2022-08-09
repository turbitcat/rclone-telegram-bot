package rclone

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// json struct:
// operations/copyurl
type CopyURLRequest struct {
	Fs           string `json:"fs"`
	Remote       string `json:"remote"`
	URL          string `json:"url"`
	AutoFilename bool   `json:"autoFilename"`
}

func RetrieveFileNameFromURL(u string) string {
	filename, err := url.QueryUnescape(path.Base(u))
	if err != nil {
		filename = path.Base(u)
	}
	return filename
}

// Download a URL's content and copy it to the destination without saving it in temporary storage.
//
// 'fs' is a remote name string e.g. "drive:". 'remotePath' is a path within that remote.
// 'filename' is the saved file's name, if empty the filename will be retrieved from url.
func (rs *RcloneServer) CopyURL(u, fs, remotePath, filename string) (int, error) {
	if filename == "" {
		filename = RetrieveFileNameFromURL(u)
	}
	if strings.HasSuffix(remotePath, "/") {
		remotePath = remotePath + filename
	} else {
		remotePath = remotePath + "/" + filename
	}
	if !strings.HasSuffix(fs, ":") {
		fs = fs + ":"
	}
	request := CopyURLRequest{Fs: fs, Remote: remotePath, URL: u, AutoFilename: false}
	jobid, err := rs.DoAsync(OperationsCopyurl, request)
	if err != nil {
		return -1, fmt.Errorf("copyurl err : %v", err)
	}
	return jobid, nil
}
