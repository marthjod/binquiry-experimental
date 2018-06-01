## separation of concerns

- split up into one-purpose minimal services
- clean package structure (domain logic-oriented)

1. frontend receives request
2. if valid, passes/routes to word type determiner
3. if valid, passes to word type parser

req -> word type? -> noun parser

*re-fetching for IDs can be done by several workers*

## 12fa

- config in the environment

## service discovery

- with k8s pods

## sre-friendly

- graceful shutdown/degradation

## grpc

- communication model
