+++ 
date = 2022-12-31T15:00:40Z
title = "Little's law"
tags = ["queues"]
categories = ["technical"]
+++

I've came across this interesting theorem by John Little about queues.
Thing with queues is that everything is a queue! Queues in a supermarket, yes.
But also requests being processed in the server, or obviously enough, a worker queue too.

The theorem is as follows:

```python
L = λ * W
```

Where
`L` = average number of items in the queue
`λ` = average rate of new items arriving in the queue
`W` = average time spent in the queue

Applying to our world, let's say we have 1000 requests per second, and we spend 0.1 second per request.

```python
L = λ * W

L = 1000 * 0.1
```

Therefore, the average number of items in the queue is 100.
Which we can then based on the item size estimate how much memory it will be used on average.

By the way, we can also shift the relation. For example, if we want to figure out
how many items can enter the queue at a time, we just need to calculate L (capacity)
divided by W (latency):

```
λ = L/W
```

Which makes it clear that to support handling more items, we need to increase the capacity
OR reduce the latency.

The beauty of this theorem is that it's simple enough that you can easily extrapolate to other questions.
For example, let's say you are doing batch processing, where L (capacity) and W (latency) is fixed.
To not overload the system you can simply... add a bigger delay between adding items to the queue,
therefore reducing `λ`.


# Sources

* [Back of the envelope estimation hacks | Roberto Vitillo's Blog](https://robertovitillo.com/back-of-the-envelope-estimation-hacks/)
* [Little’s Law, Scalability and Fault Tolerance: The OS is your bottleneck. What you can do? - High Scalability -](http://highscalability.com/blog/2014/2/5/littles-law-scalability-and-fault-tolerance-the-os-is-your-b.html)
