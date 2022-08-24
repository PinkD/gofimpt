# gofimpt

`gofimpt` is designed to format imports in a go file

It will separate imports into 3 parts:

- builtin
- 3rd-party
- module

And here's an exmple:

```go
import (
	"fmt"
	"github.com/pkg/errors" 
	// comment 1
	// comment 2
	"path"
	// import os
	"os"
	"gofimpt/errors"
)
```

The format result will be like:

```go
import (
	"fmt" 
	// import os 
	"os"
	// comment 1
	// comment 2
	"path"

	"github.com/pkg/errors"

	"gofimpt/errors"
)
```

```text
Usage:
  gofimpt [file/dir...]

    if no arg is provided, all files tracked by this git repo will formatted
    if files or directories are provided, provided files and all files under directories will be formatted
```

> NOTE: block comment is not supported

The name `gofimpt` comes from those keywords:

- **go**
- **fi**le
- **f**or**m**a**t**
- **imp**or**t**


