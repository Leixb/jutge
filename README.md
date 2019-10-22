# Jutge CLI
![GitHub](https://img.shields.io/github/license/Leixb/jutge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/Leixb/jutge)
[![Build Status](https://travis-ci.com/Leixb/jutge.svg?branch=master)](https://travis-ci.com/Leixb/jutge)
[![HitCount](http://hits.dwyl.io/Leixb/jutge.svg)](http://hits.dwyl.io/Leixb/jutge)
[![Go Report Card](https://goreportcard.com/badge/github.com/Leixb/jutge)](https://goreportcard.com/report/github.com/Leixb/jutge)
[![CodeFactor](https://www.codefactor.io/repository/github/leixb/jutge/badge)](https://www.codefactor.io/repository/github/leixb/jutge)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ab598b6be10640e5939cc4b197422acd)](https://www.codacy.com/manual/abone9999/jutge?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Leixb/jutge&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/Leixb/jutge?status.svg)](https://godoc.org/github.com/Leixb/jutge)

Easily create, test, upload and check problems from [jutge.org](https://jutge.org) without leaving the terminal.

## Install

```sh
go get https://github.com/Leixb/jutge
```

## Configuration

All configuration is done through with enviroment variables, mainly: `JUTGE_WORK_DIR`, `JUTGE_USER` and `JUTGE_PASSWORD` (Although the same options can be set with their respective flags: `--work-dir`, `--user`, `-password`.

 - `JUTGE_WORK_DIR` is the directory where you want all the problem files and data to be downloaded. 

 - `JUTGE_USER` (OPTIONAL): is the username (email address) to use when logging into jutge.
 - `JUTGE_PASSWORD` (OPTIONAL): is the password to use when logging into jutge. (!! I do not recommend to set this enviroment variable)

If no user or password are provided, the user will be prompted to enter them when needed.

It is **very important** to define `JUTGE_WORK_DIR` as an absolute path and add it to your `~/.bashrc` (or the equivalent for your shell).

 Example `~/.bashrc`:
```bash
export JUTGE_WORK_DIR="${HOME}/Documents/jutge/"
export JUTGE_USER="example@example.com"
```

If you want to use the `new` command you need to add the file `jutge.db` into your `$JUTGE_WORK_DIR` . You can download it directly from this repo with curl:

```bash
curl -o "${JUTGE_WORK_DIR}/jutge.db" https://raw.githubusercontent.com/Leixb/jutge/master/jutge.db
```

### Scripting
A very common task is to compile and then test the binary, or to compile test and then upload a file. This cannot be done directly with `jutge` but is fairly easy to do it with some basic shell scripting. For example, you can declare the following functions in your shell configuration file (`~/.bashrc`):

```bash

# Compile and test C++ program
jutgecpp() {
  name=$(basename -- "$1")
  out=$(mktemp /tmp/${name%.*}_XXXXX
  g++ "$1" -o "$out" -std=c++11 && jutge test "$out"
  rm $out
}

# Compile, test and upload if tests pass. Then wait for veredict and print it.
jutgeall() {
  jutgecpp "$1" && jutge upload "$@" --check
}
```

### Auto completion
You can add auto completion of commands by adding the following line to your shell configuration:

- If you use bash (`~/.bashrc`):
```bash
eval $(jutge --completion-script-bash)
```
- If you use zsh (`~/.zshrc`):
```zsh
eval $(jutge --completion-script-zsh)
```

## Commands

There are 6 commands:
 - Standard commands:
   - `new`: creates a new file for a problem (the filename contains the problem code followed by the problem title without accents or spaces e.g.: `P71753_ca_Maxim_de_cada_sequencia.cpp`)
   - `test`: tests an executable file against jutge samples (it will download the samples if needed)
   - `upload`: submits a problem to jutge
 - Rarely used commands:
   - `check`: check veredict of a submission
   - `download`: downloads sample test cases from jutge (usually not used since tests handles downloads when needed)
   - `db`: provies some sub commands to edit the the correspondence between codes and titles.
  
  If you want help for any of the commands just run `jutge command --help` to view all the options and their descriptions.

## Sample usage

[![asciicast](https://asciinema.org/a/dgHDQPtAqXyO5mGfDbJWzE8Ex.svg)](https://asciinema.org/a/dgHDQPtAqXyO5mGfDbJWzE8Ex)
