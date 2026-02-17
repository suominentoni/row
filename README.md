# row

Filter lines from stdin by line number. A simpler alternative to `sed -n` and `awk` for grabbing specific lines.

## Install

```
go install github.com/suominentoni/row@latest
```

Or with Homebrew:

```
brew tap suominentoni/tap
brew install row
```

Or download `.deb`/`.rpm` packages from [GitHub Releases](https://github.com/suominentoni/row/releases):

```sh
# Debian/Ubuntu
sudo dpkg -i row_*.deb

# Fedora/RHEL
sudo rpm -i row_*.rpm
```

## Usage

```
row [flags] <range>
```

### Range syntax

| Syntax | Description | Example |
|--------|-------------|---------|
| `N` | Single line | `row 5` |
| `N-M` | Inclusive range | `row 3-10` |
| `...M` | From first line to M | `row ...50` |
| `N...` | From N to end | `row 100...` |
| `A,B,C` | Multiple ranges | `row 1-3,7,20...` |

### Flags

| Flag | Description |
|------|-------------|
| `-h`, `--hide` | Invert filter: hide matching lines |
| `-s`, `--separator` | Print `---` between non-contiguous output segments |
| `-n`, `--number` | Show line numbers |
| `-v`, `--version` | Show version |
| `--help` | Show help |

## Examples

```sh
# Print lines 5 through 10
seq 1 100 | row 5-10

# Print the first 3 lines
seq 1 100 | row ...3

# Print from line 95 to end
seq 1 100 | row 95...

# Print just line 42
seq 1 100 | row 42

# Print everything except lines 3-5
seq 1 10 | row -h 3-5

# Show line numbers
seq 1 100 | row -n 5-8
# 5: 5
# 6: 6
# 7: 7
# 8: 8

# Show non-contiguous ranges with separator
seq 1 20 | row -s 3-5,10-12
# 3
# 4
# 5
# ---
# 10
# 11
# 12
```
