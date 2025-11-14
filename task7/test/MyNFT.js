const {ethers, deployments} = require("hardhat");
const {expect} = require("chai");

describe("--------- --------- --------- --------- --------- MyNFT Test", () => {
    let deployer, user1;

    let my_nft_data;
    let my_nft;
    let tokenId = 100;

    before(async () => {
        // 获取账户
        const accounts = await ethers.getSigners();
        deployer = accounts[0];
        user1 = accounts[1];

        // // 方法一：部署 MyNFT 合约
        // const MyNFT = await ethers.getContractFactory("MyNFT");
        // my_nft = await MyNFT.deploy();
        // await my_nft.waitForDeployment();
        // const my_nft_addr = await my_nft.getAddress();
        // expect(my_nft_addr).to.not.be.null;
        // console.log("MyNFT deployed at: ", my_nft_addr);

        // 方法二：部署 MyNFT 合约
        await deployments.fixture(["MyNFT"]);
        my_nft_data = await deployments.get("MyNFTData");
        expect(my_nft_data.address).to.not.be.null;
        my_nft = await ethers.getContractAt("MyNFT", my_nft_data.address);
        console.log("MyNFT deployed at: ", my_nft_data.address);

        // 创建 MyNFT
        const tx = await my_nft.mint(deployer.address, tokenId);
        await tx.wait();
    });

    it("查看 NFT 信息", async () => {
        const name = await my_nft.name();
        const symbol = await my_nft.symbol();
        expect(name).to.equal("HelloNFT");
        expect(symbol).to.equal("HNFT");
    });

    it("测试-转账", async () => {
        let balanceOfDeployer = await my_nft.balanceOf(deployer.address);
        expect(balanceOfDeployer).to.equal(1);

        // 转账给 user1
        const tx = await my_nft.safeTransferFrom(deployer.address, user1.address, tokenId);
        await tx.wait();

        balanceOfDeployer = await my_nft.balanceOf(deployer.address);
        expect(balanceOfDeployer).to.equal(0);

        let balanceOfUser1 = await my_nft.balanceOf(user1.address);
        expect(balanceOfUser1).to.equal(1);
    });
});
