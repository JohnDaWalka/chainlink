# Secure Mint

If you're reading this you're looking at the parked/secure-mint-q3-2025 branch.

This branch contains the functionality to make Secure Mint run end-to-end on Staging. 

There are a few open TODOs:

- [ ] Make sure the plugin is included in the Docker image
- [ ]  Find a better way to verify reports in the secure mint integration test so we can remove the XXX_SingletonTransmitter hack
- [ ]  Double-check that these two new integration tests run properly in CI:
  - [ ] core/services/ocr3/securemint/integrationtest/integration_test.go
  - [ ] core/capabilities/integration_tests/keystone/securemint_workflow_test.go



