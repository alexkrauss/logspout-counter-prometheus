# logspout-counter-prometheus

[Logspout](https://github.com/gliderlabs/logspout) module that can aggregate counts of certain log messages
and reports them as [Prometheus](https://prometheus.io/) metrics.

## Why?

In some setups, it is convenient to have a very lightweight mechanism that directly produces
count metrics from logs, e.g.:

* If there is no full-fledged stack for log analysis,
* if logs are not kept long enough, and we want some counts to be around for longer, 
* if everyone looks at metrics dashboards first.

## How it works

This is a very simple custom logspout module. Instead of shipping the logs elsewhere, it simply
counts occurrences of configurable text fragments. The counters are then available for scraping by
Prometheus.


## How to use it

* Use the image `alexkrauss/logspout-counter-prometheus` as a drop-in replacement for `gliderlabs/logspout`.
* Write a configuration file `counters.yaml` that looks like this:

      counters:
        some_message:
          match: "Starting handshake for type"
          level: INFO
        some_other_message:
          match: "heartbeat interval is growing too large"
          level: ERROR
    
  Here, `some_message` and `some_other_message` are arbitrary keys, and the attribute `match` determines the string
  that is matched. The attribute `level` is another arbitrary tag that is then added as a label to the metric.
  
* Mount the configuration file into the container and pass an environment variable `LOG_COUNTER_FILE` that contains the
  path to the configuration file.
  
* Start logspout with the route `counter://dummy`. Here, the host portion is not really relevant, since the logs are
  not routed anywhere. Of course you can configure another route in addition to this one.
  
* Point prometheus at the `/metrics` endpoint, where the counters are exported as the `log_count` metric.


## How to build

This repository contains the Go sources for the module, as well as a Docker build for the corresponding container,
leveraging logspout itself as a parent image. The details are explained
[here](https://github.com/gliderlabs/logspout/tree/master/custom).

## Improvement potential

This is just a first step which works well but could clearly be extended. Some improvement ideas:

* Use regex matching instead of just substring matching.
* Properly pin the logspout version that is used.
