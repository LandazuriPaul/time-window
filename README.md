# Time Window Validator

GitHub Action to validate if a timestamp (the current one by default) is inside an `allowed` or `blocked` time window.

This is to prevent a workflow from running at certain times (e.g. out of office hours, bank holidays, etc.).

## Usage

The action requires at least one of the `allowed` or `blocked` inputs. If none is provided, the action will fail.
You can provide `allowed` and `blocked` time windows in combination.

The timestamp is validated against both `allowed` and `blocked` time windows provided.

### Validation logic

**A given timestamp is considered as valid only if it respects both rules:**

- it is within an `allowed` time window,
- it is not within a `blocked` time window.

This means that the `blocked` windows always take precedence over the `allowed` ones.
This enables configurations where the blocked time windows correspond to bank holidays or non-business periods to
overtake usual business hours.

### Inputs

| Parameter            | Default value                                                                                | Description                                                                                                                     |
|----------------------|----------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|
| `allowed`            | `""`                                                                                         | The allowed time windows in YAML. See [`allowed` & `blocked`](#allowed--blocked) for more information.                          |
| `blocked`            | `""`                                                                                         | The blocked time windows in YAML. See [`allowed` & `blocked`](#allowed--blocked) for more information.                          |
| `commit-message`     | `${{ github.event.head_commit.message \|\| github.event.workflow_run.head_commit.message }}` | The commit message used to check the regular expression against. See [`commit-message`](#commit-message) for more information.  |
| `force-valid-regexp` | `force\-time\-window`                                                                        | A regular expression to check against the commit message. See [`force-valid-regepx`](#force-valid-regexp) for more information. |
| `timestamp`          | now                                                                                          | Unix timestamp (in seconds) to validate the time windows against. See [`timestamp`](#timestamp) for more information.           |

#### `allowed` & `blocked`

#### `commit-message`

#### `force-valid-regexp`

#### `timestamp`

### Outputs

| Parameter   | Type    | Description                                             |
|-------------|---------|---------------------------------------------------------|
| `is_valid`  | Boolean | If the timestamp is valid.                              |
| `error`     | String  | If an error occurred, this parameter would describe it. |
| `message`   | String  | A message explaining the result or the error.           |
| `timestamp` | Integer | The timestamp (in seconds) which has been validated.    |
| `result`    | JSON    | A JSON representation of all the above parameters.      |

### Examples

## TODOs

- Use the timezone: at the moment, the timezone isn't used and everything is in UTC.
