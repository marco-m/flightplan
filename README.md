# flight-plan

Generate Concourse pipelines programmatically with Go.

WORK IN PROGRESS. NOT USABLE.

This work stems from the realization that configuration files in YAML (or in JSON) become unmanageable as soon as the file gets long enough. In particular, Concourse pipelines tend to be big and so difficult to edit; even worse attempting to keep multiple pipelines to use the same patterns.

Generating a pipeline configuration programmatically, with a statically typed language, allows to:

- Handle complexity and refactoring with a programming language: full editor support, Language Server support.
- Handle duplication with a programming language: instead of YAML anchors, just use a language variable.
- Perform many consistency checks before hitting the Concourse runtime. This can be coupled to great advantage with Open Policy Agent, that is supported by Concourse.

In particular with Go:

- Handle common code in multiple pipelines in different repositories with the Go package facility, a natural fit (see examples). This is the same approach of [Florist](https://github.com/marco-m/florist).
- Never introduce errors possible with dynamically typed languages.
