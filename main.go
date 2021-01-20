package main

import (
	"fmt"
	"gato/configuration"
	"gato/github"
	icon "gato/icon"
	"github.com/getlantern/systray"
	"github.com/xeonx/timeago"
	"log"
	"os"
	"time"
)

var c configuration.Config

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Config path is needed as first argument")
	}
	arg := os.Args[1]
	var err error
	c, err = configuration.Load(arg)
	if err != nil {
		log.Fatal("Can't read configuration error:", err.Error())
	}
	onExit := func() {

	}

	systray.Run(onReady, onExit)
}

var slot = map[string][]*systray.MenuItem{}

func onReady() {
	systray.SetTemplateIcon(icon.DataBase, icon.DataBase)
	systray.SetTooltip("Lantern")

	for _, e := range c.Repositories {
		slot[e.Name] = make([]*systray.MenuItem, 3)
	}

	populate(true)
	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.DataBase, icon.DataBase)
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-ticker.C:
				{
					populate(false)
				}
			}
		}
	}()
}

func populate(first bool) {
	max := 0
	for _, r := range c.Repositories {
		actions, _ := github.GetLastAction(r.Name, r.Owner, c.GithubToken, c.GithubUsername)
		if first {
			systray.AddMenuItem(r.Name, r.Name)
		}
		for i, e := range actions {
			status := ""
			switch e.Status {
			case github.Running:
				status = "[Running...]"
			case github.Fail:
				status = "[Fail]"
			case github.Done:
				status = "[Ok]"
			case github.Queued:
				status = "[Queued]"
			case github.Cancelled:
				status = "[Cancel]"
			default:
				status = "[??]"
			}
			if max < e.Status {
				max = e.Status
			}
			t := e.Started
			title := fmt.Sprintf("  %s  %s %s", e.Branch, status, timeago.English.Format(t.Local()))
			if slot[r.Name][i] == nil {
				m := systray.AddMenuItem(title, title)
				slot[r.Name][i] = m
			} else {
				slot[r.Name][i].SetTitle(title)
				slot[r.Name][i].SetTooltip(title)
			}

		}
		if first {
			systray.AddSeparator()
		}
	}
	if max == github.Running {
		systray.SetTemplateIcon(icon.DataBlue, icon.DataBlue)
	} else if max == github.Fail {
		systray.SetTemplateIcon(icon.DataRed, icon.DataRed)
	} else {
		systray.SetTemplateIcon(icon.DataBase, icon.DataBase)
	}

}
