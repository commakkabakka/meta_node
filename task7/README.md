一、环境

-   Ubuntu 24.04.3 LTS
-   nodejs v24.11.0
-   hardhat 2.27.0

二、文件

    ├── README.md
    ├── contracts
    │   ├── MyNFT.sol				// ERC721 合约
    │   ├── NFTAuction.sol			// 拍卖合约
    │   ├── NFTAuctionFactory.sol	 // 拍卖工厂合约
    │   └── NFTAuctionFactoryV2.sol	 // 拍卖工厂合约V2 - 测试升级使用
    ├── deploy
    │   ├── MyNFT.js							// 部署 MyNFT 脚本
    │   ├── NFTAuction.js                          // 部署 NFTAuction 脚本
    │   ├── NFTAuctionFactory.js				 // 部署 NFTAuctionFactory 脚本
    │   ├── NFTAuctionFactory_ProxyDeploy.js	  // 使用 hardhat-upgrades 部署 NFTAuctionFactory 脚本
    │   └── NFTAuctionFactory_ProxyUpgrade.js	  // 使用 hardhat-upgrades 升级 NFTAuctionFactory 脚本到 NFTAuctionFactoryV2
    ├── hardhat.config.js
    ├── package-lock.json
    ├── package.json
    └── test
        ├── MyNFT.js					// 测试 MyNFT
        ├── NFTAuction.js				// 测试 NFTAuction
        ├── NFTAuctionFactory.js		 // 测试 NFTAuctionFactory
        └── NFTAuctionFactoryUpgrade.js	 // 测试 NFTAuctionFactory 升级

三、安装

1、安装依赖

    npm install

2、配置账户

编辑 hardhat.config.js, 在 accounts 中添加自己测试账户私钥。

四、命令

    # 1、部署 MyNFT 合约
    npx hardhat deploy --tags MyNFT --network localhost
    # 2、部署 NFTAuction 合约
    npx hardhat deploy --tags NFTAuction --network localhost
    # 3、部署 NFTAuctionFactory 合约
    npx hardhat deploy --tags NFTAuctionFactory --network localhost
    # 4、使用 hardhat-upgrades 部署 NFTAuctionFactory 合约
    npx hardhat deploy --tags NFTAuctionFactory_ProxyDeploy --network localhost
    # 5、使用 hardhat-upgrades 升级 NFTAuctionFactory 合约
    npx hardhat deploy --tags NFTAuctionFactory_ProxyUpgrade --network localhost

    # 6、测试 MyNFT 功能
    npx hardhat test test/MyNFT.js --network localhost
    # 7、测试 NFTAuction 功能
    npx hardhat test test/NFTAuction.js --network sepolia
    # 8、测试 NFTAuctionFactory 功能
    npx hardhat test test/NFTAuctionFactory.js --network localhost
    # 9、测试 NFTAuctionFactory 升级
    npx hardhat test test/NFTAuctionFactoryUpgrade.js --network localhost

-   测试 4 和 5 ：需要使用同一个网络，保证升级时可以读取部署时保存的 deployments 数据。
-   测试 7 : 这个测试需要使用 sepolia 测试网，需要使用 Chain Link 预言机。
-   测试 9 ：在 sepolia 中测试会出现缓存的情况，造成测试结果出错。建议使用 localhost 测试。
