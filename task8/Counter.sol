// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
  string public version;
  uint256 public count;

  constructor(string memory _version) {
    version = _version;
  }

  function Increment() public {
    count += 1;
  }
}