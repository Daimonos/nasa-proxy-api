# NASA API Proxy

## Description

A simple proxy with caching using BBolt for NASAs public API

## Running

Requires the library found at: https://github.com/jimdhughes/nasa

## Environment

`NASA_API_KEY` = DEMO_KEY or you can get a key from https://api.nasa.gov
`NASA_REDIS_URL` = URL for REDIS db (required) - default to `127.0.0.1:6379`
`NASA_PORT` = Port to listen on - default to `:80`
