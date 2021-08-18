package league

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type localClient struct {
	sender http.Client
	playerCache ActivePlayerResponse
	eventsCache []GameEvent
	playerStatCache PlayerStats
	onPlayerLevelChange func(old int, new int, data ActivePlayerResponse)
	onPlayerDeath func(data PlayerStats)
	onPlayerAlive func(data PlayerStats)
}

func (lc *localClient) Listen() {
	go func() {
		for {
			playerStats, err := lc.GetPlayerStats()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			go func() {
				if playerStats.IsDead != lc.playerStatCache.IsDead {
					if playerStats.IsDead == true {
						go lc.onPlayerDeath(playerStats)
					} else {
						go lc.onPlayerAlive(playerStats)
					}
				}

				lc.playerStatCache = playerStats
			}()
		}
	}()


	go func() {
		for {
			eventsResponse, err := lc.GetGameEvents()
			if err != nil {
				fmt.Println("error getting events", err.Error())
				continue
			}
			if len(eventsResponse.Events) > len(lc.eventsCache) {
				lc.eventsCache = eventsResponse.Events
				go func() {
					for _, event := range eventsResponse.Events {
						lc.HandleEvent(event)
					}
				}()
			}
			time.Sleep(time.Millisecond * 250)
		}
	}()

	go func() {
		for {
			playerData, err := lc.GetActivePlayer()
			if err != nil {
				fmt.Println("Waiting for game to start")
				continue
			}
			if playerData.Level > lc.playerCache.Level {
				go lc.onPlayerLevelChange(playerData.Level, lc.playerCache.Level, playerData)
			}

			lc.playerCache = playerData
			time.Sleep(time.Second * 1)
		}
	}()
}

func (lc *localClient) HandleEvent(event GameEvent) {
	fmt.Println(event.EventName)
}

func (lc *localClient) GetGameEvents() (GameEventResponse, error) {
	req, _ := http.NewRequest("GET", "https://127.0.0.1:2999/liveclientdata/eventdata", nil)
	resp, err := lc.sender.Do(req)
	if err != nil {
		return GameEventResponse{}, err
	}
	defer resp.Body.Close()
	bt, _ := ioutil.ReadAll(resp.Body)
	var data GameEventResponse
	err = json.Unmarshal(bt, &data)
	if err != nil {
		return data, err
	}
	return data, err
}

func (lc *localClient) GetActivePlayer() (ActivePlayerResponse, error) {
	req, _ := http.NewRequest("GET", "https://127.0.0.1:2999/liveclientdata/activeplayer", nil)
	resp, err := lc.sender.Do(req)
	if err != nil {
		return ActivePlayerResponse{}, err
	}
	defer resp.Body.Close()
	bt, _ := ioutil.ReadAll(resp.Body)
	var data ActivePlayerResponse
	err = json.Unmarshal(bt, &data)
	return data, err
}

func (lc *localClient) GetPlayerStats() (PlayerStats, error) {
	req, _ := http.NewRequest("GET", "https://127.0.0.1:2999/liveclientdata/playerlist", nil)
	resp, err := lc.sender.Do(req)
	if err != nil {
		return PlayerStats{}, nil
	}
	defer resp.Body.Close()
	bt, _ := ioutil.ReadAll(resp.Body)
	var data []PlayerStats
	err = json.Unmarshal(bt, &data)
	if err != nil {
		return PlayerStats{}, err
	}

	for _, player := range data {
		if player.SummonerName == lc.playerCache.SummonerName {
			return player, nil
		}
	}
	return PlayerStats{}, errors.New("player not found in cache")
}

func NewLocalClient() *localClient {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	lc := localClient{
		sender: http.Client{},
		playerCache: ActivePlayerResponse{},
		eventsCache: make([]GameEvent, 0),
	}
	return &lc
}

func (lc *localClient) OnPlayerLevelChange(cb func(new int, old int, data ActivePlayerResponse)) {
	lc.onPlayerLevelChange = cb
}

func (lc *localClient) OnPlayerDeath(cb func(stats PlayerStats)) {
	lc.onPlayerDeath = cb
}

func (lc *localClient) OnPlayerAlive(cb func(stats PlayerStats)) {
	lc.onPlayerAlive = cb
}

type GameEventResponse struct {
	Events []GameEvent `json:"Events"`
}

type GameEvent struct {
	Assisters []string `json:"Assisters"`
	EventID int `json:"EventID"`
	EventName string `json:"EventName"`
	EventTime float64 `json:"EventTime"`
	KillerName string `json:"KillerName"`
	VictimName string `json:"VictimName"`
}

