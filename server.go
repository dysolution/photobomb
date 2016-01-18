package main

import (
	"net/http"
	"time"

	"github.com/dysolution/espsdk"
)

func httpd() {
	http.HandleFunc("/", mw(status))
	http.HandleFunc("/config", mw(showConfig))
	http.HandleFunc("/example", mw(showExampleConfig))

	http.HandleFunc("/attack", mw(attack))
	http.HandleFunc("/cease_fire", mw(pause))
	http.HandleFunc("/once", mw(once))
	http.HandleFunc("/pause", mw(pause))
	http.HandleFunc("/warning_shot", mw(once))

	http.HandleFunc("/faster", mw(faster))
	http.HandleFunc("/slower", mw(slower))
	http.HandleFunc("/backoff", mw(backoff))
	http.HandleFunc("/speedup", mw(speedup))
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	go func() {
		log.Infof("initial interval: %v seconds", interval)
		for {
			select {
			case enabled = <-toggle:
			case d := <-intervalDelta:
				setInterval(d)
			default:
			}
			if enabled {
				log.Debugf("conducting raid")
				config.Conduct(log, espsdk.APIInvariant, warningThreshold)
				raidCount++
				log.Debugf("sleeping for %v seconds", interval)
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()

	tcpSocket := ":8080"
	log.Infof("listening on %s", tcpSocket)
	http.ListenAndServe(tcpSocket, nil)
}
