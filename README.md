# import-layers-go

Checks the direction of imports between layers of a Go project.

## Usage

Build or download the binary and run it from a project containing a `go.mod`:

```sh
import-layers-go
import-layers-go ./...
import-layers-go cmd internal
import-layers-go -c config/layers.json internal
```

The default scan starts from all Go packages under the project root. Pass paths relative to the current directory to narrow the scan. Bare directory paths are recursive. Native Go package patterns such as `./cmd` and `./cmd/...` retain their Go meaning. You can also pass individual Go source files. The analyzer follows imports into local dependencies and respects Go build constraints.

## Configuration

The command searches upward for the nearest `import_layers.json`, `import_layers.yaml`, or `import_layers.yml`, in that order within each directory. The containing directory defines the project root. Use `-c` or `--config` to select a different configuration file. Relative config paths are based on the current directory.

```json
{
  "layers": ["domain", "service", "cmd"],
  "debug": false,
  "ignoreTests": false
}
```

A package may import from its own layer and earlier layers. For example, `service` can import `domain`. A layer name matches a complete Go import path segment. A project-relative directory such as `internal/domain` or a qualified Go import path takes priority over a single segment. Go import path selectors can also target explicitly configured third-party packages.

Set `rootDir` to a directory relative to the project root to override the positional scan paths. Set `ignoreTests` to skip `*_test.go` files. The legacy spelling `ignore_tests` is also supported. With `debug: true` the command prints the relative path of each checked file.

See [the example configuration](examples/basic/import_layers.yaml).

## Output

A successful check prints `OK`. Diagnostics go to stderr and include a source path, line, and imported package. The exit status is `0` for success, `1` for layer violations, or `2` for argument, configuration, or analysis errors.

Pass `--graph` to write a JSON dependency graph with `nodes` and `edges` to stdout. Each edge connects a source file to a local Go source file in the imported package. Debug output goes to stderr in this mode. The Go vet protocol and native `-json` diagnostics remain available.
