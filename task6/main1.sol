// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MyERC20 {
    // MUST trigger when tokens are transferred, including zero value transfers.
    // A token contract which creates new tokens SHOULD trigger a Transfer event with the _from address set to 0x0 when tokens are created.
    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    // MUST trigger on any successful call to approve(address _spender, uint256 _value).
    event Approval(
        address indexed _owner,
        address indexed _spender,
        uint256 _value
    );

    // Returns the name of the token - e.g. "MyToken".
    string private _name = "My ERC20 Coin";

    // Returns the symbol of the token. E.g. “HIX”.
    string private _symbol = "MERC";

    // Returns the number of decimals the token uses - e.g. 8, means to divide the token amount by 100000000 to get its user representation.
    //     18：最常见，也是以太坊的标准（像 ETH 本身）。理由是给足够精度，方便微小交易，也符合以太坊上大多数代币标准。
    //     6：USDC、USDT 等稳定币常用，显示和美元精度一致（1 美元 = 1 单位，精确到小数点后 6 位）。
    //     8：比特币风格的代币有时用 8 位，方便兼容比特币用户习惯。
    //     0：极少数代币完全不拆分，比如某些 NFT 或只发行整数代币，直接把每个单位当作不可分割的“硬币”。
    uint8 private _decimals = 18;

    // Returns the total token supply.
    uint256 private _totalSupply = 10000 * (10 ** _decimals);

    // 合约拥有者
    address private _owner;
    // 保存用户余额
    mapping(address => uint256) private _balances;
    // 授权信息
    mapping(address => mapping(address => uint256)) _allowances;

    constructor() {
        _owner = msg.sender;

        // 初始铸币
        _balances[_owner] = _totalSupply;
        emit Transfer(address(0), _owner, _totalSupply);
    }

    modifier onlyOwner() {
        require(msg.sender == _owner, "Only owner can call this function.");
        _;
    }

    function name() public view returns (string memory) {
        return _name;
    }

    function symbol() public view returns (string memory) {
        return _symbol;
    }

    function decimals() public view returns (uint8) {
        return _decimals;
    }

    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address owner) public view returns (uint256 balance) {
        return _balances[owner];
    }

    /* 
        Transfers _value amount of tokens to address _to, and MUST fire the Transfer event. The function SHOULD throw if the message caller’s account balance does not have enough tokens to spend.
        Note Transfers of 0 values MUST be treated as normal transfers and fire the Transfer event.
    */
    function transfer(address to, uint256 value) public returns (bool success) {
        require(to != address(0), "Invalid address.");
        require(_balances[msg.sender] >= value, "Insufficient balance.");

        _balances[msg.sender] -= value;
        _balances[to] += value;
        emit Transfer(msg.sender, to, value);
        return true;
    }

    // Transfers _value amount of tokens from address _from to address _to, and MUST fire the Transfer event.
    function transferFrom(
        address from,
        address to,
        uint256 value
    ) public returns (bool success) {
        require(from != address(0) && to != address(0), "Invalid address.");
        require(
            _allowances[from][msg.sender] >= value,
            "Insufficient allowance."
        );
        require(_balances[from] >= value, "Insufficient balance.");

        _allowances[from][msg.sender] -= value;
        _balances[from] -= value;
        _balances[to] += value;
        emit Transfer(from, to, value);
        return true;
    }

    // Allows _spender to withdraw from your account multiple times, up to the _value amount. If this function is called again it overwrites the current allowance with _value.
    function approve(
        address spender,
        uint256 value
    ) public returns (bool success) {
        require(spender != address(0), "Invalid address.");

        _allowances[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    // Returns the amount which _spender is still allowed to withdraw from _owner.
    function allowance(
        address owner,
        address spender
    ) public view returns (uint256 remaining) {
        require(
            owner != address(0) && spender != address(0),
            "Invalid address."
        );

        return _allowances[owner][spender];
    }

    // 铸币 - 从账户 0 转账到 账户 account
    function mint(address account, uint256 value) internal onlyOwner {
        require(account != address(0), "Invalid address.");

        _balances[account] += value;
        _totalSupply += value;
        emit Transfer(address(0), account, value);
    }

    // 销毁 - 从账户 account 转账到 账户 0
    function burn(address account, uint256 value) internal onlyOwner {
        require(account != address(0), "Invalid address.");

        _balances[account] -= value;
        _totalSupply -= value;
        emit Transfer(account, address(0), value);
    }
}
