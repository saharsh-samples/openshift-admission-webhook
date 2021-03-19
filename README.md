# Sample OCP Mutating Admission Webhook

This app is intended to serve as an example of how to build custom **mutating admission webhooks** for OpenShift. Scroll below for instructions on building and deploying this app on an OpenShift cluster.

## What are admission webhooks?

Following is an excerpt from [official OpenShift documentation](https://docs.openshift.com/container-platform/4.7/architecture/admission-plug-ins.html) on admission webhooks.

> Admission plug-ins intercept requests to the master API to validate resource requests and ensure policies are adhered to, after the request is authenticated and authorized. For example, they are commonly used to enforce security policy, resource limitations or configuration requirements.
> 
> ...
>
> In addition to OpenShift Container Platform default admission plug-ins, dynamic admission can be implemented through webhook admission plug-ins that call webhook servers, in order to extend the functionality of the admission chain. Webhook servers are called over HTTP at defined endpoints.
>
> There are two types of webhook admission plug-ins in OpenShift Container Platform:

> - During the admission process, the mutating admission plug-in can perform tasks, such as injecting affinity labels.

> - At the end of the admission process, the validating admission plug-in can be used to make sure an object is configured properly, for example ensuring affinity labels are as expected. If the validation passes, OpenShift Container Platform schedules the object as configured.

## About this app

This app serves as a backend to two different mutating admission webhooks.

- [Namespace webhook](deploy/ns-webhook-config.yaml) is used to check whether a namespace is special. If that namespace is special then it is mutated to have an annotation that causes all pods in that namespaces to have a node selector of `workload-type=special`. This assumes there is at least one node on the cluster with the label `workload-type=special` where pods of special namespaces will be scheduled. It's important to note that only `CREATE` and `UPDATE` operations for namespaces that already have the `special.compliance.enabled=true` label will be intercepted by this webhook.

- [Pod webhook](deploy/pods-webhook-config.yaml) is used to intercept `CREATE` requests for pods running in namespaces with `special.compliance.enabled=true` label. If the pod is determined to belong to a special namespace, then a toleration for `workload-type=special` is added to the pod with the assumption that the node may have a corresponding taint to repel pods that are NOT special.

Together these webhooks essentially implement an example of the [Dedicated Nodes](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/#example-use-cases) use case. It is important to note both webhooks are served by the same service and deployment in OpenShift. However, each webhook has been given its own endpoint. The namespace webhook can be accessed at [`POST /admissions/namespaces`](build/src/routes/mutate-namespaces.go) and the pod webhook can be accessed at [`POST /admissions/pods`](build/src/routes/mutate-pods.go). Click on each link for its implementation.

## Building the app

### Build locally

The container image can be built locally using the provided [Dockerfile](build/Dockerfile).

        podman build -t saharsh-samples/maw build

### Build on OpenShift

The app can also be built on OpenShift using the provided [build configuration](build/build.yaml)

        oc process -p WEBHOOK_NAME=special-maw -f build/build.yaml | oc apply -f-
        
This command will deploy the build configuration and the target image stream to the active OpenShift cluster and namespace. A build can be started using the binary input from the local [`build`](build) directory using the following command.

        oc start-build special-maw --from-dir build

## Run on OpenShift

First go ahead and deploy the service and deployment objects for the app onto the OpenShift cluster in the desired namespace.

        oc process -p WEBHOOK_NAME=special-maw -f deploy/deploy.yaml | oc apply -f-

It is important to note that the service asks OpenShift to automatically generate a TLS serving cert using the `service.beta.openshift.io/serving-cert-secret-name` annotation. See [OpenShift documentation](https://docs.openshift.com/container-platform/4.7/security/certificates/service-serving-certificate.html) for details.

Once the app is fully up and running, register the webhook configurations with OpenShift. You'll need to have cluster admin privileges for these commands.

First, register the namespace webhook.

        oc process \
            -p WEBHOOK_NAME=special-maw \
            -p WEBHOOK_NAMESPACE=<namespace-where-maw-is-deployed> \
            -f deploy/ns-webhook-config.yaml \
            | oc apply -f-

Next, register the pod webhook

        oc process \
            -p WEBHOOK_NAME=special-maw \
            -p WEBHOOK_NAMESPACE=<namespace-where-maw-is-deployed> \
            -f deploy/pods-webhook-config.yaml \
            | oc apply -f-

It is, again, important to note that both webhooks use the `service.beta.openshift.io/inject-cabundle` annotation to ask OpenShift to inject the CA bundle behind the serving cert generated for the service in previous step.

Both webhooks are now live!

## Testing the webhooks

This repo includes some files to help test the webhooks as well. First create some namespaces.

        oc process -f test/namespaces.yaml | oc apply -f-

This will create three namespaces.

- **special-namespace**: This one is recognized by the webhook as a special namespace and is mutated to include the `openshift.io/node-selector` annotation.
- **ordinary-namespace**: This one is not recognized by the webhook as a special namespace and is not mutated.
- **ignored-namespace**: This namespace does not have the `special.compliance.enabled=true` label and is not intercepted by the webhook at all.

Next we can create pods in each of the namespaces.

        oc process -p APP_NAMESPACE=ignored-namespace -f test/deployment.yaml | oc apply -f-

The resulting pod from the `ignored-namespace` deployment will not be intercepted by the webhook at all.

        oc process -p APP_NAMESPACE=ordinary-namespace -f test/deployment.yaml | oc apply -f-

The resulting pod from the `ordinary-namespace` will be intercepted by the webhook, but it won't be mutated since `ordinary-namespace` is not special.

        oc process -p APP_NAMESPACE=special-namespace -f test/deployment.yaml | oc apply -f-

The resulting pod from the `special-namespace` will be intercepted by the webhook, and it will be mutated to include a toleration for `workload-type=special` taint.
