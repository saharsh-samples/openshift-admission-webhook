# Sample OCP Mutating Admission Webhook

This app is intended to serve as an example of how to build custom mutating admission webhooks for OpenShift

## Build

### Build locally

The container image can be built locally using the provided [Dockerfile](build/Dockerfile).

        podman build -t saharsh-samples/maw build

## Run

### Run locally

The locally built container image can be run locally as well using Podman, Docker, or similar linux container tooling.

        podman run -d --name maw -p 8080:8080 saharsh-samples/maw

