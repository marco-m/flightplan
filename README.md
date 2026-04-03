# flightplan

Generate Concourse pipelines programmatically with Go.

Zero dependencies for production code (and minimal dependencies for test code).

WORK IN PROGRESS. NOT USABLE.

This work stems from the realization that configuration files in YAML (or in JSON) become unmanageable as soon as the file gets long enough. In particular, Concourse pipelines tend to be big and so difficult to edit; even worse attempting to keep multiple pipelines to use the same patterns.

Instead, generating a pipeline configuration programmatically, with a statically typed language, allows to:

- Handle complexity and refactoring with a programming language: full editor support, Language Server support.
- Handle duplication with a programming language: instead of YAML anchors, just use a language variable.
- Handle "customizing" pieces of pipelines with a programming language, instead of YAML anchors or YAML manipulation.
- Perform many consistency checks before hitting the Concourse runtime. This can be coupled to great advantage with Open Policy Agent, that is supported by Concourse.

In particular with Go:

- Handle common code in multiple pipelines in different repositories as an external Go package: `go get site/common-code` (see examples). This is the same approach of [Florist](https://github.com/marco-m/florist).

## AI policy

This project does not accept contributions of any form (code, documentation, issues) written with the aid of AI agents.
