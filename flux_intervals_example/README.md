# Flux Intervals Enforcement Example

This example demonstrates how to ensure consumers of your kubernetes cluster don't slam your cluster's API or git repositories by configuring [flux](https://fluxcd.io/) to sync too frequently.

Policies for enforcing flux argument values are in the [policies directory](./policies)

## Arguments checked

[Flux arguments documentation](https://docs.fluxcd.io/en/latest/references/daemon/)

- `--git-poll-interval`
  - Ensure at least set to `10m`
- `--sync-interval`
  - Ensure at least set to `10m`

## Test

```bash
make test
```

> Requires docker
