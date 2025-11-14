// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./NFTAuction.sol";

// 使用插件 hardhat-upgrades 进行合约可升级性管理时，只需要继承 Initializable 即可。
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

// 使用 OpenZeppelin 提供的 ERC1967Proxy 作为 UUPS 代理合约
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract NFTAuctionFactory is
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    // 存储创建的拍卖合约 - 这里使用嵌套映射，方便根据 NFT 地址和 Token ID 查找对应的拍卖合约，也可以用 NFTAuction 管理一组拍卖合约。当前拍卖合约支持多个拍卖。
    mapping(address NFT => mapping(uint256 tokenID => NFTAuction))
        public auctions;

    // 初始化函数
    function initialize(address owner_) public initializer {
        __Ownable_init(owner_); // 初始化 owner
        __UUPSUpgradeable_init(); // 初始化升级机制
    }

    // UUPS 升级权限函数
    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}

    function initAuctionProxy(address _nftAddr, uint256 _nftToken) external {
        // 创建实现合约 & proxy
        NFTAuction auction = new NFTAuction();
        bytes memory data = abi.encodeWithSelector(
            NFTAuction.initialize.selector,
            msg.sender,
            address(this)
        );
        ERC1967Proxy proxy = new ERC1967Proxy(address(auction), data);

        // 记录
        auctions[_nftAddr][_nftToken] = NFTAuction(address(proxy));
    }

    // 创建新的拍卖合约
    function createAuction(
        address _nftAddr,
        uint256 _nftToken,
        uint256 _minPrice,
        uint256 _duration
    ) external {
        // 获取 proxy 合约地址
        NFTAuction proxy = auctions[_nftAddr][_nftToken];
        require(address(proxy) != address(0), "Proxy not initialized");

        // 调用 proxy 的 createAuction 来初始化拍卖状态（proxy 现在已持有 NFT）
        NFTAuction(address(proxy)).createAuction(
            msg.sender,
            _nftAddr,
            _nftToken,
            _minPrice,
            _duration
        );
    }

    // 获取拍卖合约地址
    function getAuction(
        address _nftAddr,
        uint256 _nftToken
    ) external view returns (NFTAuction) {
        return auctions[_nftAddr][_nftToken];
    }

    function version() public pure returns (string memory) {
        return "v1";
    }

    uint256[50] private __gap; // 保留存储槽，用于未来升级
}
