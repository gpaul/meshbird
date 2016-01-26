package discovery

import (
	"fmt"
	"github.com/nictuku/dht"
	"log"
	"os"
	"sync"
	"time"
)

type (
	service struct {
		processStopCh chan struct{}
		waitStopCh    chan struct{}
		ih            dht.InfoHash
		node          *dht.DHT
		requestPeriod time.Duration
		peers         []string
		m             sync.Mutex
		out           chan string
	}
)

var (
	logger = log.New(os.Stderr, "[discovery] ", log.LstdFlags)
	s      *service
)

func Start(config *Config) (err error) {
	Stop()

	var ih dht.InfoHash
	if ih, err = dht.DecodeInfoHash(config.NetworkID); err != nil {
		return fmt.Errorf("decode infohash err: %s", err)
	}

	dhtConfig := dht.NewConfig()
	dhtConfig.Port = config.ListenPort

	var node *dht.DHT
	if node, err = dht.New(dhtConfig); err != nil {
		return fmt.Errorf("new dht init err: %s", err)
	}

	if err = node.Start(); err != nil {
		return fmt.Errorf("dht start err: %s", err)
	}

	s = &service{
		node:          *dht.DHT,
		ih:            ih,
		requestPeriod: config.RequestPeriod,
		out:           config.Output,
		processStopCh: make(chan struct{}),
		waitStopCh:    make(chan struct{}),
	}
	go process()
	go wait()

	return nil
}

func Stop() {
	if s != nil {
		s.node.Stop()
		close(s.processStopCh)
		close(s.waitStopCh)
		s = nil
	}
}

func process() {
	t := time.NewTicker(s.requestPeriod)
	defer t.Stop()

	sendRequest()

	for {
		select {
		case <-t.C:
			sendRequest()
		case <-s.processStopCh:
			break
		}
	}
}

func sendRequest() {
	logger.Printf("sending request...")
	s.node.PeersRequest(string(s.ih), true)
}

func wait() {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case r := <-s.node.PeersRequestResults:
			for _, peers := range r {
				for _, x := range peers {
					addPeer(dht.DecodePeerAddress(x))
				}
			}
		case <-s.waitStopCh:
			break
		}
	}
}

func addPeer(host string) {
	s.m.Lock()
	post := !strInSlice(host, s.peers)
	if post {
		if len(s.peers) > 1000 {
			s.peers = s.peers[1:]
		}
		s.peers = append(s.peers, host)
	}
	s.m.Unlock()

	if post {
		logger.Printf("new peer: %s", host)
		s.out <- host
	}
}

func strInSlice(s string, sl []string) bool {
	for _, i := range sl {
		if i == s {
			return true
		}
	}
	return false
}
