# Spring AccessDecisionVoter using OPA

This directory contains a simple implementation of an [AccessDecisionVoter for Spring Security](https://docs.spring.io/spring-security/site/docs/4.2.4.RELEASE/reference/htmlsingle/#authz-voting-based) that uses OPA for making authorization decisions.

## Prerequisites

- Java (tested with 1.8)
- Maven (tested with 3.3.9)

## Usage

To build the JAR file:

```bash
mvn package
```

To use the JAR file:

```bash
mvn install:install-file -Dfile=target/voter-1.0-SNAPSHOT.jar -DpomFile=pom.xml
```

Add a dependency on the package to your project (`pom.xml`):

```xml
<dependency>
	<groupId>org.openpolicyagent</groupId>
	<artifactId>voter</artifactId>
	<version>1.0-SNAPSHOT</version>
</dependency>
```

## Web Security Configuration

To enable the voter inside your application, you must configure it. [Spring Security](https://docs.spring.io/spring-security/site/docs/4.2.4.RELEASE/reference/htmlsingle/) has sophisticated support for XML and Java-based configuration.

The example below is a simplistic Java-based configuration that you can use to test the voter. Drop this file into your project.

```java
package com.acmecorp.example.config;

import java.util.Arrays;
import java.util.List;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.access.AccessDecisionManager;
import org.springframework.security.access.AccessDecisionVoter;
import org.springframework.security.access.vote.UnanimousBased;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.*;

import org.openpolicyagent.voter.OPAVoter;

@Configuration
@EnableWebSecurity
public class WebSecurityConfig extends WebSecurityConfigurerAdapter {

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        http.authorizeRequests().anyRequest().authenticated().accessDecisionManager(accessDecisionManager());
    }

    @Bean
    public AccessDecisionManager accessDecisionManager() {
        List<AccessDecisionVoter<? extends Object>> decisionVoters = Arrays
                .asList(new OPAVoter("http://localhost:8181/v1/data/http/authz/allow"));
        return new UnanimousBased(decisionVoters);
    }

}
```

## Testing

Obtain the latest version of OPA and start your application (e.g., using `mvn sprint-boot:run`).

Create a test policy (`example.rego`):

**example.rego**:

```ruby
package http.authz

allow = true
```

Run OPA in server mode with file watching enabled:


```bash
opa run -s -w example.rego
```

Test that you can access your application's API:

```bash
curl localhost:8080
```

Modify the policy to deny all requests.

**example.rego**:

```ruby
package http.authz

allow = false
```

Test that your application's API requests are rejected:

```bash
curl localhost:8080
```
