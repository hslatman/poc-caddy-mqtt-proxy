# POC for Caddy as an MQTT proxy

A small POC using Caddy as a TLS-terminating MQTT proxy

## Description

This is a small POC using Caddy as an MQTT proxy based on the (experimental) *[Project Conncept](https://github.com/mholt/caddy-l4)* TCP/UDP app.

Caddy is configured to terminate a secure MQTT port (8883) with a self-signed certificate and to proxy the TCP traffic to Mosquitto (port 1883) running in a Docker container.

Currently no matching is done, apart from matching on a valid TLS handshake, meaning that any TCP traffic (wrapped in TLS) to port 8883 will thus be forwarded to Mosquitto.

## TODO

* Add some type of (G)UI for visualizing MQTT communication
* Implement a handler that matches MQTT traffic? (i.e. inspect first bytes of the incoming connection)
* Implement a handler that matches MQTT over WebSocket traffic?
* Implement an MQTTS matcher (like, TLS + MQTT wrapped; easier configuration)
* Add configuration to MQTT matching (e.g. protocol version, authenticated yes/no, client IDs, etc ...)
* Add some publishing and subscribing clients for demo purposes
* Improve Compose setup
* Improve README.md
* Add stats/metrics (i.e. MQTT specific metrics? Bit of a stretch, though)
* Do stuff with client certificates (mTLS)