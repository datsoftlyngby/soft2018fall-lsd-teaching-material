# Assignment: SLA + Monitoring

**Deadline**: Sunday 11th of November 23:55:00

## Hand-in:

  * An URL to your Grafana dashboard
  * Your Grafana dashboard must contain:
    - At least one widget per SLA metric (KPI)
    - The SLA itself or a link to the SLA document


## SLA

You have to agree upon an SLA with the group operating your system.

The SLA *must* contain:
 
  1. A description of what the SLA includes
  2. An uptime metric
  3. A metric on "lost" posts from the simulator (customers)
  4. A metric on serving speed of the landing page
  5. At least one other metric (preferably more) that is/are relevant for
   measuring the performance of your system.

## The Dashboard

To setup a Grafana dashboard you need to:
* Install Grafana and Prometheus
  - I recommend using Docker Compose. See the example [here](https://docs.docker.com/compose/compose-file/#compose-file-structure-and-examples).
* Expose Prometheus metrics on each of your service
  - Must exist on your endpoint with the URL `/metrics`
  - Read the docs on how to do that [here](https://prometheus.io/docs/instrumenting/clientlibs/).
* Create a Grafana dashboard with a permanent URL
  - Read a guide on how to do that [here](http://docs.grafana.org/guides/getting_started/).