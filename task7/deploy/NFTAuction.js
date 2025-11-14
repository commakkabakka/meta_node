module.exports = async ({ethers, getNamedAccounts, deployments, upgrades}) => {
    console.log("============================================ Starting NFTAuction deployment script...");
    // 1、获取账户
    const {deployer} = await getNamedAccounts();
    console.log("Deploying NFTAuction with account:", deployer);

    // 2、部署 NFTAuction 合约
    const nft_auction_info = await deployments.deploy("NFTAuction", {from: deployer, args: [], log: false});
    console.log("NFTAuction deployed to:", nft_auction_info.address);

    // 3、保存部署信息到 hardhat-deploy 管理的部署记录中
    await deployments.save("NFTAuctionData", {
        address: nft_auction_info.address,
        abi: nft_auction_info.abi
    });
};

module.exports.tags = ["NFTAuction"];
