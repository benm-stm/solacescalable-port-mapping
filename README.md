![Build and Test](https://github.com/benm-stm/solacescalable-port-mapping/actions/workflows/buildAndTest.yml/badge.svg)

# solacescalable-port-mapping
organized display of port mapping for the solacescalable operator
Example display
```
$ kubectl solmap -n solacescalable -c solacescalable
+--------------------------+--------------+-------------+--------------+
|     SERVICE NAME PUB     | SERVICE PORT | SOLACE PORT | HAPROXY PORT |
+--------------------------+--------------+-------------+--------------+
| test-botti-1028-amqp-pub |         1028 |        1100 |        32058 |
| test-botti-1029-amqp-pub |         1029 |        1100 |        32562 |
| test-default-1030-na-pub |         1030 |        1100 |        30763 |
| test-default-1031-na-pub |         1031 |        1050 |        30576 |
| test-botti-1025-mqtt-pub |         1025 |        1050 |        32263 |
| test-botti-1026-mqtt-pub |         1026 |        1050 |        31552 |
| test-botti-1027-amqp-pub |         1027 |        1100 |        32037 |
+--------------------------+--------------+-------------+--------------+
+--------------------------+--------------+-------------+--------------+
|     SERVICE NAME SUB     | SERVICE PORT | SOLACE PORT | HAPROXY PORT |
+--------------------------+--------------+-------------+--------------+
| test-botti-1025-amqp-sub |         1025 |        1100 |        32582 |
| test-default-1026-na-sub |         1026 |        1100 |        30527 |
| test-default-1027-na-sub |         1027 |        1050 |        30640 |
+--------------------------+--------------+-------------+--------------+
```

# Prerequisits
- Go 1.18
- Kubectl
# How to integrate
Build the project to generate the plugin's binary

**NOTE**: Prefix your plugin with kubectl- so it gets recognized by kubectl

```
$ cd solacescalable-port-mapping
$ go build -o kubectl-solmap
```

search for your kubectl plugins list to get the plugins dir
```
$ kubectl plugin list
kubectl plugin list

The following compatible plugins are available:

/opt/google-cloud-sdk/bin/kubectl-ports
```

Copy the generated binary to your plugins location
```
$ cp kubectl-solmap /opt/google-cloud-sdk/bin/
```
