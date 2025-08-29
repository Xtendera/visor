# Visor

Fast API monitor built for consistency and reliability. The configuration-first approach keeps it simple and flexible.

## Installation
You can currently install through go:
```
go install github.com/Xtendera/visor@latest
```
Or a specific version:
```
go install github.com/Xtendera/visor@0.0.1a6
```

You may also be required to add GOPATH into your system PATH. Use the following command to append this into your shell file. Replace `.bashrc` with the file specific to your shell:
```
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
```
## Usage

```bash
visor run <CONFIG_FILE>.json
```

## Configuration

The config is stored in a JSON file. Please look at `config.example.json` for a functional configuration file. **NOTE:** You will also need the `example/` directory in the same location as the example configuration for the project to run. You will find the config and the example directory in this git project.

Here are the possible properties for the config:

`root` (string): The base url for all endpoints to run on. This should contain the scheme alongside the absolute hostname WITHOUT a leading slash (`/`).

`headers` (object array): Each object contains a key and value property (both of which are type string) correlating with the key and value of a header.

`endpoints` (object array): Describes each endpoint to test. For all properties, see the documentation below.

`cookies` (object array): Each object contains a `key` and `value` property (both strings) representing cookies to be sent with each request.

### Endpoints

Each endpoint object contains the following additional properties:

`name` (string): The endpoint name (used for logging purposes).

`path` (string): The relative path for this endpoint. It should be prefixed with a trailing slash (`/`).

`method` (string): The HTTP method to use. Must be one of: `GET`, `POST`, `PUT`, `HEAD`, `DELETE`, `OPTIONS`, `PATCH`.

`headers` (object array, optional): Additional headers specific to this endpoint. Each object contains a `key` and `value` property.

`body` (object, optional): The request body to send. If provided, it will be sent as JSON if the value is an object or array, otherwise as plain text.

`bodyFile` (string, optional): A file which will be read and sent as the request body. It will always be sent as plaintext, you can manually set the header through the headers property to circumvent this.

`acceptStatus` (array of integers): List of HTTP status codes that are considered successful for this endpoint. At least one value is required.

`schema` (string, optional): Path to a JSON schema file used to validate the response body.


## Licensing

This project is under the [Mozilla Public License 2.0](https://github.com/Xtendera/Visor/blob/main/LICENSE).

Development is sponsored by [ĀYŌDÈ](https://ayode.org). In addition to the rights granted under the MPL 2.0, I hereby grant ĀYŌDÈ a non-exclusive, irrevocable, royalty-free license to use, modify, sublicense, and distribute this code under any terms, including proprietary licenses, without restriction.
