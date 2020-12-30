# POC for Caddy as an MQTT proxy

A small POC using Caddy as a TLS-terminating MQTT proxy

## Description

This is a small POC using Caddy as an MQTT proxy based on the (experimental) *[Project Conncept](https://github.com/mholt/caddy-l4)* TCP/UDP app.

Caddy is configured to terminate a secure MQTT port (8883) with a self-signed certificate and to proxy the TCP traffic to Mosquitto (port 1883) running in a Docker container.

In addition to MQTT, this POC also shows how MQTT over Secure WebSockets (wss) can be used.
The `wss` endpoint can be reached at `https://localhost:8443` and forwarded to Mosquitto port 9001.

Currently no matching is done, apart from matching on a valid TLS handshake, meaning that any TCP traffic (wrapped in TLS) to port 8883 or 8443 will thus be forwarded to Mosquitto.

## Enforcing MQTT Traffic

This repository also contains a custom Caddy Handler that terminates connections that don't look like MQTT.
Ideally this would have been a Matcher, but when matching an MQTT connection wrapped with TLS, it is not possible to inspect the contents without terminating the TLS connection first, which is done in the TLS Handler.
The MQTT Handler should thus be executed after the TLS termination Handler.
In `config.json` an example usage of the MQTT Handler is shown.

## TODO

* Add some type of (G)UI for visualizing MQTT communication
* Implement a handler that matches MQTT over WebSocket traffic?
* Implement an MQTTS matcher (like, TLS + MQTT wrapped; easier configuration)?
* Add configuration to MQTT matching (e.g. protocol version, authenticated yes/no, client IDs, etc ...)
* Add some publishing and subscribing clients for demo purposes
* Improve Compose setup
* Improve README.md
* Add stats/metrics (i.e. MQTT specific metrics? Bit of a stretch, though)
* Do stuff with client certificates (mTLS)
* Do something with ALPN support?