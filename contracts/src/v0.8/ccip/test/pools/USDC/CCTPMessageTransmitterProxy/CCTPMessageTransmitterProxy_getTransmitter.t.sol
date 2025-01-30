// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_getTransmitter is CCTPMessageTransmitterProxySetup {
  function test_getTransmitter() public {
    assertEq(address(s_cctpMessageTransmitterProxy.getTransmitter()), address(s_cctpMessageTransmitter));
  }
}
