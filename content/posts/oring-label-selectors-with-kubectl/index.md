+++
date = "2021-04-05"
title = "ORing label selectors with kubectl"
slug = "oring-label-selectors-with-kubectl"
categories = ["technical"]
tags = ["kubernetes"]
+++

Let's say you want, using the `kubectl`, to get all `pods` with label `tier=backend` **OR** `tier=frontend`.

Did you know you could do the following?
```bash
kubectl get pods -l 'tier in (backend, frontend)'
```

We at $WORK did not.


It's actually even more powerful, you can use `notin`, you can combine with `,`
to perform AND. I recommend reading the docs for more examples.

Also, if you are at the exploration phase, don't forget to pass `--show-labels`
to have a understanding of the existing labels.

Or, if you want to reduce the clutter and only display the fields you are interested in:
```bash
kubectl get pods -l 'tier in (backend, frontend)' \
  -o custom-columns=name:.metadata.name,labels:.metadata.labels

NAME                           LABELS
admin-699855c9c6-s8fc6         map[app:admin tier:frontend]
authserver-646864ddf8-jt2f6    map[app:authserver tier:backend]
client-79678bf8c-qgp8j         map[app:client tier:frontend]
(...)
```

Here's a confession: sometimes I give up and just pipe to classic unix tools:
```bash
kubectl get pods -l 'tier in (backend, frontend)' --show-labels | \
  awk '{ print $1, $6 }' | column -t

> NAME                        LABELS
admin-699855c9c6-s8fc6        app=admin,tier=frontend
authserver-646864ddf8-jt2f6   app=authserver,tier=backend
client-79678bf8c-qgp8j        app=client,tier=frontend
(...)
```

# Links
[Labels and Selectors | Kubernetes](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#set-based-requirement)
