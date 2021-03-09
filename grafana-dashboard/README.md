# Open Policy Agent Grafana dashboard

![](images/dashboard-image-1.png)

This folder hosts an example on the Grafana dashboard for Open Policy Agent.

It visualizes most of the metrics recorded by Open Policy Agent as documented in 
https://www.openpolicyagent.org/docs/latest/monitoring/

The version of OPA that this is made for is 0.26.

**Notes:** Before you use this dashboard, please mind the variables and the prometheus labels.
It may not be applicable in your case or your K8s cluster may use another name than mine.
For example: the staging/production datasource, the namespace of the system, etc.

## Remaining issues
- Avg response time and http response time doesnt seem to be correct. The units of avg response time is really weird. Need more investigation in OPA source code as well as my queries.
This should be fixed in OPA versions later than 0.26.0 thanks to this [pull request](https://github.com/open-policy-agent/opa/pull/3214)


This dashboard is published on Grafana at https://grafana.com/grafana/dashboards/13965. 


