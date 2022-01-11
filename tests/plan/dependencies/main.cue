package main

import (
	"alpha.dagger.io/europa/dagger/engine"
)

engine.#Plan & {
	actions: {
		pull: engine.#Pull & {
			source: "alpine:3.15.0@sha256:e7d88de73db3d3fd9b2d63aa7f447a10fd0220b7cbf39803c803f2af9ba256b3"
		}

		sayHello: engine.#Exec & {
			input: pull.output
			args: [
				"sh", "-c",
				#"""
					echo -n Hello Europa! > /out.txt
					"""#,
			]
		}

		readfile: engine.#ReadFile & {
			input: sayHello.output
			path:  "/out.txt"
		} & {
			// assert result
			contents: "Hello Europa!"
		}

		#things: {
			require: [...]
			sayHello2: engine.#Exec & {
				input:     pull.output
				"require": require
				args: [
					"sh", "-c",
					#"""
						echo -n Hello Europa! > /out.txt
						"""#,
				]
			}
			readfile: engine.#ReadFile & {
				input: sayHello2.output
				path:  "/out.txt"
			} & {
				// assert result
				contents: "Hello Europa!"
			}

			sayHello3: engine.#Exec & {
				input: sayHello2.output
				args: [
					"sh", "-c",
					#"""
						echo -n Hello from Dagger > /world.txt
						"""#,
				]
			}
			verifyWorld: engine.#ReadFile & {
				input: sayHello3.output
				path:  "/world.txt"
			} & {
				// assert result
				contents: "Hello from Dagger"
			}
		}

		todo: #things & {
			require: [readfile]
		}
	}
}
