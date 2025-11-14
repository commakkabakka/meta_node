module.exports = async ({ethers, getNamedAccounts, deployments, upgrades}) => {
    console.log("============================================ Starting NFTAuctionFactory deployment script...");
    // 1、获取账户
    const {deployer} = await getNamedAccounts();
    console.log("Deploying NFTAuctionFactory with account:", deployer);

    // 2、部署 NFTAuctionFactory 合约
    const nft_auction_factory_info = await deployments.deploy("NFTAuctionFactory", {from: deployer, args: [], log: false});
    console.log("NFTAuctionFactory deployed to:", nft_auction_factory_info.address);

    // 3、保存部署信息到 hardhat-deploy 管理的部署记录中
    await deployments.save("NFTAuctionFactoryData", {
        address: nft_auction_factory_info.address,
        abi: nft_auction_factory_info.abi
    });
};

module.exports.tags = ["NFTAuctionFactory"];
