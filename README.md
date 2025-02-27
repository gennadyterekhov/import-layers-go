# import-layers-go

import-layers-go is a golang checker to enforce abstraction layers
so that they don't mix up, and you can create coherent layered architecture  
  
in other words - check that higher layer packages do not depend on lower layer packages (dependency rule from clean architecture)

## example

config

    {
        "layers": [
            "high",
            "low",
        ]
    }


go

    package high

    import (
        "low" // returns error: `cannot import package from lower layer`
    )

## running in your repo

download bin from releases and place it in `analyzers` for example  
run

    ./analyzers/import-layers-go ./...

## config

See example config in [examples/basic/import_layers.yaml ](https://github.com/gennadyterekhov/import-layers-go/blob/main/examples/basic/import_layers.yaml)  
Config file must be in the same directory as a `go.mod` file.  
It must be named `import_layers.yaml`.  
Config file name and location are not configurable.  

