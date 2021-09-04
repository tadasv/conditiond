# conditiond

`conditioned` is an arbitrary constraint and policy evaluator.

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

## Condition spec

Conditions are designed after Lisp's S-Expression but encoded as JSON (because
everyone uses JSON these days).

TODO write up what's allowed and what's not.

## Supported expressions

TODO list and describe available expressions.

## TODO

- [ ] http interface
