const {ethers, deployments} = require("hardhat");
const {expect} = require("chai");

describe("--------- --------- --------- --------- --------- NFTAuction Test", () => {
    let deployer, user1, user2;

    let nft_auction_data;
    let nft_auction;

    let my_nft_data;
    let my_nft;
    let tokenId = 100;

    before(async () => {
        // 获取账户
        const accounts = await ethers.getSigners();
        deployer = accounts[0];
        user1 = accounts[1];
        user2 = accounts[2];

        // 部署 NFTAuction 合约
        await deployments.fixture(["NFTAuction"]);
        nft_auction_data = await deployments.get("NFTAuctionData");
        expect(nft_auction_data.address).to.not.be.null;
        nft_auction = await ethers.getContractAt("NFTAuction", nft_auction_data.address);
        console.log("NFTAuction deployed at: ", nft_auction_data.address);

        // 部署 MyNFT 合约
        await deployments.fixture(["MyNFT"]);
        my_nft_data = await deployments.get("MyNFTData");
        expect(my_nft_data.address).to.not.be.null;
        my_nft = await ethers.getContractAt("MyNFT", my_nft_data.address);
        console.log("MyNFT deployed at: ", my_nft_data.address);

        // 创建 MyNFT
        const tx = await my_nft.mint(deployer.address, tokenId);
        await tx.wait();
    });

    it("初始化", async () => {
        // 将 factory_ 参数设置为 deployer 方便测试。
        const tx = await nft_auction.initialize(deployer.address, deployer.address);
        await tx.wait();
    });

    it("创建拍卖", async () => {
        // 创建拍卖
        const tx1 = await my_nft.approve(nft_auction_data.address, tokenId);
        await tx1.wait();
        const tx2 = await nft_auction.createAuction(deployer.address, my_nft_data.address, tokenId, 100000000, 60);
        await tx2.wait();

        // 查看拍卖信息
        const auction = await nft_auction.auctions(my_nft_data.address, tokenId);
        expect(auction.seller).to.equal(deployer.address);
        expect(auction.nftAddr).to.equal(my_nft_data.address);
        expect(auction.nftToken).to.equal(tokenId);
        expect(auction.minPrice).to.equal(100000000);
    });

    it("出价", async () => {
        // user1 出价
        const erc20 = await ethers.getContractAt("IERC20", "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238");
        const tx1 = await erc20.connect(user1).approve(nft_auction_data.address, 100);
        await tx1.wait();
        const tx2 = await nft_auction.connect(user1).placeBid(my_nft_data.address, tokenId, "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238", 2, {value: 0});
        await tx2.wait();

        // 查看拍卖信息
        const auction = await nft_auction.auctions(my_nft_data.address, tokenId);
        expect(auction.highestBidder).to.equal(user1.address);
        expect(auction.highestBidAmount).to.equal(2);

        // user2 出价
        const tx3 = await nft_auction.connect(user2).placeBid(my_nft_data.address, tokenId, "0x0000000000000000000000000000000000000000", 0, {value: 1000000000000000});
        await tx3.wait();

        // 查看拍卖信息
        const auction2 = await nft_auction.auctions(my_nft_data.address, tokenId);
        expect(auction2.highestBidder).to.equal(user2.address);
        expect(auction2.highestBidAmount).to.equal(1000000000000000);
    });

    it("结束拍卖", async () => {
        // 增加时间，结束拍卖
        await new Promise((resolve) => setTimeout(resolve, 60000));

        // 结束拍卖
        const tx = await nft_auction.settleAuction(my_nft_data.address, tokenId);
        await tx.wait();

        // 查看 NFT 归属
        const owner = await my_nft.ownerOf(tokenId);
        expect(owner).to.equal(user2.address);
    });
});