type ActivePlayerResponse struct {
	Abilities struct {
		E struct {
			AbilityLevel   int    `json:"abilityLevel"`
			DisplayName    string `json:"displayName"`
			ID             string `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"E"`
		Passive struct {
			DisplayName    string `json:"displayName"`
			ID             string `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"Passive"`
		Q struct {
			AbilityLevel   int    `json:"abilityLevel"`
			DisplayName    string `json:"displayName"`
			ID             string `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"Q"`
		R struct {
			AbilityLevel   int    `json:"abilityLevel"`
			DisplayName    string `json:"displayName"`
			ID             string `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"R"`
		W struct {
			AbilityLevel   int    `json:"abilityLevel"`
			DisplayName    string `json:"displayName"`
			ID             string `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"W"`
	} `json:"abilities"`
	ChampionStats struct {
		AbilityHaste                 float64 `json:"abilityHaste"`
		AbilityPower                 float64 `json:"abilityPower"`
		Armor                        float64 `json:"armor"`
		ArmorPenetrationFlat         float64 `json:"armorPenetrationFlat"`
		ArmorPenetrationPercent      float64 `json:"armorPenetrationPercent"`
		AttackDamage                 float64 `json:"attackDamage"`
		AttackRange                  float64 `json:"attackRange"`
		AttackSpeed                  float64 `json:"attackSpeed"`
		BonusArmorPenetrationPercent float64 `json:"bonusArmorPenetrationPercent"`
		BonusMagicPenetrationPercent float64 `json:"bonusMagicPenetrationPercent"`
		CritChance                   float64 `json:"critChance"`
		CritDamage                   float64 `json:"critDamage"`
		CurrentHealth                float64 `json:"currentHealth"`
		HealShieldPower              float64 `json:"healShieldPower"`
		HealthRegenRate              float64 `json:"healthRegenRate"`
		LifeSteal                    float64 `json:"lifeSteal"`
		MagicLethality               float64 `json:"magicLethality"`
		MagicPenetrationFlat         float64 `json:"magicPenetrationFlat"`
		MagicPenetrationPercent      float64 `json:"magicPenetrationPercent"`
		MagicResist                  float64 `json:"magicResist"`
		MaxHealth                    float64 `json:"maxHealth"`
		MoveSpeed                    float64 `json:"moveSpeed"`
		Omnivamp                     float64 `json:"omnivamp"`
		PhysicalLethality            float64 `json:"physicalLethality"`
		PhysicalVamp                 float64 `json:"physicalVamp"`
		ResourceMax                  float64 `json:"resourceMax"`
		ResourceRegenRate            float64 `json:"resourceRegenRate"`
		ResourceType                 string  `json:"resourceType"`
		ResourceValue                float64 `json:"resourceValue"`
		SpellVamp                    float64 `json:"spellVamp"`
		Tenacity                     float64 `json:"tenacity"`
	} `json:"championStats"`
	CurrentGold float64 `json:"currentGold"`
	FullRunes   struct {
		GeneralRunes []struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"generalRunes"`
		Keystone struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"keystone"`
		PrimaryRuneTree struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"primaryRuneTree"`
		SecondaryRuneTree struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"secondaryRuneTree"`
		StatRunes []struct {
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
		} `json:"statRunes"`
	} `json:"fullRunes"`
	Level        int    `json:"level"`
	SummonerName string `json:"summonerName"`
}

type PlayerStats struct {
	ChampionName string `json:"championName"`
	IsBot        bool   `json:"isBot"`
	IsDead       bool   `json:"isDead"`
	Items        []struct {
		CanUse         bool   `json:"canUse"`
		Consumable     bool   `json:"consumable"`
		Count          int    `json:"count"`
		DisplayName    string `json:"displayName"`
		ItemID         int    `json:"itemID"`
		Price          int    `json:"price"`
		RawDescription string `json:"rawDescription"`
		RawDisplayName string `json:"rawDisplayName"`
		Slot           int    `json:"slot"`
	} `json:"items"`
	Level           int     `json:"level"`
	Position        string  `json:"position"`
	RawChampionName string  `json:"rawChampionName"`
	RespawnTimer    float64 `json:"respawnTimer"`
	Runes           struct {
		Keystone struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"keystone"`
		PrimaryRuneTree struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"primaryRuneTree"`
		SecondaryRuneTree struct {
			DisplayName    string `json:"displayName"`
			ID             int    `json:"id"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"secondaryRuneTree"`
	} `json:"runes"`
	Scores struct {
		Assists    int     `json:"assists"`
		CreepScore int     `json:"creepScore"`
		Deaths     int     `json:"deaths"`
		Kills      int     `json:"kills"`
		WardScore  float64 `json:"wardScore"`
	} `json:"scores"`
	SkinID         int    `json:"skinID"`
	SummonerName   string `json:"summonerName"`
	SummonerSpells struct {
		SummonerSpellOne struct {
			DisplayName    string `json:"displayName"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"summonerSpellOne"`
		SummonerSpellTwo struct {
			DisplayName    string `json:"displayName"`
			RawDescription string `json:"rawDescription"`
			RawDisplayName string `json:"rawDisplayName"`
		} `json:"summonerSpellTwo"`
	} `json:"summonerSpells"`
	Team string `json:"team"`
}