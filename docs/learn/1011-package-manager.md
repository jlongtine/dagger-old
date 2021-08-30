---
slug: /1011/package-manager/
---

# Manage packages using the package manager

This tutorial illustrates how to install and upgrade packages using Dagger package manager.

## Installing a package

### Initializing workspace

Create an empty directory for your new Dagger workspace:

```shell
mkdir workspace
cd workspace
```

As described in the previous tutorials, initialize your Dagger workspace:

```shell
dagger init
dagger new test
```

That will create 2 directories: `.dagger` and `cue.mod` where our package will reside:

```shell
.
├── cue.mod
│   ├── module.cue
│   ├── pkg
│   └── usr
├── .dagger
│   └── env
│       └── test
```

### Install

In our example we will use `gcpcloudrun` module from [github](https://github.com/tjovicic/dagger-modules/blob/main/gcpcloudrun/source.cue)
Let's first add it to our `source.cue` file:

```cue title="./source.cue"
package main

import (
  "github.com/tjovicic/dagger-modules/gcpcloudrun"
)

run: gcpcloudrun.#Run
```

To install it just run

```shell
dagger mod get github.com/tjovicic/dagger-modules/gcpcloudrun@v0.1
```

It should pull the `v0.1` version from Github, leave a copy in `cue.mod/pkg` and reflect the change in
`cue.mod/dagger.mod.cue` file:

```shell
cue.mod/pkg/github.com/
└── tjovicic
    └── dagger-modules
        └── gcpcloudrun
            ├── cue.mod
            ├── README.md
            └── source.cue
```

```cue title="./cue.mod/dagger.mod"
github.com/tjovicic/dagger-modules/gcpcloudrun v0.1
```

Querying the current setup with `dagger query` should return a valid result:

```json
{
  "run": {
    "creds": {
      "username": "oauth2accesstoken"
    },
    "deploy": {
      "platform": "managed",
      "port": "80"
    },
    "push": {
      "auth": {
        "username": "oauth2accesstoken"
      },
      "push": {}
    }
  }
}
```

### Upgrading

Now that you've successfully installed a package, let's try to upgrade it.

```shell
dagger mod get github.com/tjovicic/dagger-modules/gcpcloudrun@v0.2
```

You should see similar output:

```shell
12:25PM INF system | downloading github.com/tjovicic/dagger-modules:v0.2
```

And `cue.mod/dagger.mod.cue` should reflect the new version:

```cue title="./cue.mod/dagger.mod"
github.com/tjovicic/dagger-modules/gcpcloudrun v0.2
```