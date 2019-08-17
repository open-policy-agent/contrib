# Managing IPTable Rules with OPA

IPTables is a useful tool available to Linux kernel for filtering network packets. **opa-iptables** extension provides the management of IPTables rules with Rego policy. The purpose of this extension is to manage rules using OPA. Here OPA is used as a centralized location for storing rules and write a context-aware policy to insert/delete rules to Linux host.

## Getting Started

- If you want to build opa-iptables right away, you need a working [Go environment](https://golang.org/doc/install). It requires Go version 1.12 and above.
```
$ git clone https://github.com/open-policy-agent/contrib.git
$ cd opa-iptables
$ make go-build
```

- You can build opa-iptables on any system that has Docker running as follows:
```
$ git clone https://github.com/open-policy-agent/contrib.git
$ cd opa-iptables
$ make docker-build
```

## Tutorial

- For more information on how to use this extension with OPA to manage iptables rules, checkout [this tutorial document](./docs/tutorial.md).

## How To Run?

opa-iptables contain a web server that runs on port `33455`. iptables is a Linux kernel facility that is used for managing network traffic of Network Layer(L3/L4). This is a reason, it's only working with Linux kernel. The following is a list of all the command-line flags for configuration.

> Note: Because of this extension need to access iptables Linux utility, It requires root privileges. Therefore you need to run it using `sudo`.

```
sudo ./opa-iptables -h

Usage of ./opa-iptables:
 -controller-host string
  controller host (default "0.0.0.0")
 -controller-port string
  controller port on which it listen on (default "33455")
 -log-format string
  set log format. i.e. text | json | json-pretty (default "text")
 -log-level string
  set log level. i.e. info | debug | error (default "info")
 -opa-endpoint string
  endpoint of opa in form of http://ip:port i.e. http://192.33.0.1:8181 (default "http://127.0.0.1:8181")
 -v  show version information
```

**Run As Docker Container:**

```
docker run --rm --net host --cap-add=NET_ADMIN urvil38/opa-iptables:0.0.1-dev -log-level debug -opa-endpoint http://127.0.0.1:8181
```

## API

## **Insert Rule**

```
POST /v1/iptables/insert
Content-Type: application/json
```
```
{
 "input": ...
}
```

Insert rules into the kernel.

The request body contains an object that specifies a value for [The input Document](https://www.openpolicyagent.org/docs/latest/how-does-opa-work#the-input-document).

#### Query Parameters

- **q** - path to OPA policy's rule

#### Status Code

- **200 OK** - Successfully inserted given iptables rules

- **400 Bad Request** - If provided query path didn't resolve to any defined OPA policy rule or server fails to parse JSON payload

- **404 Not Found** - OPA policy didn't return any iptables rules

- **500 Server Error** - Fail to insert given iptables rules

## **Delete Rule**

```
POST /v1/iptables/delete
Content-Type: application/json
```
```
{
 "input": ...
}
```

Delete rules from the kernel.

The request body contains an object that specifies a value for [The input Document](https://www.openpolicyagent.org/docs/latest/how-does-opa-work#the-input-document).

#### **Query Parameters**

- **q** - path to OPA policy's rule

#### **Status Code:**

- **200 OK** - Successfully deleted given iptables rules

- **400 Bad Request** - If provided query path didn't resolve to any defined OPA policy rule or server fails to parse JSON payload

- **404 Not Found** - OPA policy didn't return any iptables rules

- **500 Server Error** - Fail to delete given iptables rules

## **List Rules**

```
GET /v1/iptables/list/{table}/{chain}
```

List the rules from the specified **table** and **chain**.

## **List All Rules**

```
GET /v1/iptables/list/all?verbose=
```

#### Query Parameters

- **verbose** - If parameter is **true**, List iptables rules with more detailed output.

List the rules from all tables and chains.

## **IPTable rules to JSON converter**

```
POST /v1/iptables/json
```

The request body contains `\n` delimited iptables rules. It will returns iptables rules represented in JSON. For more information on how it's works, checkout [this document](./docs/converter.md).

# **Contribution**

If you have any suggestions or issues then please open GitHub issue prefix with **`[opa-iptables]`**. Any pull request is most welcome.

## **Contribution History**

- List of all commits associated with this project: https://github.com/open-policy-agent/contrib/commits?author=urvil38

- List of Pull Request associated with this project: https://github.com/open-policy-agent/contrib/pulls?utf8=%E2%9C%93&q=is%3Apr+author%3Aurvil38+is%3Amerged