const { network } = require("hardhat");

module.exports = async ({ getNamedAccounts, deployments }) => {
  const { deploy, log, get } = deployments;
  const { deployer } = await getNamedAccounts();

  log("------------------------------");
  // // 获取现有代理合约
  // const auctionProxy = await get("Auction");
  
  // // 部署新的实现合约
  // const newAuctionImplementation = await deploy("Auction", {
  //   from: deployer,
  //   args: [],
  //   log: true,
  //   waitConfirmations: network.config.blockConfirmations || 1,
  // });

  // // 升级代理合约
  // const auctionContract = await ethers.getContractAt("Auction", auctionProxy.address);
  // await auctionContract.upgradeTo(newAuctionImplementation.address);

  //log(`拍卖合约已升级至新实现：${newAuctionImplementation.address}`);

  // 获取新的实现合约工厂
  const Auction = await ethers.getContractFactory("Auction");
  
  // 安全升级代理合约（使用官方升级插件）
  const upgradedAuction = await upgrades.upgradeProxy(
    (await deployments.get("Auction")).address, // 现有代理地址
    // '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512',
    Auction // 新的实现合约
  );
  
  // 等待升级交易确认
  await upgradedAuction.waitForDeployment();

  log(`拍卖合约已升级至新实现：${upgradedAuction.target}`);
    const logicAddress = await upgrades.erc1967.getImplementationAddress(upgradedAuction.target); 
    console.log('逻辑合约地址:', logicAddress);
};

module.exports.tags = ["upgrade"];