package main

import (
	"crypto/tls"
	"fmt"
	"github.com/amimof/huego"
	"goh/hoo"
	"goh/league"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	leagueClient := league.NewLocalClient()
	leagueClient.OnPlayerLevelChange(func(new int, old int, data league.ActivePlayerResponse) {

		hoo.Client.DoColorSequence([][]float64{ { 255, 0, 0 }, { 0, 255, 0 }, { 0, 0, 255 }, { 100, 100, 100 } })
		fmt.Println("Leveled up to", new)
	})
	
	leagueClient.OnPlayerDeath(func(stats league.PlayerStats) {
		hoo.Client.DoColorSequence([][]float64{ { 255, 0, 0 } })
	})

	leagueClient.OnPlayerAlive(func(stats league.PlayerStats) {
		hoo.Client.DoColorSequence([][]float64{{ 50, 50, 50 }})
	})
	hoo.Start("Stephen Room")
	something()
	leagueClient.Listen()
	for {

	}
}

func something() {
	client := http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("GET", "https://127.0.0.1:2999/liveclientdata/playerlist", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error req", err.Error())
		return
	}

	bt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bt))

}

func initUser() {
	bridge, err := huego.Discover()
	if err != nil {
		fmt.Println(err.Error())
	}
	connected := false
	var user string
	for connected == false {
		user, err = bridge.CreateUser("cyntaax")
		if err == nil {
			connected = true
		}
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("Connected!", user)
	if err != nil {
		fmt.Println("user err", err.Error())
	}
}