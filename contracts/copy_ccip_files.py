import os

x = '''"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/nonce_manager"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_remote"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_0_0/rmn_proxy_contract"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/price_registry"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/mock_rmn_contract"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/rmn_contract"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/usdc_token_pool"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/registry_module_owner_custom"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/burn_from_mint_token_pool"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/burn_mint_token_pool"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/burn_with_from_mint_token_pool"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/lock_release_token_pool"
"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/token_pool"'''

# split the string to take only the part after the version
# the two last parts are the version and the pkgName
# copy
# from: ../core/gethwrappers/ccip/generated/latest/{pkgName}/{pkgName}_zksync.go for the file to be copied
# to: ../core/gethwrappers/ccip/generated/{version}/{pkgName}/{pkgName}_zksync.go
lines = x.splitlines()
for line in lines:
    parts = line.strip('"').split("/")
    version = parts[-2]
    pkg_name = parts[-1]
    src = f"../core/gethwrappers/ccip/generated/latest/{pkg_name}/{pkg_name}_zksync.go"
    dest = f"../core/gethwrappers/ccip/generated/{version}/{pkg_name}/{pkg_name}_zksync.go"
    
    # Ensure the destination directory exists
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    
    # Copy the file
    if os.path.exists(src):
        with open(src, "r") as src_file:
            content = src_file.read()
        with open(dest, "w") as dest_file:
            dest_file.write(content)
        print(f"Copied {src} to {dest}")
    else:
        print(f"Source file {src} does not exist")
