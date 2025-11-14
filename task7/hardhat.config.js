require("@nomicfoundation/hardhat-toolbox");
require("hardhat-deploy");
require("@openzeppelin/hardhat-upgrades");

task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
    const accounts = await hre.ethers.getSigners();

    for (const account of accounts) {
        console.log(account.address);
    }
});

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
    solidity: "0.8.28",
    networks: {
        sepolia: {
            url: "https://sepolia.infura.io/v3/-",
            accounts: ["-", "-", "-", "-"]
        }
    },
    namedAccounts: {
        deployer: 0,
        user1: 1,
        user2: 2,
        user3: 3
    },
    mocha: {
        timeout: 200000 // 200s
    }
};
