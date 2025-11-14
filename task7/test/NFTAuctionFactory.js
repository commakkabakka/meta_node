const {ethers} = require("hardhat");
const {expect} = require("chai");

describe("--------- --------- --------- --------- --------- NFTAuctionfactory Test", async () => {
    let deployer, user1, user2;

    let nft_auction_factory_data;
    let nft_auction_factory;

    let my_nft_data;
    let my_nft;
    let tokenId = 100;

    before(async () => {
        // 获取账户
        const accounts = await ethers.getSigners();
        deployer = accounts[0];
        user1 = accounts[1];
        user2 = accounts[2];

        // 部署 NFTAuctionFactory 合约
        await deployments.fixture(["NFTAuctionFactory"]);
        nft_auction_factory_data = await deployments.get("NFTAuctionFactoryData");
        expect(nft_auction_factory_data.address).to.not.be.null;
        nft_auction_factory = await ethers.getContractAt("NFTAuctionFactory", nft_auction_factory_data.address);
        console.log("NFTAuctionFactory deployed at: ", nft_auction_factory_data.address);

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
        // 初始化 NFTAuctionFactory 合约
        const tx = await nft_auction_factory.initialize(deployer.address);
        await tx.wait();
    });

    it("创建拍卖", async () => {
        // 创建 NFTAuction 合约 proxy
        const tx1 = await nft_auction_factory.initAuctionProxy(my_nft_data.address, tokenId);
        await tx1.wait();

        // 获取拍卖合约地址
        const auction = await nft_auction_factory.getAuction(my_nft_data.address, tokenId);
        expect(auction).to.not.be.null;

        // 授权
        const tx2 = await my_nft.approve(auction, tokenId);
        await tx2.wait();

        // 创建拍卖
        const tx3 = await nft_auction_factory.createAuction(my_nft_data.address, tokenId, 100000000, 60);
        await tx3.wait();
    });
});
