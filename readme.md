# GridWorker

**A Go library that aids in distributing tasks across grids.**

## What?

Grid worker allows you to distribute tasks across worker grids. Similar to background jobs only for jobs that require responsese.

## That's pretty cool? Is this production ready?

Probably not. This project was written for another project that will eventually be a load balancer. If you want to use this in production go right ahead, but we cannot promise it will be bug free.

## How does it work?

There are two types of machines in the GridWorker model. Requesters and Workers.

- **Requesters**: Aptly named, they request work to be done. They generally have enrmous work queues or are streaming requests in from another location (like a HTTP server). They make the requests, and get the response back from the workers.

- **Workers**: Workers actually perform the task. Workers are stateless and can be added or removed.

## Tasks

Tasks are sections of code. Requesters don't necessarily need to know how the work is being done, they just need to be able to expect the format of the output and process it from there.

Tasks should be reasonably short lived functions. A task can only be performed on a single worker. At the same time, however, you must keep in mind the effects of network latency.

For instance, it's a bad idea to have a task that performs something simple like multiplication of two numbers. Instead you're better off segmenting the list of numbers you want to multiply into larger chunks on your requester, and then sending the chunks to the workers.

When in doubt, test, benchmark, analyze.

## Examples

Examples can be found on our Wiki.

https://github.com/BlueprintProject/GridWorker/wiki/Examples
