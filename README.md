# photobomb

Photobomb performs workflow tests against the ESP REST API.

It uses the client and structs provided by the ESP SDK package.

# TODO

- auto-refresh token once it expires
- SPA front-end using a Javascript framework
- long-running server-to-browser connections via Websockets or HTTP/2
- auto-reload config.json
- graceful shutdown by trapping signals from OS
- track (and display in UI) per-bomb and per-missile stats
- send logs to Splunk/Logstash
- configurable latency threshold(s) to pass a health check
  - any workflow can have a health check that passes only once the workflow execution limit falls under the threshold ("warm-up")
  - a secondary health check can fail when the execution limit is exceeded, optionally auto-scaling
