package cli

import (
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/swarm/discovery"
	ilockerjoin "github.com/yansmallb/ilocker/join"
	"strings"
)

func checkAddrFormat(addr string) bool {
	m, _ := regexp.MatchString("^[0-9a-zA-Z._-]+:[0-9]{1,5}$", addr)
	return m
}

func join(c *cli.Context) {
	dflag := getDiscovery(c)
	if dflag == "" {
		log.Fatalf("discovery required to join a cluster. See '%s join --help'.", c.App.Name)
	}

	addr := c.String("advertise")
	if addr == "" {
		log.Fatal("missing mandatory --advertise flag")
	}
	if !checkAddrFormat(addr) {
		log.Fatal("--advertise should be of the form ip:port or hostname:port")
	}

	hb, err := time.ParseDuration(c.String("heartbeat"))
	if err != nil {
		log.Fatalf("invalid --heartbeat: %v", err)
	}
	if hb < 1*time.Second {
		log.Fatal("--heartbeat should be at least one second")
	}
	ttl, err := time.ParseDuration(c.String("ttl"))
	if err != nil {
		log.Fatalf("invalid --ttl: %v", err)
	}
	if ttl <= hb {
		log.Fatal("--ttl must be strictly superior to the heartbeat value")
	}

	d, err := discovery.New(dflag, hb, ttl, getDiscoveryOpt(c))
	if err != nil {
		log.Fatal(err)
	}

	ilocker_join(c.Args()[0], addr, hb)
	for {
		log.WithFields(log.Fields{"addr": addr, "discovery": dflag}).Infof("Registering on the discovery service every %s...", hb)
		if err := d.Register(addr); err != nil {
			log.Error(err)
		}
		time.Sleep(hb)
	}
}

func ilocker_join(etcdpath string, addr string, heartbeat time.Duration) {
	etcdpath = strings.Replace(etcdpath, "etcd://", "", 1)
	endpoints := strings.Split(etcdpath, ",")
	etcdpath = ""
	for index := range endpoints {
		etcdpath += "http://" + endpoints[index]
		if index < len(endpoints)-1 {
			etcdpath += ","
		}
	}
	reg := regexp.MustCompile(":[0-9]{1,5}$")
	addr = reg.ReplaceAllString(addr, ":2374")
	go ilockerjoin.Join(etcdpath, addr, heartbeat)
}
