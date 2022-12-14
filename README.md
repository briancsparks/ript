# Hello Fooprojectname

## TODO: version number

(And how to do isActiveDevelopment)
You can use the linker to set variables inside the app.

See:

- Explaining how the Go build tool helps for this
  - https://github.com/golang/go/wiki/GcToolchainTricks#including-build-information-in-the-executable
- Lots of examples
  - https://stackoverflow.com/questions/11354518/application-auto-build-versioning
- Article(s)
  - https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
  - Using git tags:
    - https://www.forkingbytes.com/blog/dynamic-versioning-your-go-application/


## Ript is Replace-In-File sans the Template Engine

Template engines are great, but are overkill in many situations. If you don't need the power -
like when generating a project's files - it is convenient to just have files - sans template engine.
For example, consider this Go file:

```go
package riptprojectname

/* Copyright © 2022 RIPTENV_IMFO_USER_IRL_NAME <RIPTENV_IMFO_GITHUB_USEREMAIL> -- MIT (see LICENSE file) */

import "github.com/RIPTENV_IMFO_GITHUB_USERNAME/riptprojectname/cmd"

func main() {
  cmd.Execute()
}
```

* Tools can parse 100% of the file
  * No getting tripped up on the template engine's markup.
* All the items that start `ript` will be replaced with a simple text replacement.
  * All your template files are valid as-is, so all your tools will work with them.

### Simple Rules

1. Any tokens that start `RIPTENV_` will be replaced by that environment variable.
   1. If it is absent in the environment, a default is read from a YAML config file.
2. Any tokens that start `ript` are replaced by whatever the user entered on the
   command-line.

