# Plugins

This package provides the infrastructure for registering PhotoPrism plugins.

Currently the plugin system provides hooks into the following internals:

- `indexing`: Allowing to run different tagging and face detection algorithms or adding and modifying media metadata.

## Implementation

Check the demo plugin implementation in [demo.go](./demo/demo.go) and the `build-plugin-demo` make target in the [Makefile](../../Makefile) on how to compile the plugin solib.

## Configuration

The plugin can be configured using environment variables, which will be passed as key-value pairs to the plugin's `Configure` method.
All configuration parameters will be namespaced and the plugin will only have access to the variables prefixed with `PHOTOPRISM_PLUGIN_[UPPERCASE_PLUGIN_NAME]_`. The prefix will be stripped from the resulting parameters and the the rest of the parameter key will be lowercased.

For example a `demo` plugin might need a hostname and a port as configuration parameters, which can be achieved using the following environment variables:

```sh
PHOTOPRISM_PLUGIN_DEMO_HOSTNAME=localhost
PHOTOPRISM_PLUGIN_DEMO_PORT=1234
```

This configuration will then be passed to the plugin as the following key-value pairs:

```go
map[string]string{"hostname": "localhost", "port": "1234"}
```

This way you can safely pass any configuration parameters to your plugin.

## Activation

PhotoPrism will look in the `plugins` folder within the `PHOTOPRISM_STORAGE_PATH` folder and load all solib (.so) files found there. In order for the plugin to be loaded, it needs to fulfill few criterias:

- There should be `Plugin` variable that exposes the plugin implementation.
- Your plugin must implement all of the required interface method, even if you dont use them.
- Your plugin must be explicitly activated using the `enabled` configuration parameter: `PHOTOPRISM_PLUGIN_[UPPERCASE_PLUGIN_NAME]_ENABLED=yes/true/on/enable`.
