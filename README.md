# Overview
<img src="webcmd.svg" height="100">

> "webcmd is a simple webserver with the aim to bridge HTTP and shell scripts"

It aims to be
- Easy to configure
- Fast
- Reliable
- Secure

> [!NOTE]  
> This documentation is currently worked on and will be filled in the next weeks.

# Configuration

## Route

### Path

### Method

### Parameters
For webcmd, route parameters are the interface between http requests and the executing process. For a given route, you can define a list of parameters that are used to translate HTTP request input into environment variables - readable by shell scripts and executables alike. Values can be retrieved from a request's route, query or header. Additionally you can define constants that are injected as environment variables as well.

To define which parameters a user can provide, you need to declare them in the `parameters` list of a `Route`. For each `RouteParameter` you can provide the following values:
| Field   | Optional | Description |
| ------- | -------- | ----------- |
| source  | yes      | The parameter source. |
| name    | no       | The name of a parameter, mostly used to retrieve a parameters value from its source. | 
| as      | yes      | Can be used to define the environment variable name for this parameter. By default it is `WC_{upper(.name)}`. | 
| default | yes      | Value that is used when the user input is empty - or for constants where no user input is ever provided. |

These are the valid sources:
| Source   | Value from        | Example |
| -------- | ------------ | ------- |
| `route`    | HTTP route by `name` (case-sensitive)   | `/hello/{planet}` -> `name` needs to be "planet". |
| `query`    | HTTP request query by `name` (case-sensitive) | `/hello?planet=mars` -> `name` needs to be "planet" and will have the value "mars". |
| `header`   | HTTP request header by `name` (case-insensitive) | `User-Agent: curl/7.81.0` -> `name` needs to be "user-agent". |
| `""`       | When the `source` field is omitted or set to `""` the parameter is a constant whose value is coming from `default`. | |

An example can be found in [/examples/parameters](/examples/parameters/server.config.yaml)

### Exec


# Contribution

# Issues