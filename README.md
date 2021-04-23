# RHPAM Kogito Operator

The RHPAM Kogito Operator is based on the [Kogito Operator](https://github.com/kiegroup/kogito-operator).

## Productization Requirements
A separate `image-prod.yaml` has been prepared for Productization that **requires CEKit 3.11+**. The `org.kie.kogito.gomoddownloader` module has been replaced with Cachito configuration. Cachito is setup in the OSBS build pod and provides the rhpam-kogito-operator repository at the configured revision along with dependencies so ensure the revision is updated before building.

The image can be built using the following command:

```bash
$ cekit --redhat --descriptor image-prod.yaml build osbs```

