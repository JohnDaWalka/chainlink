# Griddle Devenv


## Getting started - Deploying Workflow DON on Kind

### Prerequisites
1. Create overrides-kind file in configs folder with the following content:

```toml
[infra]
type = "griddle-devenv"

[infra.devenv]
namespace = "devenv-local"

[infra.devenv.metadata]
account = ""
owner   = ""
contact = ""
project = ""
service = ""
```

2. Clone devenv repo and copy reusable baseline configs to configs dir

```shell
mkdir configs/griddle-devenv
cp -R $CODE_DIR/devenv/example/flux-native/deploy/config/base configs/griddle-devenv/
```

3. Initialize griddle devenv kind cluster
```shell
dctl init
```

### Deploying Workflow DON on Kind

```shell
CTF_CONFIGS=./configs/workflow-don-crib.toml,./configs/overrides-kind.toml go run main.go env start
```


### Todo:
- [x] Deploy Anvil chain in deployBlockchain step
- [ ] Configure Telepresence in bootstrap step (kind)
  - [ ] Quit existing telepresenece before start 
- [ ] Run init step to setup kind cluster, test after nuke
- [ ] Sync baseline configs using vendir, eliminate manual copy step
- [ ] Configure Workflow DON in deployDons step