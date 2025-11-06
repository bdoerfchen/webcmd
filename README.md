# Overview
<img src="assets/logo/webcmd.svg" height="100">

> "webcmd is your fast and simple bridge between the web and shell scripts with focus on customizability and security!"

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

### Exec Mode

### Parameters
For webcmd, route parameters are the interface between http requests and the executing process. For a given route, you can define a list of parameters that are used to translate HTTP request input into environment variables - readable by shell scripts and normal executables alike. Values can be retrieved from a request's route, query or header. Additionally you can define constants that are injected as environment variables as well.

#### Define Parameter

To define which parameters a user can provide, you need to declare them in the `parameters` list of a `Route`. For each `RouteParameter` you can provide the following values:
| Field     | Mandatory | Description |
| --------- | --------- | ----------- |
| `source`  | no        | The parameter source. |
| `name`    | yes       | The name of a parameter, mostly used to retrieve a parameters value from its source. | 
| `as`      | no        | Can be used to define the environment variable name for this parameter. By default it is `WC_{upper(.name)}`. | 
| `default` | no        | Value that is used when the user input is empty - or for constants where no user input is ever provided. |
| `disableSanitization` | no | Disable input value sanitization for this parameter. Refer to [Input Sanitization](#input-sanitization).

These are the valid sources:
| Source   | Value from        | Example |
| -------- | ------------ | ------- |
| `route`    | HTTP route by `name` (case-sensitive)   | `/hello/{planet}` -> `name` needs to be "planet". |
| `query`    | HTTP request query by `name` (case-sensitive) | `/hello?planet=mars` -> `name` needs to be "planet" and will have the value "mars". |
| `header`   | HTTP request header by `name` (case-insensitive) | `User-Agent: curl/7.81.0` -> `name` needs to be "user-agent". |
| `""`       | When the `source` field is omitted or set to `""` the parameter is a constant whose value is coming from `default`. | |

> [!TIP]  
> You can find an example configuration in [/examples/parameters](/examples/parameters/server.config.yaml)

#### Using the Parameters
For every incoming request, the parameter configuration is used to read the user-provided values into environment variables of the execution environment.

_But what are the names of the environment variables?_  
In general, the names have to follow this pattern: `[A-Za-z_][0-9A-Za-z_]`. What does it mean? The name must only consist of ASCII characters, numbers and underscores - but the first character must not be a number.

This is how webcmd maps parameter names into variable names:  
- With no custom name defined with `as`, the name is set to `WC_{uppercase(param.name)}`
- Otherwise the value from `as` is used
- Prohibited characters are replaced with `_` (whitespace, dashes, etc.)
- If the custom name in `as` starts with a number, the name is prefixed with `WC_`



# Security
"Bridging shell scripts and the web" is powerful but also comes with a risk. These guidelines may help you to reduce the risk of an attack:
- Run webcmd with the least amount of required permissions. Avoid running as root.
- Expect malicious user input when writing commands.
- Check user input.
- Run webcmd in an isolated environment. In a container for example.
- Only install necessary packages into your execution environment.

## Input Sanitization
Executing shell commands with custom user input does not only sound dangerous - it is. That is why all user input is undergoing santization before being exported as an environment variable. This comes at a small performance cost, but is worth it as a first line of defense to prevent users from executing arbitrary code on the server.  
It is worth reading more about [Command Injection](https://owasp.org/www-community/attacks/Command_Injection).

For webcmd the sanitization is quite simple and works by removing certain characters: 
> ; # % $ " ` ' & |

If this leads to some unwanted behaviour, sanitization can be turned off for each parameter individually by setting `disableSanitization` to `true`.

> [!CAUTION]  
> Disabling santization is a huge risk with the `shell` executer and thus routes where it is disabled should only use `proc` for execution!  
> Find more information in [/examples/attack](/examples/attack/server.config.yaml)



# Contribution



# Issues