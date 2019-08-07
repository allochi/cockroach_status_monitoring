# Cockroach DB cluster status gRPC service and CLI

This is a demo project and should not be used in production.

The goal is to create a small services that aim to monitor the health status of a cockroachdb cluster.  https://www.cockroachlabs.com/ 

If possible, with less framework as possible than gRPC and with respect of go idiomatism.

The service should offer the following GRPC endpoint(s): 

Enough information to determine the health and state of the cluster.

In addition to the deamon : 

Create a command line client (in go) allowing to retrieve the monitoring information(s). The goal of such cli tool would be to plug them within an old-shool monitoring system like nagios & co and also allow a human to use it.
