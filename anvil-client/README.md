## Troubleshoot CapabilitiesRegistry

- `node`
```
const { ethers } = require("ethers");
const provider = new ethers.JsonRpcProvider("http://localhost:8545");
const blockNumber = await provider.getBlockNumber();
console.log(blockNumber);

const abi = JSON.parse(`<copy from chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0/capabilities_registry.go#91>`);
const address = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";

const capReg = new ethers.Contract(
  address,
  abi,
  provider
);

var x = await capReg.getNodeOperators();
console.dir(x, { depth: null });

var x = await capReg.getNodes();
console.dir(x, { depth: null });

var x = await capReg.getCapabilities();
console.dir(x, { depth: null });

var hashedCapabilityId = await capReg.getHashedCapabilityId(
  "securemint-trigger",
  "1.0.0",
);
console.log("securemint-trigger hashedCapabilityId", hashedCapabilityId);

x = await capReg.isCapabilityDeprecated(hashedCapabilityId);
console.dir(x, { depth: null });

x = await capReg.getDONs();
console.dir(x, { depth: null });

x = await capReg.getCapabilityConfigs(2, hashedCapabilityId);
console.dir(x, { depth: null });
```

