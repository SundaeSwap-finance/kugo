# Kugo

[![Licence](https://img.shields.io/github/license/SundaeSwap-finance/kugo)](https://github.com/SundaeSwap-finance/kugo/blob/main/LICENSE)

Kugo is a golang client library for interacting with [Kupo](https://github.com/CardanoSolutions/kupo), a chain-indexer for Cardano.

## QuickStart

### Getting started

```console

$ go get github.com/SundaeSwap-finance/kugo

```

```golang

ctx := context.Background()
k := kugo.New(WithEndpoint("http://localhost:1442"))

matches, err := k.Matches(ctx,
    OnlyUnspent(),
    Address("addr1vyc29pvl2uyzqt8nwxrcxnf558ffm27u3d9calxn8tdudjgz4xq9p"),
)
if err != nil {
    return fmt.Errorf("failed to fetch matches: %w", err)
}

```

### How to use

We don't yet have our own documentation set up, but the code is fairly simple to follow.

We currently support three methods:
 - `.Matches(...)`: query UTXOs that match various supported criteria
 - `.Patterns(...)`: query what patterns are currently being indexed by Kupo
 - `.Checkpoints(...)`: query "checkpoints" (block/slot combinations) that kupo is aware of

For more information, please refer to the [Kupo documentation.](https://cardanosolutions.github.io/kupo)

## Contributing

Want to contribute? See [CONTRIBUTING.md](./CONTRIBUTING.md) to know how.