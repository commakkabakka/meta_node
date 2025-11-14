module.exports = async ({ethers, getNamedAccounts, deployments, upgrades}) => {
    console.log("============================================ Starting NFTAuctionFactory proxy deployment script...");
    // 1、获取账户
    const {deployer} = await getNamedAccounts();
    console.log("Deploying NFTAuctionFactory with account:", deployer);

    // 2、用 OpenZeppelin 插件创建代理并部署实现合约
    // 获取 NFTAuctionFactory 合约工厂
    const NFTAuctionFactory = await ethers.getContractFactory("NFTAuctionFactory");
    // 创建代理并部署
    const nft_auction_factory_proxy = await upgrades.deployProxy(NFTAuctionFactory, [deployer], {initializer: "initialize"});
    await nft_auction_factory_proxy.waitForDeployment();
    // 获取代理合约地址
    const nft_auction_factory_proxy_address = await nft_auction_factory_proxy.getAddress();
    console.log("NFTAuctionFactory Proxy deployed to:", nft_auction_factory_proxy_address);
    // 获取实现合约地址
    const nft_auction_factory_implementation_address = await upgrades.erc1967.getImplementationAddress(nft_auction_factory_proxy_address);
    console.log("NFTAuctionFactory Implementation deployed to:", nft_auction_factory_implementation_address);

    // 3、保存部署信息到 hardhat-deploy 管理的部署记录中
    await deployments.save("NFTAuctionFactory_ProxyDeploy_Data", {
        address: nft_auction_factory_proxy_address,
        abi: nft_auction_factory_proxy.interface.format("json")
    });
};

module.exports.tags = ["NFTAuctionFactory_ProxyDeploy"];
