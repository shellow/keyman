package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/shellow/keyman"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
)

var Logger *zap.SugaredLogger

var HOSTURL string
var KEY string

func main() {
	app := cli.NewApp()
	app.Name = "Key manage"
	app.Usage = "Key manage"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		//cli.IntFlag{
		//	Name:  "port, p",
		//	Value: 8000,
		//	Usage: "listening port",
		//},
		cli.StringFlag{
			Name:        "surl, s",
			Value:       "http://127.0.0.1",
			Usage:       "server url",
			Destination: &HOSTURL,
		},
		cli.StringFlag{
			Name:        "key, k",
			Value:       "key",
			Usage:       "server key",
			Destination: &KEY,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:     "list",
			Usage:    "list keys",
			Category: "manage",
			Action:   listkey,
		},
		{
			Name:     "add",
			Usage:    "add key",
			Category: "manage",
			Action:   addkey,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hkey, hk",
					Value: "1",
					Usage: "key for add",
				},
				cli.StringFlag{
					Name:  "hkeyname, kn",
					Value: "key",
					Usage: "key name",
				},
			},
		},
		{
			Name:     "enable",
			Usage:    "enable key",
			Category: "manage",
			Action:   enablekey,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hkey, hk",
					Value: "1",
					Usage: "key for add",
				},
				cli.StringFlag{
					Name:  "day",
					Value: "10",
					Usage: "Time limit",
				},
				cli.StringFlag{
					Name:  "num",
					Value: "10",
					Usage: "Limit of times",
				},
			},
		},
		{
			Name:     "get",
			Usage:    "get key",
			Category: "manage",
			Action:   getkey,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hkey, hk",
					Value: "1",
					Usage: "key for get",
				},
			},
		},
		{
			Name:     "dis",
			Usage:    "dis key",
			Category: "manage",
			Action:   diskey,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hkey, hk",
					Value: "1",
					Usage: "key for dis",
				},
			},
		},
		{
			Name:     "del",
			Usage:    "del key",
			Category: "manage",
			Action:   delkey,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hkey, hk",
					Value: "1",
					Usage: "key for del",
				},
			},
		},
		{
			Name:     "token",
			Usage:    "get token",
			Category: "manage",
			Action:   gettoken,
		},
		{
			Name:     "address",
			Usage:    "get key address",
			Category: "manage",
			Action:   keyaddr,
		},
		{
			Name:     "getownkey",
			Usage:    "get own key",
			Category: "manage",
			Action:   getownkey,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func listkey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/listkey"
	req, err := http.NewRequest("GET", murl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func addkey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/addkey"
	var b bytes.Buffer
	var key keyman.HKey
	key.Key = c.String("hkey")
	key.Name = c.String("hkeyname")
	bj, err := json.Marshal(key)
	if err != nil {
		return err
	}
	b.Write(bj)
	req, err := http.NewRequest("POST", murl, &b)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func enablekey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/enable"
	var b bytes.Buffer
	var key keyman.Key
	key.Key = c.String("hkey")
	key.Expday = c.Int("day")
	key.Number = c.Int64("num")
	bj, err := json.Marshal(key)
	if err != nil {
		return err
	}
	b.Write(bj)
	req, err := http.NewRequest("POST", murl, &b)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func getkey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/getkey"
	var b bytes.Buffer
	var key keyman.HKey
	key.Key = c.String("hkey")
	bj, err := json.Marshal(key)
	if err != nil {
		return err
	}
	b.Write(bj)
	req, err := http.NewRequest("POST", murl, &b)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func diskey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/diskey"
	var b bytes.Buffer
	var key keyman.HKey
	key.Key = c.String("hkey")
	bj, err := json.Marshal(key)
	if err != nil {
		return err
	}
	b.Write(bj)
	req, err := http.NewRequest("POST", murl, &b)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func delkey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/delkey"
	var b bytes.Buffer
	var key keyman.HKey
	key.Key = c.String("hkey")
	bj, err := json.Marshal(key)
	if err != nil {
		return err
	}
	b.Write(bj)
	req, err := http.NewRequest("POST", murl, &b)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func gettoken(c *cli.Context) error {
	murl := c.GlobalString("surl")
	req, err := http.NewRequest("PUT", murl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	fmt.Println(res.Header.Get("token"))
	return nil
}

func keyaddr(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/keyaddr"
	req, err := http.NewRequest("GET", murl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func getownkey(c *cli.Context) error {
	murl := c.GlobalString("surl")
	murl = murl + "/keymem/getownkey"
	req, err := http.NewRequest("GET", murl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("key", KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
