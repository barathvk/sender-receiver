# Sender-Receiver with Fallback

## Introduction
For this exercise, I have chosen to use `golang` as my preferred language.

The purpose of this exercise was to create a simple stateful distributed application with automated fallback.
I have chosen to go about this task in 4 phases:

  1. An initial naive implementation of a sender and receiver with no fallback logic
  1. An in process monitor to trigger creation of a fallback process
  1. Offload the state management to an external system to reduce strong coupling between the sender and receiver
  1. Offload the fallback process creation to an external system to remove strong coupling between the sender and receiver

The end result is a loosely coupled stateful distributed system with automated fallback.

## Part 1 - Sender and receiver with no fallback
[v1.0.0](https://github.com/barathvk/sender-receiver/tree/v1.0.0) contains the code for this section.

This phase of the solution contains two primary portions
  * A receiver process that implements a RESTful API that accepts a payload as defined in [common/model.go](common/model.go) at the endpoint `/count`, returns a `204: Accepted` response, and logs the message out to the console. The receiver process accepts an `appId` argument and a `port` argument for the API server to listen on.
  * A sender process that makes HTTP POST requests to the receiver REST API with the agreed payload every second. The sender process accepts an `appId` as an initialization argument and an `initialValue` argument. A `nodeId` is generated to identify the sender using github.com/segmentio/ksuid

### Testing the implementation
To test this implementation:

  1. Clone this repository
  1. `git checkout v1.0.0`
  1. Start the receiver `go run .`
  1. Start the sender `go run . --sender`

## Part 2.1 - Initial naive failover
[v2.0.0](https://github.com/barathvk/sender-receiver/tree/v2.0.0) contains the code for this section

This phase builds on [v1.0.0](https://github.com/barathvk/sender-receiver/tree/v1.0.0) by designating the receiver as a state manager and controller for the sender process with a naive in-process failover.

The receiver process implements a [heartbeat](https://github.com/barathvk/sender-receiver/blob/82e094b9dc31ebc57884bcce38b518304bc03d38/receiver/receiver.go#L29) `goroutine` that keeps track of the [lastRequest](https://github.com/barathvk/sender-receiver/blob/82e094b9dc31ebc57884bcce38b518304bc03d38/receiver/receiver.go#L15) time that the receiver has received a request from the sender. If a heartbeat is not detected within 1 second, the sender process is assumed to have died and a [new sender process](https://github.com/barathvk/sender-receiver/blob/82e094b9dc31ebc57884bcce38b518304bc03d38/receiver/receiver.go#L19) is started by the receiver.

While this failover is automated and works, the sender and receiver process are strongly coupled since the receiver manages state and failover.

### Testing the implementation

  1. Clone this repository
  1. `git checkout v2.0.0`
  1. `go build .`
  1. Start the receiver `go run .`
  1. In a separate terminal, run `curl -X POST http://localhost:8080/stop`.
  1. You should notice the original sender dies and a new sender process is started with the last count

## Part 2.2 - Offload state management

[v3.0.0](https://github.com/barathvk/sender-receiver/tree/v3.0.0) contains the code for this section

This phase replaces the internal state management in [v2.0.0](https://github.com/barathvk/sender-receiver/tree/v2.0.0) with an external system. For the purposes of this exercise, I will be using Redis and its [keyspace notifications](https://redis.io/topics/notifications) feature to act as the state management system. This simplifies the implementation to a large extent as we do not need to maintain the current state of the system in our code. Essentially, the sender connects to Redis and [sends a message](https://github.com/barathvk/sender-receiver/blob/40d0dcaea085850fbb2e79fbe5241f42e8e11cf7/sender/sender.go#L13) every second. The heartbeat is still maintained by the receiver and the sender process is still controlled by the receiver. The receiver subscribes to key changes from Redis and logs any messages received.

This phase allows us to remove the `/count` from the REST API, but the `/stop` endpoint still remains to allow us to control the sender process.

In a future implementation, we could replace Redis with a more robust durable message system like Apackhe Kafka.

### Testing the implementation

  1. Clone this repository
  1. `git checkout v3.0.0`
  1. `go build .`
  1. `docker-compose up`
  1. Start the receiver `go run .`
  1. In a separate terminal, run `curl -X POST http://localhost:8080/stop`.
  1. You should notice the original sender dies and a new sender process is started with the last count

## Part 2.3 - Offload automated fallback

[v4.0.0](https://github.com/barathvk/sender-receiver/tree/v4.0.0) contains the code for this section

This phase simplifies the implementation and removes all strong coupling between the sender and the receiver, however, the complexity in terms of infrastructure increases.

Kubernetes specializes in maintaining required state. It makes a lot of sense to offload this functionality to Kubernetes. When a sender pod is deleted, a new one is created by Kubernetes since the `replicas` spec of the deployment is [set to `1`](https://github.com/barathvk/sender-receiver/blob/edf5f2d18226908b6fce041458eddd4ecada5c55/deploy/sender.tf#L7).
This phase removes the need for a REST API.

### Infrastructure
The code required to deploy the infrastructure to GCP can be found [here](https://github.com/barathvk/sender-receiver/tree/v4.0.0/infrastructure). The infrastructure is fairly bare bones with:

  1. A simple single node Kubernetes cluster on GKE
  1. Some secuirty settings to enforce TLS on my preferred DNS provider (Cloudflare)
  1. A cluster wide deployment of Redis
  1. Nginx ingress for any ingress needs

### Testing the infrastructure

  1. Clone this repository
  1. `git checkout v4.0.0`
  1. `cd infrastructure && terraform init && terraform apply`
  1. `gcloud container clusters get-credentials sender-receiver --zone europe-west3-a`
  1. Use the `.github/workflows` folder to deploy to the kubernetes cluster via Github Actions
  1. Once deployed, watch the logs of the receiver with `kubectl logs -f --namespace sender-receiver master-receiver-<random_id>`
  1. Delete the running sender pod with `kubectl delete pod --namespace sender-receiver master-sender-<random_id>`
  1. Note that a new sender pod is created with a new `nodeId` and starts sending count messages to the receiver from where the previous sender left off.

## Continuous integration

For this exercise, I have chosen to use Github Actions as the CI system. The CI files can be found under [.github/workflows](https://github.com/barathvk/sender-receiver/tree/v4.0.0/.github/workflows). I have also chosen to work with [Trunk based development](https://trunkbaseddevelopment.com/).

  * Each feature is developed on a dedicated branch
  * Each branch is deployed for testing new features in isolation
  * When the branch is merged to `master` via a pull request, the review deployment is deleted and the new `master` is deployed to the cluster as the `staging` deployment
  * When a new tag is created, it is deployed as the latest `production` deployment

The CI system performs 3 tasks:

  1. Build a Docker image of the application and publish it to Dockerhub
  1. Deploy the image with 3 deployments to the kubernetes cluster
      1. sender (`./sender-receiver --sender`)
      2. receiver (`./sender-receiver`)
      3. redis
  1. Once a branch is deleted, the review deployment is deleted

The deployment `terraform` files can be reviewed [here](https://github.com/barathvk/sender-receiver/tree/v4.0.0/deploy)

## Additional resources

The CI runs and logs can be seen [here](https://github.com/barathvk/sender-receiver/actions)
The pull requests for each feature can be seen [here](https://github.com/barathvk/sender-receiver/pulls?q=is%3Apr+is%3Aclosed)