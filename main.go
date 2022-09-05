package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/turbitcat/rclone-telegram-bot/config"
	"github.com/turbitcat/rclone-telegram-bot/rclone"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func main() {
	releaseIfNotExiest("download_and_unzip.sh", download_and_unzip_sh)
	var cfg config.Config
	cfg.ReadAll()
	if !fileExists("config.yml") {
		cfg.WriteFile()
	}
	fmt.Printf("config: %+v\n", cfg)
	rs := rclone.NewRcloneServer(cfg.Rclone.BaseURL, cfg.Rclone.User, cfg.Rclone.Password)

	pref := telebot.Settings{
		Token:  cfg.TelegramBot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	remoteList, _ := rs.ListRemotes()
	var remoteFs string
	if len(remoteList) > 0 {
		remoteFs = remoteList[0]
	}
	remoteDir := "rclone-bot-downloads"

	adminOnly := b.Group()

	if cfg.TelegramBot.AdminIDs != "" {
		adminIDsStrings := strings.Split(cfg.TelegramBot.AdminIDs, ",")
		adminIDs := make([]int64, len(adminIDsStrings))
		for i, adminIDsString := range adminIDsStrings {
			n, err := fmt.Sscan(adminIDsString, &adminIDs[i])
			if err != nil {
				log.Fatal("error splite adminIDs: ", err)
				os.Exit(2)
			}
			if n != 1 {
				log.Fatal("error splite adminIDs")
				os.Exit(2)
			}
		}
		adminOnly.Use(middleware.Whitelist(adminIDs...))
	} else {
		adminOnly.Use(middleware.Whitelist())
	}

	b.Handle("/hello", func(c telebot.Context) error {
		c.Message()
		return c.Send(fmt.Sprintf("Hello, %v", c.Chat().ID))
	})

	adminOnly.Handle("/listremote", func(c telebot.Context) error {
		remoteList, err := rs.ListRemotes()
		if err != nil {
			c.Send(fmt.Sprintf("%s", err))
			return err
		}
		return c.Send(fmt.Sprintf("%s", remoteList))
	})

	adminOnly.Handle("/setdefaultfs", func(c telebot.Context) error {
		fs := c.Message().Payload
		remoteList, err := rs.ListRemotes()
		if err != nil {
			c.Send(fmt.Sprintf("%s", err))
			return err
		}
		if fs != "" && listContains(remoteList, fs) {
			remoteFs = fs
			return c.Send(fmt.Sprintf("Default fs set to '%s'.", fs))
		} else if fs == "" {
			return c.Send(fmt.Sprintf("Default fs is '%s'.", remoteFs))
		} else {
			return c.Send(fmt.Sprint(remoteList))
		}
	})

	adminOnly.Handle("/unzipanduploadurl", func(c telebot.Context) error {
		return unzipanduploadurl(rs, remoteFs, c.Message().Payload, c, b, remoteDir, cfg.TempPath.Download)
	})

	adminOnly.Handle("/uploadurl", func(c telebot.Context) error {
		return uploadurl(rs, remoteFs, c.Message().Payload, c, b, remoteDir)
	})

	newURLSelector := func(s string) (*telebot.ReplyMarkup, *telebot.Btn, *telebot.Btn) {
		var (
			selector  = &telebot.ReplyMarkup{}
			btnUpload = selector.Data("upload", "upload", s)
			btnUnzip  = selector.Data("unzip", "unzip", s)
		)
		selector.Inline(
			selector.Row(btnUpload, btnUnzip),
		)
		return selector, &btnUpload, &btnUnzip
	}

	_, btnUpload, btnUnzip := newURLSelector("")
	// btnUnzip := selector.Data("unzip", "unzip", url)

	urlCache := stringCatch{}

	adminOnly.Handle(btnUpload, func(c telebot.Context) error {
		c.Respond()
		k := c.Data()
		if k != "" {
			if url, ok := urlCache.Get(k); ok {
				return uploadurl(rs, remoteFs, url, c, b, remoteDir)
			}
		}
		return nil
	})

	adminOnly.Handle(btnUnzip, func(c telebot.Context) error {
		c.Respond()
		k := c.Data()
		if k != "" {
			if url, ok := urlCache.Get(k); ok {
				return unzipanduploadurl(rs, remoteFs, url, c, b, remoteDir, cfg.TempPath.Download)
			}
		}
		return nil
	})

	adminOnly.Handle(telebot.OnText, func(c telebot.Context) error {
		text := c.Text()
		if strings.HasPrefix(text, "http:") || strings.HasPrefix(text, "https:") || strings.HasPrefix(text, "ftp:") {
			k := urlCache.New(text)
			selector, _, _ := newURLSelector(k)
			return c.Send("It's a URL!", selector)
		}
		return nil
	})
	b.Start()

}
