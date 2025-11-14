module.exports = async ({getNamedAccounts, deployments}) => {
    console.log("============================================ Starting NFT deployment script...");
    // 1、获取账户
    const {deployer} = await getNamedAccounts();
    console.log("Deploying MyNFT with account:", deployer);

    // 2、部署 MyNFT 合约
    const my_nft_info = await deployments.deploy("MyNFT", {from: deployer, args: [], log: false});
    console.log("MyNFT deployed to:", my_nft_info.address);

    // 3、保存部署信息到 hardhat-deploy 管理的部署记录中
    await deployments.save("MyNFTData", {
        address: my_nft_info.address,
        abi: my_nft_info.abi
    });
};

module.exports.tags = ["MyNFT"];
