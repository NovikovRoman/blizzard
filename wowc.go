package blizzard

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/FuzzyStatic/blizzard/wowc"
)

const (
	wowPath               = "/wow"
	wowUserPath           = wowPath + "/user"
	wowUserCharactersPath = wowUserPath + "/characters"
	wowAchievementPath    = wowPath + "/achievement"
	wowAuctionDataPath    = wowPath + "/auction" + dataPath
	wowBossPath           = wowPath + "/boss"
)

// WoWUserCharacters returns all characters for users Access Token
func (c *Config) WoWUserCharacters(accessToken string) (*wowc.Profile, error) {
	var (
		dat wowc.Profile
		req *http.Request
		res *http.Response
		b   []byte
		err error
	)

	err = c.updateAccessTokenIfExp()
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest("GET", c.apiURL+wowUserCharactersPath+"?access_token="+accessToken, nil)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	res, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			return
		}
	}()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil, errors.New(res.Status)
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return nil, err
	}

	return &dat, nil
}

// WoWAchievement returns data about an individual achievement
func (c *Config) WoWAchievement(achievementID int) (*wowc.Achievement, error) {
	var (
		dat wowc.Achievement
		b   []byte
		err error
	)

	b, err = c.getURLBody(c.apiURL + wowAchievementPath + "/" + strconv.Itoa(achievementID) + "?" + localeQuery + c.locale)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return nil, err
	}

	return &dat, nil
}

// WoWAuctionFiles returns list of auction URLs containing auction data
func (c *Config) WoWAuctionFiles(realm string) (*wowc.AuctionFiles, error) {
	var (
		dat wowc.AuctionFiles
		b   []byte
		err error
	)

	b, err = c.getURLBody(c.apiURL + wowAuctionDataPath + "/" + realm + "?" + localeQuery + c.locale)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return nil, err
	}

	return &dat, nil
}

// WoWAuctionData returns auction data for realm
func (c *Config) WoWAuctionData(realm string) ([]*wowc.AuctionData, error) {
	var (
		af   *wowc.AuctionFiles
		adad []*wowc.AuctionData
		err  error
	)

	af, err = c.WoWAuctionFiles(realm)
	if err != nil {
		return nil, err
	}

	for _, file := range af.Files {
		var (
			dat wowc.AuctionData
			req *http.Request
			res *http.Response
			b   []byte
			err error
		)

		req, err = http.NewRequest("GET", file.URL, nil)
		if err != nil {
			return nil, err
		}

		res, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = res.Body.Close()
			if err != nil {
				return
			}
		}()

		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		switch res.StatusCode {
		case http.StatusNotFound:
			return nil, errors.New(res.Status)
		}

		err = json.Unmarshal(b, &dat)
		if err != nil {
			return nil, err
		}

		adad = append(adad, &dat)
	}

	return adad, nil
}

// WoWBossMasterList returns a list of all supported bosses. A "boss" in this context should be considered a boss encounter, which may include more than one NPC
func (c *Config) WoWBossMasterList() (*wowc.BossMasterList, error) {
	var (
		dat wowc.BossMasterList
		b   []byte
		err error
	)

	b, err = c.getURLBody(c.apiURL + wowBossPath + "/?" + localeQuery + c.locale)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return nil, err
	}

	return &dat, nil
}

// WoWBoss  provides information about bosses. A "boss" in this context should be considered a boss encounter, which may include more than one NPC
func (c *Config) WoWBoss(bossID int) (*wowc.Boss, error) {
	var (
		dat wowc.Boss
		b   []byte
		err error
	)

	b, err = c.getURLBody(c.apiURL + wowBossPath + "/" + strconv.Itoa(bossID) + "?" + localeQuery + c.locale)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return nil, err
	}

	return &dat, nil
}