# Low Latency framework for Game Telemetry

This repo is intended to provide a low-latency framework and helper functions for processing low-latency game telemetry data. 

Goals of the framework: 

1.  Provide a configurable interface to accept a variety of protocols & standard interfaces.

2.  Process events with no (or very minima impact) on the game server.

3.  Scale up to handle 100,000+ concurrent requests.

4.  The base framework should process events under 5ms. (added analytics & data processing could increase this time).
