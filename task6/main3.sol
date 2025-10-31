// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract {
    event Donation(address indexed from, uint amount);
    event Withdraw(uint amount);

    address public _owner;
    mapping(address => uint) private _all_donates;

    // 前三名数组
    address[3] public _topDonaters;

    uint256 public _startTime;
    uint256 public _endTime;

    modifier onlyOwner() {
        require(msg.sender == _owner, "Must owner can do.");
        _;
    }

    modifier withinTime() {
        require(
            block.timestamp >= _startTime && block.timestamp <= _endTime,
            "Donation not allowed at this time"
        );
        _;
    }

    constructor(uint256 start_time, uint256 end_time) {
        _owner = msg.sender;

        _startTime = start_time;
        _endTime = end_time;
    }

    function donate() public payable withinTime {
        require(msg.value > 0, "Donation must be greater than 0.");

        _all_donates[msg.sender] += msg.value;
        emit Donation(msg.sender, msg.value);

        // 计算排名
        uint cur_amount = _all_donates[msg.sender];
        // 如果当前捐赠金额大于前三名中的某一个，则更新前三名数组
        uint index = 3; // 默认为未入榜
        for (uint i = 0; i < 3; i++) {
            uint idx_amount = _all_donates[_topDonaters[i]];
            if (cur_amount > idx_amount) {
                index = i;
                break;
            }
        }

        // 该用户取得前三名。
        if (index < 3) {
            // 如果排名是 0 或 1 进行移动
            for (uint i = 2; i > index; i--) {
                _topDonaters[i] = _topDonaters[i - 1];
            }
            // 更新前三名
            _topDonaters[index] = msg.sender;
        }
    }

    function getDonation(address donater) public view returns (uint amount) {
        require(
            donater != address(0),
            "Donater address cannot be zero address."
        );

        return _all_donates[donater];
    }

    // payable 修饰符只在接收 ETH的函数上才需要。
    function withdraw() public onlyOwner {
        // require(msg.sender ==_owner, "Must owner can withdraw.");

        uint amount = address(this).balance;

        // 方式一
        payable(msg.sender).transfer(amount);

        // 方式二
        // (bool ret,) = (msg.sender).call{value : amount}(""); // send all balance to owner
        // require(ret, "Transfer failed.");

        emit Withdraw(amount);
    }

    function getTops3() public view returns (address[3] memory tops) {
        return _topDonaters;
    }
}
