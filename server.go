package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func serve() {
	http.HandleFunc("/favicon.ico", mw(favicon))
	http.HandleFunc("/", mw(root))
	http.HandleFunc("/config", mw(showConfig))
	http.HandleFunc("/config/new", mw(showConfigNew))
	http.HandleFunc("/example", mw(showExampleConfig))

	http.HandleFunc("/attack", mw(attack))
	http.HandleFunc("/cease_fire", mw(pause))
	http.HandleFunc("/pause", mw(pause))

	http.HandleFunc("/backoff", mw(backoff))
	http.HandleFunc("/speedup", mw(speedup))
	http.HandleFunc("/status", mw(getStatus))
	// TODO http.HandleFunc("/refresh_token", refreshToken)

	// go beginMission(cfg.Mission, reporter)
	var reporter = NewReporter()
	go reporter.listen()
	go cfg.Mission.Begin(config)

	go func() {
		for {
			select {
			case enabled := <-toggle:
				fmt.Printf("toggled state to: %v", enabled)
				cfg.Mission.EnabledCh <- enabled
			case intervalDelta := <-intervalDeltaCh:
				fmt.Printf("got an interval delta: %v", intervalDelta)
				cfg.Mission.IntervalDeltaCh <- intervalDelta
			default:
			}
		}
	}()

	tcpSocket := ":8080"
	log.WithFields(map[string]interface{}{
		"socket": tcpSocket,
		"status": "listening",
	}).Info()
	http.ListenAndServe(tcpSocket, nil)
}
