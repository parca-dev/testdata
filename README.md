[![Apache 2 License](https://img.shields.io/badge/license-Apache%202-blue.svg)](LICENSE)
[![Built with Devbox](https://jetpack.io/img/devbox/shield_moon.svg)](https://jetpack.io/devbox/docs/contributor-quickstart/)
[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
![CI](https://github.com/parca-dev/testdata/workflows/ci/badge.svg?branch=main)
![CI](https://github.com/parca-dev/testdata/actions/workflows/ci.yml/badge.svg)

This repo has some small custom binaries and vendored executables we use for testing some of Parca Agent's features. This lives in a separate repository as we we wanted to commit executables and unwind tables that can be very large.

## Structure

- `tables/normal` contains a textual representation of the unwind tables. the names follow the following format `"<producer>_<the executable name>.txt"`. `<producer>` is the name of the tool that produced the table, "ours" for our implementation.
- `tables/compact` as above, but using the compact unwind format we use in BPF.
- `tables/tables` similar to compact tables, but after any processing. These are the ones that would be loaded in the BPF maps.
- `jitdump` contains the jitdump files for certain executables. These are used to test the JIT dump parser.
- `out/` has the produced executables for our custom programs.
- `src/` contains the source for our custom programs.
- `vendored/` hosts prebuilt executables.


## Building

To build the executables, run `mage build:all` or `make build`. This will build the executables and place them in `out/`.

## Adding new tables

To add a new table, run `mage generate:all` or `make generate`. This will build the tables and place them in `tables/normal/`,`tables/compact/` and `final_tables/`.


# TODO(kakkoyun):
# dagger mod install github.com/tsirysndr/daggerverse/devbox@0dd781fc462998725cb8f8eef88158a98945e7a8
# dag.DevBox().Run(...)
