package hoo

import (
	"fmt"
	"github.com/amimof/huego"
	"github.com/lucasb-eyer/go-colorful"
	"time"
)

type hoo struct {
	Bridge *huego.Bridge
	Group *huego.Group
}

func (h *hoo) DoColorSequence(sequence [][]float64) {
	for _, colors := range sequence {
		x, y , _ := colorful.LinearRgb(colors[0], colors[1], colors[2]).Xyy()
		h.Group.Xy([]float32{float32(x), float32(y)})
		time.Sleep(time.Millisecond * 200)
	}
}

var Client = hoo{}

func Start(groupName string) {
	bridge := huego.New("", "")
	Client.Bridge = bridge
	groups, err := bridge.GetGroups()
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	for _, group := range groups {
		fmt.Println(group.Name)
		if group.Name == groupName {
			group.DisableStreaming()
			Client.Group = &group
			break
		}
	}
	fmt.Println("Hoo started")
}