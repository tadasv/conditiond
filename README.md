# conditiond

`conditioned` is a generic constraint and policy evaluator.

This tool lets you define constraints in *data* and evaluate them at run time.
It's designed to be run as a container sidecar but it can also be used from the
command line and integrate with your shell scripts.

## Installation

```
go install github.com/tadasv/conditiond/cmd/conditiond@latest
```

Then run `conditioned -h`.

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

TODO list and describe available expressions.
