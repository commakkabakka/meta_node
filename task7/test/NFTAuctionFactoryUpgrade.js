const {ethers, deployments, upgrades} = require("hardhat");
const {expect} = require("chai");

describe("--------- --------- --------- --------- --------- NFTAuctionfactory Upgrade Test", () => {
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
        await deployments.fixture(["NFTAuctionFactory_ProxyDeploy"], {keepExistingDeployments: false});
        nft_auction_factory_data = await deployments.get("NFTAuctionFactory_ProxyDeploy_Data");
        expect(nft_auction_factory_data.address).to.not.be.null;
        nft_auction_factory = await ethers.getContractAt("NFTAuctionFactory", nft_auction_factory_data.address);

        // 部署 MyNFT 合约
        await deployments.fixture(["MyNFT"]);
        my_nft_data = await deployments.get("MyNFTData");
        expect(my_nft_data.address).to.not.be.null;
        my_nft = await ethers.getContractAt("MyNFT", my_nft_data.address);

        // 创建 MyNFT
        const tx = await my_nft.mint(deployer.address, tokenId);
        await tx.wait();
    });

    // 被自动调用了。
    // it("初始化", async () => {
    //     // 初始化 NFTAuctionFactory 合约
    //     const tx = await nft_auction_factory.initialize(deployer.address);
    //     await tx.wait();
    // });

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

    it("升级 NFTAuctionFactory 合约到 V2", async () => {
        // 获取代理合约地址
        const nft_auction_factory_proxy_address = nft_auction_factory_data.address;
        console.log("NFTAuctionFactory Proxy address:", nft_auction_factory_proxy_address);

        // 获取实现合约地址
        const nft_auction_factory_implementation_address_v1 = await upgrades.erc1967.getImplementationAddress(nft_auction_factory_proxy_address);
        console.log("NFTAuctionFactoryV1 Implementation address:", nft_auction_factory_implementation_address_v1);
        // const implV1 = await upgrades.erc1967.getImplementationAddress(nft_auction_factory_proxy_address);
        // console.log("V1 code:", await ethers.provider.getCode(implV1));

        // 获取拍卖数据。
        const auction1 = await nft_auction_factory.auctions(my_nft_data.address, tokenId);

        // 升级合约
        // await deployments.fixture(["NFTAuctionFactory_ProxyDeploy", "NFTAuctionFactory_ProxyUpgrade"]);
        const NFTAuctionFactoryV2 = await ethers.getContractFactory("NFTAuctionFactoryV2");
        const nft_auction_factory_proxy_v2 = await upgrades.upgradeProxy(nft_auction_factory_proxy_address, NFTAuctionFactoryV2);
        await nft_auction_factory_proxy_v2.waitForDeployment();

        // 获取实现合约地址
        const nft_auction_factory_implementation_address_v2 = await upgrades.erc1967.getImplementationAddress(nft_auction_factory_proxy_address);
        console.log("NFTAuctionFactoryV2 Implementation address:", nft_auction_factory_implementation_address_v2);
        // const implV2 = await upgrades.erc1967.getImplementationAddress(nft_auction_factory_proxy_address);
        // console.log("V2 code:", await ethers.provider.getCode(implV2));

        // 获取拍卖数据。
        const auction2 = await nft_auction_factory.auctions(my_nft_data.address, tokenId);

        // 验证实现合约地址已更改
        expect(nft_auction_factory_implementation_address_v1).to.not.equal(nft_auction_factory_implementation_address_v2);

        // 5、【验证】验证拍卖数据未变
        expect(auction1.seller).to.equal(auction2.seller);
        expect(auction1.nftAddress).to.equal(auction2.nftAddress);
        expect(auction1.tokenId).to.equal(auction2.tokenId);
        expect(auction1.startingPrice).to.equal(auction2.startingPrice);
        expect(auction1.highestBid).to.equal(auction2.highestBid);
        expect(auction1.endTime).to.equal(auction2.endTime);
    });
});
