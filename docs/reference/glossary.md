# Glossary

## Builder

A *builder* is the implementation for a given object of the *builder pattern*,
allowing to construct complex objects incrementally.

## Builder option

A *builder option* is a method exposed by a builder, responsible for manipulating
a property of the object being built.

## Intermediate representation

The intermediate representation (IR) is the result of parsing the schemas listed
as inputs of the codegen pipeline.

The IR is what will be modified by schema and builder transformations.

*Jennies* use this IR to generate code. 

## Jennies

`cog` calls "*jennies*" its components that are responsible for generating code
in a specific language.

Example: the component generating types in Go is called the "Go types" jenny.

## Pipeline

A pipeline – or *codegen pipeline* – is the configuration file used by `cog`
to generate code.

A typical pipeline contains schemas used as inputs and configuration describing
the desired outputs.
