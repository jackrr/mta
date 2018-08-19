# Plan

## Process

1. Investigate feed data
  - Structure (how do i identify a station? a train?)
  - Reception mechanism (stream? poll? how does it work?)
2. Determine how to satisfy data needs:
  - Stations matching a query string
  - All lines for a given station
  - All upcoming trains for a given station
  - Service changes for lines
3. Determine storage mechanism
  -

## Architecture

- Feed ingestion (ProtoBuf stream?)
- Data storage
- HTTP server (API). Endpoints:
  - Stations matching a query string
  - Expected train arrivals at a specified station
  - Service change information for a specific train line

## Resources

- [List of all train feeds](http://datamine.mta.info/list-of-feeds)
