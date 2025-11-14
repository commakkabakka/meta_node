module.exports = async ({ethers, getNamedAccounts, deployments, upgrades}) => {
    console.log("============================================ Starting NFTAuctionFactory proxy upgrade script...");
    // 1、获取账户
    const {deployer} = await getNamedAccounts();
    console.log("Upgrading NFTAuctionFactory with account:", deployer);

    // 2、用 OpenZeppelin 插件升级 NFTAuctionFactory 合约
    // 获取 NFTAuctionFactoryV2 合约工厂
    const NFTAuctionFactoryV2 = await ethers.getContractFactory("NFTAuctionFactoryV2");
    // 获取已部署的 NFTAuctionFactory 代理合约地址
    const nft_auction_factory_data = await deployments.get("NFTAuctionFactory_ProxyDeploy_Data");
    const nft_auction_factory_proxy_address = nft_auction_factory_data.address;
    // 升级代理合约到新的实现合约
    const nft_auction_factory_proxy_v2 = await upgrades.upgradeProxy(nft_auction_factory_proxy_address, NFTAuctionFactoryV2);
    await nft_auction_factory_proxy_v2.waitForDeployment();
    // 获取代理合约地址
    const nft_auction_factory_proxy_address_v2 = await nft_auction_factory_proxy_v2.getAddress();
    console.log("NFTAuctionFactoryV2 Proxy upgraded to:", nft_auction_factory_proxy_address_v2);
    // 获取新的实现合约地址
    const nft_auction_factory_implementation_address = await upgrades.erc1967.getImplementationAddress(nft_auction_factory_proxy_address_v2);
    console.log("NFTAuctionFactoryV2 Implementation upgraded to:", nft_auction_factory_implementation_address);

    // 3、保存部署信息到 hardhat-deploy 管理的部署记录中
    await deployments.save("NFTAuctionFactory_ProxyUpgrade_DataV2", {
        address: nft_auction_factory_proxy_address_v2,
        abi: nft_auction_factory_proxy_v2.interface.format("json")
    });
};

module.exports.tags = ["NFTAuctionFactory_ProxyUpgrade"];
