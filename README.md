This repo has some small custom binaries and vendored executables we use for testing some of Parca Agent's features. This lives in a separate repository as we we wanted to commit executables and unwind tables that can be very large.

## Structure

- `src/` contains the source for our custom programs
- `out/` has the produced executables for our custom programs
- `vendor/` hosts prebuilt executables
- `tables/` contains a textual representation of the unwind tables. The names follow the following format `"<producer>_<the executable name>.txt"`. `<producer>` is the name of the tool that produced the table, "ours" for our implementation 


## TODO
- Improve Makefile
- Don't use the system compiler / tools
