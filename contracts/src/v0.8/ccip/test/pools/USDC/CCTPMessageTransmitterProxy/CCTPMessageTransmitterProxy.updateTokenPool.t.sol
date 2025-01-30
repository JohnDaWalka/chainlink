// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../../shared/access/Ownable2Step.sol";
import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_updateTokenPool is CCTPMessageTransmitterProxySetup {
  function test_updateTokenPool() public {
    s_cctpMessageTransmitterProxy.updateTokenPool(s_usdcTokenPool);
    assertEq(s_cctpMessageTransmitterProxy.s_tokenPool(), s_usdcTokenPool);
  }

  // Revert cases
  function test_updateTokenPool_RevertWhen_NotOwner() public {
    changePrank(makeAddr("RANDOM"));
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_cctpMessageTransmitterProxy.updateTokenPool(s_usdcTokenPool);
  }
}
