# TBMS
TBM BadAss Management System - Simple Graph Database


## Dev dependencies

```bash
# for debugging in vscode
go get -u github.com/derekparker/delve/cmd/dlv

```

## How to run

Clone project:

```bash

# move to go workspace
cd ~/go/src/github.com

# clone project to alikhil/TBMS dir
git clone git@github.com:alikhil/TBMS.git alikhil/TBMS

# get all dependencies, it also should clone alikhil/distributed-fs into go workspace
go get ./...

```

Prepare distributed file system(DFS):

**Note:** if you want to run DFS in several physical machines you need to manually clone repository on them. In other case, it's enough to run separate proccesses in different terminals

```bash
# go to DFS project dir
cd ~/go/src/github.com/alikhil/distributed-fs

# run master (terminal #1)
go run cmd/master/*.go -peers=3 # You can put another number here.
# after you run it. it will print endpoint for peer connection to stdout, copy it


# run peer 1 (terminal #2)
go run cmd/peer/*.go -fsdir=peer1 -port=5021 -endpoint={MASTER_ENDPOINT} # use endpoint from master log

# run peer 2 (terminal #3)
go run cmd/peer/*.go -fsdir=peer2 -port=5022 -endpoint={MASTER_ENDPOINT} # use endpoint from master log

# run peer 3 (terminal #4)
go run cmd/peer/*.go -fsdir=peer3 -port=5023 -endpoint={MASTER_ENDPOINT} # use endpoint from master log


```