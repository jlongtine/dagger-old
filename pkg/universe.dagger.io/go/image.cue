package go

import (
	"universe.dagger.io/docker"
)

// Go image default version
_#DefaultVersion: "1.18"

// Build a go base image
#Image: {
	version: *_#DefaultVersion | string

	packages: [pkgName=string]: version: string | *""

	// FIXME Basically a copy of alpine.#Build with a different image
	// Should we create a special definition?
	docker.#Build & {
		steps: [
			docker.#Pull & {
				source: "index.docker.io/golang:\(version)-alpine"
			},
			for pkgName, pkg in packages {
				docker.#Run & {
					command: {
						name: "apk"
						args: ["add", "\(pkgName)\(pkg.version)"]
						flags: {
							"-U":         true
							"--no-cache": true
						}
					}
				}
			},
		]
	}
}
