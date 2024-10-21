# Time Window Validator

GitHub Action to validate if a timestamp (now by default) is inside a time window.

This is useful to be able to prevent a workflow from running at certain times (e.g. out of office hours, bank holidays,
etc.).

## Usage

### Inputs

| Parameter            | Default value                                                                                | Description                                                                                                                     |
|----------------------|----------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|
| `allowed`            | `""`                                                                                         | The allowed time windows in YAML. See [`allowed` & `blocked`](#allowed--blocked) for more information.                          |
| `blocked`            | `""`                                                                                         | The blocked time windows in YAML. See [`allowed` & `blocked`](#allowed--blocked) for more information.                          |
| `commit-message`     | `${{ github.event.head_commit.message \|\| github.event.workflow_run.head_commit.message }}` | The commit message used to check the regular expression against. See [`commit-message`](#commit-message) for more information.  |
| `force-allow-regexp` | `force\-time\-window`                                                                        | A regular expression to check against the commit message. See [`force-allow-regepx`](#force-allow-regexp) for more information. |
| `timestamp`          | now                                                                                          | Unix timestamp (in seconds) to validate the time windows against. See [`timestamp`](#timestamp) for more information.           |

#### `allowed` & `blocked`

#### `commit-message`

#### `force-allow-regexp`

#### `timestamp`

### Outputs

## TODOs

- Use the timezone: at the moment, the timezone isn't used and everything is in UTC.
