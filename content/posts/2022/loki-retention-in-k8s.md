+++
date = "2022-11-29"
title = "loki retention in k8s"
tags = ["loki", "kubernetes"]
categories = ["technical"]
+++

Requirements:
* running the [loki-stack chart](https://github.com/grafana/helm-charts/tree/7d1aefde3fac41904398a00e9371ff3b1979f0bc/charts/loki-stack), tested with version 2.8.7
* filesystem persistence enabled (ie not shipping to s3, gcs etc), which is the default in the chart

To enable retention you need to have the following configs:

```yaml
loki:
  config:
    compactor:
      retention_enabled: true
    limits_config:
      retention_period: 336h
```
