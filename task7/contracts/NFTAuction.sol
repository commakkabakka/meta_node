// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 引入 OpenZeppelin 合约库
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

// UUPS
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

// 使用 Chainlink 价格预言机
import {
    AggregatorV3Interface
} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

// NFT 拍卖合约
contract NFTAuction is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    event AuctionCreated(
        address indexed nftAddr,
        uint256 indexed nftToken,
        address seller,
        uint256 minPrice, // minPrice 使用 USD 计价,包含 8 位小数。
        uint256 endTime
    );

    event AuctionBided(
        address indexed nftAddress,
        uint256 indexed tokenId,
        address bidder,
        uint256 value // value 使用 USD 计价,包含 8 位小数。
    );

    event AuctionSettled(
        address indexed nftAddress,
        uint256 indexed tokenId,
        address winner,
        uint256 value // value 使用 USD 计价,包含 8 位小数。
    );

    // 拍卖结构体
    struct Auction {
        bool settled; // 拍卖是否已结算
        address seller; // 拍卖卖家
        address nftAddr; // NFT 合约地址
        uint256 nftToken; // NFT ID
        uint256 minPrice; // 最低竞拍价 - 使用 USD 计价,包含 8 位小数。
        uint256 endTime; // 拍卖结束时间 - 单位：秒
        address highestBidder; // 最高竞拍者
        address highestBidAddr; // 最高竞拍者 - 出价 ERC20 代币 合约地址
        uint256 highestBidAmount; // 最高竞拍者 - 出价 ERC20 代币数量
    }
    mapping(address => mapping(uint256 => Auction)) public auctions;

    // 存储每个代币地址对应的 Chainlink 数据源合约地址 Feed Address
    mapping(address => AggregatorV3Interface) public dataFeeds;

    // 只能由 factory 调用
    address public factory;
    modifier onlyFactory() {
        require(msg.sender == factory, "Only factory can call");
        _;
    }

    // UUPS
    function initialize(address owner_, address factory_) public initializer {
        __Ownable_init(owner_); // 初始化 owner
        __UUPSUpgradeable_init(); // 初始化升级机制
        factory = factory_;

        // 默认支持 ETH 和 USDC 作为竞拍货币
        // ETH=>USD
        setDataFeed(address(0), 0x694AA1769357215DE4FAC081bf1f309aDC325306);
        // USDC=>USD
        setDataFeed(
            0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238,
            0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E
        );
    }

    // UUPS
    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}

    function setDataFeed(address bidAddr, address dataFeed) public {
        dataFeeds[bidAddr] = AggregatorV3Interface(dataFeed);
    }

    function getChainlinkDataFeedLatestAnswer(
        address bidAddr
    ) public view returns (int256) {
        AggregatorV3Interface dataFeed = dataFeeds[bidAddr];
        (, int256 answer, , , ) = dataFeed.latestRoundData();
        return answer;
    }

    // 创建拍卖
    function createAuction(
        address _seller,
        address _nftAddr,
        uint256 _nftToken,
        uint256 _minPrice,
        uint256 _duration
    ) external onlyFactory {
        // 验证
        require(
            auctions[_nftAddr][_nftToken].seller == address(0),
            "Auction already exists"
        );

        // 将 NFT 从卖家转入 proxy（Factory 来执行转移）
        IERC721(_nftAddr).transferFrom(_seller, address(this), _nftToken);

        // 创建拍卖
        auctions[_nftAddr][_nftToken] = Auction({
            settled: false,
            seller: msg.sender,
            nftAddr: _nftAddr,
            nftToken: _nftToken,
            minPrice: _minPrice,
            endTime: block.timestamp + _duration,
            highestBidder: address(0),
            highestBidAddr: address(0),
            highestBidAmount: 0
        });

        // 触发拍卖创建事件
        emit AuctionCreated(
            _nftAddr,
            _nftToken,
            msg.sender,
            _minPrice,
            block.timestamp + _duration
        );
    }

    function placeBid(
        address _nftAddr,
        uint256 _nftToken,
        address _bidAddr,
        uint256 _amount
    ) external payable {
        // 获取拍卖信息
        Auction storage auction = auctions[_nftAddr][_nftToken];
        require(
            auctions[_nftAddr][_nftToken].seller != address(0),
            "Auction does not exist"
        );
        require(
            !auction.settled && block.timestamp < auction.endTime,
            "Auction ended"
        );

        // 根绝出价类型进行处理
        if (_bidAddr == address(0)) {
            // 使用 ETH 出价

            // 计算出价的 USD 价值
            int256 ethToUsd = getChainlinkDataFeedLatestAnswer(address(0));
            uint256 bidValueInUSD = (uint256(ethToUsd) * msg.value) / 1e18;

            // 计算之前最高出价的 USD 价值
            int256 oldEthToUsd;
            uint256 highestBidValueInUSD;
            if (auction.highestBidAddr == address(0)) {
                oldEthToUsd = ethToUsd;
                highestBidValueInUSD =
                    (uint256(oldEthToUsd) * auction.highestBidAmount) /
                    1e18;
            } else {
                oldEthToUsd = getChainlinkDataFeedLatestAnswer(
                    auction.highestBidAddr
                );
                highestBidValueInUSD =
                    uint256(oldEthToUsd) *
                    auction.highestBidAmount;
            }

            // 验证出价有效性
            require(
                bidValueInUSD >= auction.minPrice &&
                    bidValueInUSD > highestBidValueInUSD,
                "Bid too low"
            );

            // 使用 ETH 出价 - 已通过 msg.value 传入
            // ...

            // 退还之前的最高出价者
            if (auction.highestBidder != address(0)) {
                if (auction.highestBidAddr == address(0)) {
                    // 之前最高价是 ETH
                    payable(auction.highestBidder).transfer(
                        auction.highestBidAmount
                    );
                } else {
                    // 之前最高价是 ERC20 代币
                    IERC20 oldErc20 = IERC20(auction.highestBidAddr);
                    oldErc20.transfer(
                        auction.highestBidder,
                        auction.highestBidAmount
                    );
                }
            }

            // 更新最高出价信息
            auction.highestBidder = msg.sender;
            auction.highestBidAddr = _bidAddr;
            auction.highestBidAmount = msg.value;

            // 触发出价事件
            emit AuctionBided(_nftAddr, _nftToken, msg.sender, bidValueInUSD);
        } else {
            // 使用 ERC20 代币出价

            // 计算出价的 USD 价值
            int256 tokenToUsd = getChainlinkDataFeedLatestAnswer(_bidAddr);
            uint256 bidValueInUSD = uint256(tokenToUsd) * _amount;

            // 计算之前最高出价的 USD 价值
            int256 oldTokenToUsd;
            uint256 highestBidValueInUSD;
            if (auction.highestBidAddr == address(0)) {
                oldTokenToUsd = getChainlinkDataFeedLatestAnswer(address(0));
                highestBidValueInUSD =
                    (uint256(oldTokenToUsd) * auction.highestBidAmount) /
                    1e18;
            } else {
                oldTokenToUsd = getChainlinkDataFeedLatestAnswer(
                    auction.highestBidAddr
                );
                highestBidValueInUSD =
                    uint256(oldTokenToUsd) *
                    auction.highestBidAmount;
            }

            // 验证出价有效性
            require(
                bidValueInUSD >= auction.minPrice &&
                    bidValueInUSD > highestBidValueInUSD,
                "Bid too low"
            );

            // 使用 ERC20 代币出价
            IERC20 erc20 = IERC20(_bidAddr);
            erc20.transferFrom(msg.sender, address(this), _amount);

            // 退还之前的最高出价者
            if (auction.highestBidder != address(0)) {
                if (auction.highestBidAddr == address(0)) {
                    // 之前最高价是 ETH
                    payable(auction.highestBidder).transfer(
                        auction.highestBidAmount
                    );
                } else {
                    // 之前最高价是 ERC20 代币
                    IERC20 oldErc20 = IERC20(auction.highestBidAddr);
                    oldErc20.transfer(
                        auction.highestBidder,
                        auction.highestBidAmount
                    );
                }
            }

            // 更新最高出价信息
            auction.highestBidder = msg.sender;
            auction.highestBidAddr = _bidAddr;
            auction.highestBidAmount = _amount;

            // 触发出价事件
            emit AuctionBided(_nftAddr, _nftToken, msg.sender, bidValueInUSD);
        }
    }

    function settleAuction(address _nftAddr, uint256 _nftToken) external {
        // 获取拍卖信息
        Auction storage auction = auctions[_nftAddr][_nftToken];
        require(
            auctions[_nftAddr][_nftToken].seller != address(0),
            "Auction does not exist"
        );
        require(block.timestamp >= auction.endTime, "Auction not yet ended");
        require(!auction.settled, "Auction already settled");
        auction.settled = true;

        // 处理拍卖结算
        IERC721 nft = IERC721(_nftAddr);
        // 有人竞拍成功
        if (auction.highestBidder != address(0)) {
            // 将 NFT 转移给最高出价者
            nft.transferFrom(address(this), auction.highestBidder, _nftToken);
            // 将资金转移给卖家
            if (auction.highestBidAddr == address(0)) {
                // ETH
                payable(auction.seller).transfer(auction.highestBidAmount);
            } else {
                // ERC20 代币
                IERC20 erc20 = IERC20(auction.highestBidAddr);
                erc20.transfer(auction.seller, auction.highestBidAmount);
            }

            // 计算最高出价的 USD 价值
            uint256 bidValueInUSD;
            if (auction.highestBidAddr == address(0)) {
                // ETH 出价

                int256 ethToUsd = getChainlinkDataFeedLatestAnswer(address(0));
                bidValueInUSD =
                    (uint256(ethToUsd) * auction.highestBidAmount) /
                    1e18;
            } else {
                // ERC20 代币出价
                int256 tokenToUsd = getChainlinkDataFeedLatestAnswer(
                    auction.highestBidAddr
                );
                bidValueInUSD = uint256(tokenToUsd) * auction.highestBidAmount;
            }
            // 触发拍卖结算事件
            emit AuctionSettled(
                _nftAddr,
                _nftToken,
                auction.highestBidder,
                bidValueInUSD
            );
        } else {
            // 无人竞拍成功，退还 NFT 给卖家
            nft.transferFrom(address(this), auction.seller, _nftToken);

            // 触发拍卖结算事件
            emit AuctionSettled(_nftAddr, _nftToken, address(0), 0);
        }
    }

    uint256[50] private __gap; // 保留存储槽，用于未来升级
}
