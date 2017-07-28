# x86db - x86 instruction database

## Installation

```bash
# If GOPATH is set correctly, you may remove "GOPATH=`pwd`"
# from each line.

GOPATH=`pwd` go get github.com/dlespiau/x86db
GOPATH=`pwd` go install github.com/dlespiau/x86db/cmd/x86db-gogen

# Add `pwd`/bin/x86db-gogen to your path or use it directly
# from `pwd`/bin.
```

## Usage and examples

```bash
# List all instructions:
./bin/x86db-gogen list
# => very long output

# List only SSE3 instructions:
./bin/x86db-gogen list --extension SSE3
# ADDSUBPD  xmmreg,xmmrm  rm  66 0f d0 /r  PRESCOTT,SSE3,SO
# ADDSUBPS  xmmreg,xmmrm  rm  f2 0f d0 /r  PRESCOTT,SSE3,SO
# HADDPD    xmmreg,xmmrm  rm  66 0f 7c /r  PRESCOTT,SSE3,SO
# HADDPS    xmmreg,xmmrm  rm  f2 0f 7c /r  PRESCOTT,SSE3,SO
# HSUBPD    xmmreg,xmmrm  rm  66 0f 7d /r  PRESCOTT,SSE3,SO
# HSUBPS    xmmreg,xmmrm  rm  f2 0f 7d /r  PRESCOTT,SSE3,SO
# LDDQU     xmmreg,mem    rm  f2 0f f0 /r  PRESCOTT,SSE3,SO
# MOVDDUP   xmmreg,xmmrm  rm  f2 0f 12 /r  PRESCOTT,SSE3
# MOVSHDUP  xmmreg,xmmrm  rm  f3 0f 16 /r  PRESCOTT,SSE3
# MOVSLDUP  xmmreg,xmmrm  rm  f3 0f 12 /r  PRESCOTT,SSE3

# List SS2 instructions that are not tested AND
# do not list instructions that can take MMX operands:
./bin/x86db-gogen list --extension SSE2 --not-tested --not-mmx
# CLFLUSH  mem            m   np 0f ae /7  WILLAMETTE,SSE2
# MOVSD    xmmreg,xmmreg  rm  f2 0f 10 /r  WILLAMETTE,SSE2
# MOVSD    xmmreg,xmmreg  mr  f2 0f 11 /r  WILLAMETTE,SSE2
# MOVSD    mem64,xmmreg   mr  f2 0f 11 /r  WILLAMETTE,SSE2
# MOVSD    xmmreg,mem64   rm  f2 0f 10 /r  WILLAMETTE,SSE2
```

```
# Help output x86db-gogen

Usage:

  x86db-gogen command [options]

List of commands:

  help    print this help
  list    list x86 instructions

Filtering options:

  -extension string
    	select instructions by extension
  -known
    	select instructions already known by the go assembler
  -not-known
    	select instructions not already known by the go assembler
  -not-mmx
    	do not select instructions taking MMX operands
  -not-tested
    	select instructions with no test case in the go assembler
  -tested
    	select instructions with test cases in the go assembler
```
