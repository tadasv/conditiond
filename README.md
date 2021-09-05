[![Go](https://github.com/tadasv/conditiond/actions/workflows/go.yml/badge.svg)](https://github.com/tadasv/conditiond/actions/workflows/go.yml)

# conditiond

`conditiond` is a generic constraint and policy evaluator.

This tool lets you define constraints in *data* and evaluate them at run time.
It's designed to be run as a container sidecar but it can also be used from the
command line and integrate with your shell scripts.

## Installation

```
go install github.com/tadasv/conditiond/cmd/conditiond@latest
```

Then run `conditiond -h`.

## Example

```sh
$ cat input
{
    "condition": {
        "and": [
            {"if": [
                {"eq": [{"context": ["user_id"]}, 123]},
                true
            ]}
        ]
    },
    "context": {
        "user_id": 123
    }
}
{
    "condition": {
        "and": [
            {"if": [
                {"eq": [{"context": ["user_id"]}, 123]},
                true,
                false
            ]}
        ]
    },
    "context": {
        "user_id": "not 123"
    }
}
$ cat input | ./conditiond -cli
{"error":null,"result":true}
{"error":null,"result":false}
```

Above example passes in two condition definitions and context associated with
each of them. The first condition definition checks whether the `user_id`
matches `123` and returns `true` if that's the case. The second one is the
same, but we have a different user_id in the provided context which will result
in a different result value.

We can achieve the same by invoking the evaluator via HTTP RPC:

```sh
$ ./conditiond &
[1] 21780
$ 2021/09/04 10:25:05 starting conditiond server on :9000
$ curl -d @input localhost:9000/evaluate
{"error":null,"result":true}
{"error":null,"result":false}
```

## Why do we need this?

Sometimes we want to create our own policies or constraints, but manage them in
data instead of updating code and shipping a new release. `conditiond` enables
that. Such rules can be managed by other people outside of engineering via
some nice UI requiring almost no code changes once integration with your backend
is complete.

`conditiond` can serve as a building block for

- Access control policies.
- AB tests and experiments.
- Feature flag toggles.

For example, we could setup an experiment where we assign 50% of the users to
`cohort-a` and another 50% to `cohort-b`:

```
{
  "condition": {
    "if": [
      {
        "gt": [
          5,
          {
            "sha1mod": [
              {"context": ["user"]},
              10
            ]
          }
        ]
      },
      "cohort-a",
      "cohort-b"
    ]
  },
  "context": {
    "user": "some user id"
  }
}
```


## The Protocol

`conditiond` operates on a stream of JSON messages. These messages can be
passed in via CLI or HTTP RPC. A stream is created by simply concatenating
several JSON messages. The messages **may** be evaluated out of order but the
result messages will always be returned in the same order as the input so they
can be indexed the same way.

For a given two message input stream

```
{ input message 1 }
{ input message 1 }
```

We are going to return a stream of results

```
{ results for message 1 }
{ results for message 2 }
```

The request message is a JSON object:

```
{
  "condition": ...
  "context": ...
}
```

Here, `condition` contains an expression (see Expression Specification).
Optionally, a `context` object can be provided. This context object is passed
into every expression at evaluation time so expressions that need outside
information can utilize it.

Here's an example of a full request message:

```
{
  "condition": {
    "gt": [{"context": ["monthly_spend"]}, 10000]
  },
  "context": {
    "user_id": "123",
    "monthly_spend": 5555
  }
}
```

The result message is a JSON object of the following form:

```
{
  "error": ...
  "result": ...
}
```

The `error` key will be set to a string containing an error message if
evaluation failed for some reason. The `error` will be null otherwise and
`result` key will contain `condition` evaluation result.

## Expression specification

Expressions in `conditiond` are designed after
[S-Expressions](https://en.wikipedia.org/wiki/S-expression) but encoded as a
subset of JSON.

An expression takes a form of a JSON object:

```
{
  "<expression-name>": [<expression-argument>, ...]
}
```

The object **must** contain a single key, `<expression-name>`. The key must
point to a JSON array of 0 or more `<expression-argument>` values.

`<expression-argument>` can be another expression object or any of the JSON
literals (string, number, boolean or null).


## Available expressions

### and

Returns `true` when all arguments evaluate to `true`. If argument list is empty
the result will be true.

Examples:

```
{
  "and": [true, false]
}
```

### or

Returns `true` when some of the arguments evaluate to `true`. If argument list
is empty the result will be `false`.

Examples:

```
{
  "or": [true, false]
}
```

### not

Negates the evaluation result of it's argument. Requires exacly one argument to
be passed in.

Examples:

```
{
    "not": [true]
}
```

### gt

Returns `true` if the first argument is greater than the second argument. It
requires exactly two arguments, which when evaluated must return numbers.

Examples:

```
{
    "gt": [123, 321]
}
```

### lt

Returns `true` if the first argument is less than the second argument. It
requires exactly two arguments, which when evaluated must return numbers.

Examples:

```
{
    "lt": [123, 321]
}
```

### gte

Returns `true` if the first argument is greater or equal to the second argument. It
requires exactly two arguments, which when evaluated must return numbers.

Examples:

```
{
    "gte": [123, 321]
}
```

### lte

Returns `true` if the first argument is less or equal to the second argument. It
requires exactly two arguments, which when evaluated must return numbers.

Examples:

```
{
    "lte": [123, 321]
}
```

### eq

Returns `true` if two arguments are equal. It requires exactly two arguments to
be passed in.

Examples:

```
{
    "eq": ["123", "123"]
}
```

NOTE This function does not perform type coersion. E.g.

```
{
    "eq": ["123", 123]
}
```

Will return `false`.

### sha1mod

Takes two arguments. The first argument is hashed with SHA1. Second argument is
used to perform a mod operation with the SHA1 output. The remainder of the mod
operation is returned as a result.

Examples:
```
{
    "sha1mod": ["some data", 15]
}
```

### context

Extracts value from a provided context. Arguments represent path to the field
we want to extract. The extracted value is returned as is and no type coersion
is performed. It returns null value if no data exists at the path.

Examples:

```
{
    "context": ["key", 1, "key2"]
}
```

With provided context:

```
{
  "key": [
    123,
    {
      "key2": "value"
    },
    "test"
  ]
}
```

Will return `value` string.

### if

Requires 2 or 3 arguments and returns second argument if the first argument
evaluates to `true`. Otherwise returns 3 argument or null value if the first
argument evaluates to `false`.

In other languages this could be written as:

```
if (arg1) {
    return arg2
} else {
    return arg3
}
```

Example:

```
{
    "if": [true, "value1", "value2"]
}
```
