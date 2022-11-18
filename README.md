# Software Development Kit for indexers
SDK for creation indexers by DipDup. It's a set of package which can be used for building indexers.

## Messages

It's package for inner-communication between indexer's components. Package implements PubSub pattern via channels. Detailed docs can be found [here](/pkg/messages/).

## Storage

Abstract layer of data storage is described in the package. Detailed docs can be found [here](/pkg/storage/).

## Modules

The workflow is builded by modules.

### gRPC

gRPC module where realized default client and server. Detailed docs can be found [here](/pkg/modules/grpc/).

### Cron

Cron module implements cron scheduler. Detailed docs can be found [here](/pkg/modules/cron/).