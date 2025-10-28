// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/*
    题目描述：创建一个名为Voting的合约，包含以下功能：
        一个mapping来存储候选人的得票数
        一个vote函数，允许用户投票给某个候选人
        一个getVotes函数，返回某个候选人的得票数
        一个resetVotes函数，重置所有候选人的得票数
*/

contract Voting {
    mapping(address => uint256) private _votes;
    address[] private _candidates;

    function vote(address candidate) public {
        require(candidate != address(0), "address is zero");

        // 票数为 0 说明后渲染还未添加到候选人数组中。
        if (_votes[candidate] == 0) {
            _candidates.push(candidate);
        }
        _votes[candidate] += 1;
    }

    function getVotes(address candidate) public view returns (uint256) {
        return _votes[candidate];
    }

    function resetVotes() public {
        // 清空 map
        for (uint256 i = 0; i < _candidates.length; ++i) {
            address key = _candidates[i];
            delete _votes[key];
        }

        // 清空 array
        delete _candidates;
    }
}
