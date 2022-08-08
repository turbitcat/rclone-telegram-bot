package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"example.com/rclone-tgbot/rclone"
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func lastSubstring(s string, n int) string {
	return s[maxInt(0, len(s)-n):]
}

func bytesToString(b int) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case b < 0:
		return "?"
	case b < kb:
		return fmt.Sprintf("%dB", b)
	case b < mb:
		return fmt.Sprintf("%.1fKB", float32(b)/kb)
	case b < gb:
		return fmt.Sprintf("%.1fMB", float32(b)/mb)
	default:
		return fmt.Sprintf("%.1fGB", float32(b)/gb)
	}
}

func runCommendAndGetOutput(name string, arg []string, outBuffer *bytes.Buffer) error {
	cmd := exec.Command(name, arg...)
	w := io.MultiWriter(outBuffer)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func uploadurl(sc *rclone.ServerConfig, remoteFs, url string, c telebot.Context, b *telebot.Bot, remoteDir string) error {
	jobid, err := sc.CopyURL(url, remoteFs, remoteDir, "")
	if err != nil {
		c.Send(err.Error())
		return err
	}
	header := fmt.Sprintf("Uploading %.512s to %s:%s...", rclone.RetrieveFileNameFromURL(url), remoteFs, remoteDir)
	return checkJobStatusAndUpdateMessage(c.Chat(), header, nil, sc, b, jobid)
}

func unzipanduploadurl(sc *rclone.ServerConfig, remoteFs, url string, c telebot.Context, b *telebot.Bot, remoteDir, tempDownloadDir string) error {
	buffer := &bytes.Buffer{}
	tempdir, err := os.MkdirTemp(tempDownloadDir, "temp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempdir)
	msg, err := b.Send(c.Chat(), "Downloading & Unzipping...")
	if err != nil {
		return err
	}
	cmdDone := make(chan struct{})
	go func() {
		for {
			b.Edit(msg, lastSubstring(buffer.String(), 1024))
			select {
			case <-cmdDone:
				b.Edit(msg, lastSubstring(buffer.String(), 1024)+"\nDownloaded and upzipped.")
				return
			case <-time.After(5 * time.Second):
			}
		}
	}()
	if err := runCommendAndGetOutput("bash", []string{"download_and_unzip.sh", url, tempdir}, buffer); err != nil {
		cmdDone <- struct{}{}
		b.Edit(msg, err.Error())
		return err
	}
	cmdDone <- struct{}{}
	jobid, err := sc.Copy(tempdir, remoteFs+":"+remoteDir)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("Uploading %.512s(%.64s) to %s:%s...", rclone.RetrieveFileNameFromURL(url), url, remoteFs, remoteDir)
	return checkJobStatusAndUpdateMessage(c.Chat(), header, msg, sc, b, jobid)
}

func checkJobStatusAndUpdateMessage(to telebot.Recipient, initMsg string, msg *telebot.Message, sc *rclone.ServerConfig, b *telebot.Bot, jobid int) error {
	if msg == nil {
		var err error
		msg, err = b.Send(to, initMsg)
		if err != nil {
			return err
		}
	}
	for {
		status, err := sc.CheckJobStatus(jobid)
		if err != nil {
			b.Edit(msg, err.Error())
			return err
		}
		message := initMsg
		for _, st := range status.Transferring {
			message += fmt.Sprintf(
				"\n%s\n%s/%s  %d%%\nSpeed: %s/s  ETA: %.0f s",
				st.Name,
				bytesToString(st.TransferredBytes),
				bytesToString(st.ToTalBytes),
				st.Percentage,
				bytesToString(int(st.Speed)),
				st.Eta,
			)
		}
		b.Edit(msg, message)
		if status.Finished {
			if status.Success {
				b.Edit(msg, initMsg+"\nSucceed.")
				return nil
			} else {
				b.Edit(msg, initMsg+"\nFailed."+"\n"+status.Error)
				return fmt.Errorf("job failed: %v", status.Error)
			}

		}
		time.Sleep(5 * time.Second)
	}
}

func releaseIfNotExiest(name, data string) error {
	if _, err := os.Stat(name); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path.Dir(name), fs.FileMode(0660)); err != nil {
			return err
		}
		return ioutil.WriteFile(name, []byte(data), fs.FileMode(0660))

	} else {
		return err
	}
}

func containsString(a []string, s string) bool {
	for _, t := range a {
		if t == s {
			return true
		}
	}
	return false
}

type stringCatch map[string]string

func (sc stringCatch) New(s string) string {
	key := uuid.NewString()
	sc[key] = s
	return key
}

func (sc stringCatch) Get(k string) (string, bool) {
	v, ok := sc[k]
	return v, ok
}

func (sc stringCatch) Pop(k string) (string, bool) {
	v, ok := sc[k]
	delete(sc, k)
	return v, ok
}
