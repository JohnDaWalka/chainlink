// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../libraries/Internal.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";
import {IRMN} from "../../../interfaces/IRMN.sol";

contract RMNRemote_verify_withConfigSet is RMNRemoteSetup {
  event BlessedRootsGenerated(Internal.MerkleRoot[] blessedMerkleRoots);
  event DigestCalculated(bytes32 digest);

  error UnexpectedSigner();
  error OutOfOrderSignatures();
  error InvalidSignature();

  bytes32 private constant RMN_V1_6_ANY2EVM_REPORT = keccak256("RMN_V1_6_ANY2EVM_REPORT");

  function setUp() public override {
    super.setUp();

    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 3});
    s_rmnRemote.setConfig(config);
    _generatePayloadAndSigs(2, 4);
  }

  uint8 private constant ECDSA_RECOVERY_V = 27;

  function test_verify_real() public {
    RMNRemote rmnRemote = new RMNRemote(3478487238524512106, IRMN(address(0)));
    RMNRemote.Signer[] memory signers = new RMNRemote.Signer[](3);
    signers[0] = RMNRemote.Signer({onchainPublicKey: 0x996Fa60CE9A71Dd40F0E096b71938BFB079C82Ca, nodeIndex: 0});
    signers[1] = RMNRemote.Signer({onchainPublicKey: 0x858589216956F482A0f68B282a7050Af4cd48ed2, nodeIndex: 1});
    signers[2] = RMNRemote.Signer({onchainPublicKey: 0x7C5E94162C6FAbbdeb3BFE83Ae532846e337bfAE, nodeIndex: 2});
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: 0x000b25bf471d5b79368a0fac0f78940789eff42068d2014bb3c5560799bca728, signers: signers, fSign: 1});
    rmnRemote.setConfig(config);

    bytes memory rawReport = hex"0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000560000000000000000000000000000000000000000000000000000000000000058000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000008f90b8876dee6538000000000000000000000000000000b9a42e0000000000000000000070482e2b0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000048810ec3e431431f00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000014b900000000000000000000000000000000000000000000000000000000000015b86fc15777b4dff4afaaf0affe46dfb2f0ee01bad9bbadb716bb38748cdccd4b7f000000000000000000000000000000000000000000000000000000000000002000000000000000000000000020c8649fde48fc795154e1f9f5d9deef0719693c0000000000000000000000000000000000000000000000008f90b8876dee653800000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000c2d0000000000000000000000000000000000000000000000000000000000000ce785e086c8b83bba62564ac4b5f82100d8bfa9dcfe3a3d370abd7508952703de79000000000000000000000000000000000000000000000000000000000000002000000000000000000000000020c8649fde48fc795154e1f9f5d9deef0719693c000000000000000000000000000000000000000000000000de41ba4fc9d91ad900000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000085d000000000000000000000000000000000000000000000000000000000000088f509344edf33b0387d93d5f52dbed3d9dd5945b0bb9c73be5e70721af5f01b7830000000000000000000000000000000000000000000000000000000000000020000000000000000000000000024a6804c0afb97189aeffde5371e6835c57c6d3000000000000000000000000000000000000000000000000e1f4423f1bf587cd00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000097a00000000000000000000000000000000000000000000000000000000000009c58bd7eeff7a824801d2ce60b8f08d1fe4464eb7d12404c77336a51065868e4961000000000000000000000000000000000000000000000000000000000000002000000000000000000000000020c8649fde48fc795154e1f9f5d9deef0719693c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000027afd55aed1baa04d225e7ae28206a86696716348936426d45c06378b4661115d77a616c8fd6bfb1479e66269f63dc05f216de613b9275e28678ad43cfbfd0ddd0c3e8d02c0a6298eb7c0275b3c3f72fb35145f6affee2de40d407626f108688cbd71535dc14df2843b26665a817db34531eb2868cc9ecfc1a1a6a35d1bc1bd64";
    OffRamp.CommitReport memory commitReport = abi.decode(rawReport, (OffRamp.CommitReport));
    address offRamp = 0x9B73Aa7b0b92B8d8c79b9A85Ef9b810e73c1e901;

    emit BlessedRootsGenerated(commitReport.blessedMerkleRoots);

    bytes32 digest = keccak256(
      abi.encode(
        RMN_V1_6_ANY2EVM_REPORT,
        RMNRemote.Report({
          destChainId: 421614,
          destChainSelector: 3478487238524512106,
          rmnRemoteContractAddress: 0xD59F70d849876692D738D8c8185a5E9e3FBC7601,
          offrampAddress: offRamp,
          rmnHomeContractConfigDigest: 0x000b25bf471d5b79368a0fac0f78940789eff42068d2014bb3c5560799bca728,
          merkleRoots: commitReport.blessedMerkleRoots
        })
      )
    );

    emit DigestCalculated(digest);

    address prevAddress;
    address signerAddress;
    for (uint256 i = 0; i < commitReport.rmnSignatures.length; ++i) {
      signerAddress = ecrecover(digest, ECDSA_RECOVERY_V, commitReport.rmnSignatures[i].r, commitReport.rmnSignatures[i].s);
      if (signerAddress == address(0)) revert InvalidSignature();
      if (prevAddress >= signerAddress) revert OutOfOrderSignatures();

      bool found = false;
      for (uint256 j = 0; j < s_signers.length; ++j) {
        if (s_signers[j].onchainPublicKey == signerAddress) {
          found = true;
          break;
        }
      }
      if (!found) revert UnexpectedSigner();
      prevAddress = signerAddress;
    }
  }

  function test_verify() public view {
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_InvalidSignature() public {
    s_signatures[s_signatures.length - 1].r = 0x0;

    vm.expectRevert(RMNRemote.InvalidSignature.selector);

    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_OutOfOrderSignatures_not_sorted() public {
    IRMNRemote.Signature memory sig1 = s_signatures[s_signatures.length - 1];
    IRMNRemote.Signature memory sig2 = s_signatures[s_signatures.length - 2];

    s_signatures[s_signatures.length - 1] = sig2;
    s_signatures[s_signatures.length - 2] = sig1;

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_OutOfOrderSignatures_duplicateSignature() public {
    s_signatures[s_signatures.length - 1] = s_signatures[s_signatures.length - 2];

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_UnexpectedSigner() public {
    _setupSigners(4); // create new signers that aren't configured on RMNRemote
    _generatePayloadAndSigs(2, 4);

    vm.expectRevert(RMNRemote.UnexpectedSigner.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_ThresholdNotMet() public {
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 2}); // 3 = f+1 sigs required
    s_rmnRemote.setConfig(config);

    _generatePayloadAndSigs(2, 2); // 2 sigs generated, but 3 required

    vm.expectRevert(RMNRemote.ThresholdNotMet.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }
}
