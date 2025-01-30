// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_getTokenPool is CCTPMessageTransmitterProxySetup {
  function test_getTransmitter() public {
    assertEq(address(s_cctpMessageTransmitterProxy.i_transmitter()), address(s_cctpMessageTransmitter));
  }
}
