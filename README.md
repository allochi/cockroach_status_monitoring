# Cockroach DB cluster status gRPC service and CLI

**This is a demo project and should not be used in production.**

The goal is to create a small services that aim to monitor the health status of a cockroachdb cluster.  https://www.cockroachlabs.com/ 

If possible, with less framework as possible than gRPC and with respect of go idiomatism.

The service should offer the following GRPC endpoint(s): 

Enough information to determine the health and state of the cluster.

In addition to the deamon : 

Create a command line client (in go) allowing to retrieve the monitoring information(s). The goal of such cli tool would be to plug them within an old-shool monitoring system like nagios & co and also allow a human to use it.

## Instructions

- You can compile the whole solution using only `make`
- The `dist` folder will hold the server, the client and the cli tool
- The cli too doesn't require the server and it serializes all nodes state into `json` and print them out
- Both the server and the cli require an initial address and then it's up to them to discover the cluster
- The server periodically calls the cluster for it's state and cache that state, the duration between updates can be modified
- If the entry host is not specified it's assumed to be `localhost:8008`
- If the entry host dies, the server automatically tries other addresses from its internal cache and promote the first successful address to be the new entry address


## Notes

- I started with the idea of using the CLI of `cockroachdb`, unfortunately it doesn't report on node memory
- I modified the code to call `_status/vars`, but then, `cockroachdb node status` doesn't provide the http address so I had to replace it all with http request
- The server and cli assume that they can access the nodes `_status_` endpoints
