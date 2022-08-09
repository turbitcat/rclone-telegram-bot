package rclone

import (
	"fmt"
)

// json struct:
// operations/copyurl
type CopyRequest struct {
	SrcFs string `json:"srcFs"`
	DstFs string `json:"dstFs"`
}

// Copy a directory from source remote to destination remote
//
// srcFs - a remote name string e.g. "drive:src" for the source
// dstFs - a remote name string e.g. "drive:dst" for the destination
func (rs *RcloneServer) Copy(src, dst string) (int, error) {
	request := CopyRequest{SrcFs: src, DstFs: dst}
	jobid, err := rs.DoAsync(SyncCopy, request)
	if err != nil {
		return -1, fmt.Errorf("copy err : %v", err)
	}
	return jobid, nil
}
