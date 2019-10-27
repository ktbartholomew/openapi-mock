# OpenAPI Mock

OpenAPI mock reads an OpenAPI spec file and generates a server that emulates the behavior described in the spec file. Templates can be used within the spec file to make the output of responses somewhat dynamic. Clients can also use HTTP headers to control the output of the mocked API, such as triggering specific error responses for an endpoint.

## Installation

The application is distrubuted as a Docker container, so use `docker pull` to install it. All examples below assume you're running the application in a container.

```console
$ docker pull quay.io/ktbartholomew/openapi-mock
```

Alternatively, clone this repo and run `make build` to build the application locally, then copy it to a location in your PATH:

```console
$ make build
$ mv ./openapi-mock /usr/local/bin/openapi-mock
```

## Usage

```console
$ docker run -v /path/to/swagger.yml:/swagger.yml -p 3000:3000 quay.io/ktbartholomew/openapi-mock --spec-path /swagger.yml
```

### Flags

- `--spec-path`: A filesystem path to the OpenAPI spec YAML file
- `--spec-url`: A URL from which to fetch the OpenAPI spec YAML file. If `--spec-path` is set, this option is ignored.
- `--listen-addr`: The address (IP and port) on which the application will listen. Defaults to `0.0.0.0:3000`.

### Request Headers

The mock looks for certain request headers that a user can send to adjust the mock's output. Examples include simulating latency and requesting failure or other non-standard responses from an endpoint.

- `X-Mock-Latency`: A number of milliseconds to delay before the mock sends a response
- `X-Mock-Response`: An HTTP status code number like 202 or 403 to receive. If not provided, the lowest-numbered status code from the spec will be returned. Requesting a status code not defined in the spec will return a 400 (Bad Request) response.
- `X-Mock-Count`: The number of items to be inluded in response collections. This requires corresponding [response templates](#) in order to be effective.

### Templating OpenAPI Specs

The response body of each endpoint in the mock is derived from the [example](https://swagger.io/docs/specification/adding-examples/) field. This works best when the example string is a single string, and not the object property format that OpenAPI also supports:

```yaml
responses:
  '200':
    description: A list of user objects
    content:
      application/json:
        example: |
          [
            {
              "id": 1,
              "email": "albert@example.com"
            },
            {
              "id": 2,
              "email": "brenda@example.com"
            },
            {
              "id": 3,
              "email": "cathy@example.com"
            }
          ]
```

Maintaining examples with more complicated data structures or larger collections can be very tedious to do by hand in a single, static string. To address this, the mock supports templating with Go templates to make realistic API outputs easier to mock. Using this template capability, the example above can be written like this instead:

```yaml
responses:
  '200':
    description: A list of user objects
    content:
      application/json:
        example: |
          {{ .JSONArray {{ .ItemCount }} `{
            "id": {{ .Index }},
            "email": "{{ .ToLowerCase .RandomFirstName }}@example.com"
          }` }}
```

All templates receive a pipeline that exposes a number of useful methods for generating API output.

#### `{{ .JSONArray <n> <string> }}`

Repeats `<string>` `<n>` times, wrapping the entire output in square brackets and adding commas between each entry so that the resulting output is a valid JSON array. The contents of `<string>` are also parsed as a template, so properties like `{{ .Index }}` can be used within the string to increment numbers for each element in the array.

#### `{{ .ToLower <string> }}`

Returns `<string>`, converted to lower-case.

#### `{{ .RandomFirstName }}`

Returns a random first name.

#### `{{ .RandomFrom <...string> }}`

Returns one of the provided strings at random

#### `{{ .RandomPassword <len> }}`

Returns a random string (like a strong password) of length `<len>`.