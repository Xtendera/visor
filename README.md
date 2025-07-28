# Visor

Fast API monitor built for consistency and reliability. The configuration-first approach keeps it simple and flexible.

## Usage

```bash
visor run <CONFIG_FILE>.json
```

## Configuration

The config is stored in a JSON file. Please look at `config.example.json` for a functional configuration file.

Here are the possible properties for the config:

`root` (string): The base url for all endpoints to run on. This should countain the scheme alongside the absolute hostname WITHOUT a leading slash (`/`).

`headers` (object array): Each object contains a key and value property (both of which are type string) corrolating with the key abd value of a header.

`endpoints` (object array): Describes each endpoint to test. For all properties, see the documentation below.

### Endpoints

Each endpoint object contains the following properties.

`name` (string): The endpoint name (used for logging purposes).

`path` (string): The relative path for this endpoint. It should be prefixed with a trailing slash (`/`).

## Licensing

This project is under the [Mozilla Public License 2.0](https://github.com/Xtendera/Visor/blob/main/LICENSE).

Development is sponsored by [AYODE](https://ayode.org). In addition to the rights granted under the MPL 2.0, I hereby grant AYODE a non-exclusive, irrevocable, royalty-free license to use, modify, sublicense, and distribute this code under any terms, including proprietary licenses, without restriction.
