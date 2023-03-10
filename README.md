[![License: GPL v3](https://img.shields.io/github/license/anthonyme00/go-envloader)](https://github.com/anthonyme00/go-envloader/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/anthonyme00/go-envloader)](https://pkg.go.dev/github.com/anthonyme00/go-envloader)
![Travis (.com)](https://app.travis-ci.com/anthonyme00/go-envloader.svg?branch=main)
[![Coverage Status](https://coveralls.io/repos/github/anthonyme00/go-envloader/badge.svg)](https://coveralls.io/github/anthonyme00/go-envloader)
[![Go Report Card](https://goreportcard.com/badge/github.com/anthonyme00/go-envloader)](https://goreportcard.com/report/github.com/anthonyme00/go-envloader)

# go-envloader
Easily load your go configuration structs from your environment variables

# Import Guide:

Run the following command from your terminal

```shell
go get "github.com/anthonyme00/go-envloader"
```
After which you can import this package with the following line
```go
import (
    envloader "github.com/anthonyme00/go-envloader"
)
```

# Usage
You will need to define environment variable keys to your existing struct fields by using the `env` tag.

Currently, there are 2 supported attributes:

1. `key` - **REQUIRED**     : The key for getting your environment variable
2. `default` - *OPTIONAL*   : Default value for this field in case your environment variable doesn't contain an entry with key = `key`

Attributes uses the colon (:) for definition, and are separated with semicolon (;)

Example:
```go
type ExampleStruct struct {
    AppName string  `env:"key:APP_NAME;default:UNIPROJECT"`
    AppPort int     `env:"key:APP_PORT"`
}
```

This library also supports nested structs whether it's inline or not

Example:
```go
package main

import (
	"github.com/anthonyme00/go-envloader"
)

type AppConfig struct {
	AppName string `env:"key:APP_NAME;default:UNIPROJECT"`
	AppPort int    `env:"key:APP_PORT;default:8081"`
}

type DBConfig struct {
	// you can use colons after the identifying atrribute key
	// colons after the first one are ignored and treated as
	// part of the string
	DBUri string `env:"key:DB_URI;default:192.168.0.1:8080"`
}

type Config struct {
	App AppConfig
	DB  DBConfig
}

func main() {
	config := Config{}

	loader := envloader.New(nil)

	errs := loader.Load(&config)
	if len(errs) > 0 {
		// Error handling here
	}

	// config is ready to be used :)
}
```

Slices are also supported. Slice elements are separated with comma (,)
Example:
```go
type ExampleStruct struct {
    StringSlice []string  `env:"key:STRING;default:STRING 1,STRING 2"`
    Int64Slice  []int64   `env:"key:INT64;default:3,5,7"`
}
```

You are also able to create default configuration files using the `Stringify` function

Example:
```go
package main

import (
	"fmt"

	"github.com/anthonyme00/go-envloader"
)

type AppConfig struct {
	AppName string `env:"key:APP_NAME;default:UNIPROJECT"`
	AppPort int    `env:"key:APP_PORT;default:8081"`
}

type DBConfig struct {
	DBUri string `env:"key:DB_URI;default:192.168.0.1:8080"`
}

type Config struct {
	App AppConfig
	DB  DBConfig
}

func main() {
	config := Config{}

	loader := envloader.New(nil)

	env, _ := loader.Stringify(&config)
	fmt.Println(env)
}
```

Output:

```

APP_NAME = UNIPROJECT
APP_PORT = 8081

DB_URI = 192.168.0.1:8080

```

# Supported types
- `int64` : single, slice
- `int32` : single, slice
- `int16` : single, slice
- `int8` : single, slice
- `int` : single, slice
- `float64` : single, slice
- `float32` : single, slice
- `string` : single, slice
- **custom** `struct` : single

## Special Types
These types are casted from string to its respective type. Follow this guideline for these types
- `bool` : single, slice <br>
	values :
	- `"true"`= `true`
	- `"false"` = `false`

# Limitations
- Does not support pointer types
