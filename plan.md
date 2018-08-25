# Plan

## Process

2. Determine how to satisfy data needs:
  - Stations matching a query string
  - All lines for a given station
  - All upcoming trains for a given station
  - Service changes for lines

## Architecture

- Feed ingestion (ProtoBuf stream?)
- Data storage
- HTTP server (API). Endpoints:
  - Stations matching a query string
  - Expected train arrivals at a specified station
  - Service change information for a specific train line

## Resources

- [List of all train feeds](http://datamine.mta.info/list-of-feeds)
