# counter-machine

`counter-machine` is an in-memory database that stores a counter (positive integer)

## Features

Below are some of the features I have in mind for this DB

- [ ] Store one counter (positive integer)
- [ ] Benchmarking tool / script for users to test
    - [ ] Checked with locust or similar too, with Docker Container with resource limits
- [ ] Low memory footprint
    - [ ] Checked with Docker Container with resource limits
- [ ] Low CPU usage
    - [ ] Checked with Docker Container with resource limits
- [ ] Store multiple named counters - key-value store with key as string and value as integer
- [ ] Write Ahead Log for durability in case
- [ ] Optional Persistence based on configuration
- [ ] Replication - Allow for another instance to replicate from one database
- [ ] Replication - Allow for multiple instances to replicate from one database
- [ ] HTTP Clients for the Database in different server side languages, like Java, TypeScript (NodeJs) etc
- [ ] gRPC server with gRPC clients

## Future Ideas

- [ ] Load test and check CPU and RAM usage
    - [ ] Profile the program
    - [ ] Load test using Locust
    - [ ] Load test in resource constrained environments. Simulate such an environment with docker containers with limited resources

- [ ] Test the correctness of the implementation
    - [ ] If the counter value is correct when many (large scale) concurrent requests are sent with different values
