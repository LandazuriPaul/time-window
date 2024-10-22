# Time Window Validator

A GitHub Action to validate if a timestamp (the current one by default) is inside an `allowed` or `blocked` time window.

This is to prevent a workflow from running at certain times (e.g. out of office hours, bank holidays, etc.).

If needed, the validation can be by-passed with a commit message matching the `force-valid-regexp` (by
default, `force\-time\-window`, e.g. `hotfix: deploy this now, let's force-time-window!`).

## Usage

The action requires the `allowed` input and can have an optional `blocked` input. Unless the commit message forces the
validation, the timestamp is validated against both `allowed` and `blocked` time windows provided.

### Validation logic

A given timestamp is considered as valid only if it respects **both** rules:

- **it is within an `allowed` time window, and**
- **it is not within a `blocked` time window**.

This means that the `blocked` windows always take precedence over the `allowed` ones.
This enables configurations where the blocked time windows correspond to bank holidays or non-business periods to
overtake usual business hours.

### Inputs

| Parameter            | Default value                                                                                | Description                                                                                                                         |
|----------------------|----------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| `allowed`            | `""`                                                                                         | The allowed time windows in YAML. See [`allowed` & `blocked`](#allowed--blocked) for more information.                              |
| `blocked`            | `""`                                                                                         | The blocked time windows in YAML. See [`allowed` & `blocked`](#allowed--blocked) for more information.                              |
| `commit-message`     | `${{ github.event.head_commit.message \|\| github.event.workflow_run.head_commit.message }}` | The commit message against which the regular expression is validated. See [`commit-message`](#commit-message) for more information. |
| `force-valid-regexp` | `force\-time\-window`                                                                        | A regular expression to check against the commit message. See [`force-valid-regepx`](#force-valid-regexp) for more information.     |
| `timestamp`          | now                                                                                          | Unix timestamp (in seconds) to validate the time windows against. See [`timestamp`](#timestamp) for more information.               |

#### `allowed` & `blocked`

Each of these inputs receives an array of time windows defined in a YAML structure with the following fields:

- `name`: The name of the time window. Useful to understand the result.
- `cronExpression`: The CRON expression defining the beginning of the time window. You can use
  the [crontab guru](https://crontab.guru/) to generate your expression.
- `duration`: The duration of the time window.

> N.B.: The action currently doesn't support timezones and everything is computed using UTC. See the [TODOs](#todos)
> below.

#### `commit-message`

If you want to provide a specific commit message. By default, this is inferred from the workflow with the following
GitHub
Action expression:

```
${{ github.event.head_commit.message || github.event.workflow_run.head_commit.message }}
```

#### `force-valid-regexp`

The regular expression used to validate the commit message. If there is a match, the result is automatically set
to `is_valid=true`.

#### `timestamp`

The Unix timestamp (in seconds) to validate. By default, this is the current one when the action runs.

### Outputs

| Parameter   | Type    | Description                                                        |
|-------------|---------|--------------------------------------------------------------------|
| `is_valid`  | Boolean | If the timestamp is valid (i.e. is within an allowed time window). |
| `error`     | String  | If an error occurred, this parameter would describe it.            |
| `message`   | String  | A message explaining the result or the error.                      |
| `timestamp` | Integer | The timestamp (in seconds) which has been validated.               |
| `result`    | JSON    | A JSON representation of all the above parameters.                 |

### Examples

Some example jobs to give you an idea of how this action can be used:

```yaml
name: Example workflow
on:
  push:
jobs:
  simple_example:
    name: Simple example
    runs-on: ubuntu-latest
    steps:
      - name: Time Window Validator
        id: time_window_validator
        uses: github.com/landazuripaul/time-window-validator
        with:
          allowed: |
            - name: Office hours
              cronExpression: "0 9 * * 1-5"
              duration: 8h

      - name: Checkout
        # the time_window_validator fails on an invalid timestamp
        # you can use the outputs, to recover
        # if: steps.time_window_validator.outputs.is_valid == true
        uses: actions/checkout@v4

      # ...
```

## TODOs

- Tests!
    - in Go code
    - in CI
- Use the timezone: at the moment, the timezone isn't used and everything is in UTC.
