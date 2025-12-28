// 引入hardhat的network模块，用于获取当前网络相关配置
const { network, upgrades } = require("hardhat");

// 定义Sepolia测试网中Chainlink的ETH/USD价格预言机合约地址
const ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";

// 导出一个异步函数，该函数会在部署时被Hardhat执行
// 函数接收一个包含getNamedAccounts和deployments属性的对象作为参数
module.exports = async ({ getNamedAccounts, deployments }) => {
  // 从deployments对象中解构出deploy和log方法
  // deploy用于部署合约，log用于在控制台输出日志
  const { deploy,save, log } = deployments;
  
  // 调用getNamedAccounts方法获取命名账户，这里获取部署者账户
  const { deployer } = await getNamedAccounts();

  // 输出分隔线，用于区分不同部署步骤的日志
  log("------------------------------");
     // 通过代理合约部署拍卖合约
    const Auction = await ethers.getContractFactory("Auction");
    // const nftAuctionProxy = await upgrades.deployProxy(NFTAuction, [deployer.address, nftContract, 2, 10000, tokenId, ethers.ZeroAddress, ethers.ZeroAddress, 0], {
    //     initializer: "initialize",
    // })


  const auctionProxy = await upgrades.deployProxy(
    await ethers.getContractFactory('Auction'), 
    [ETH_USD_PRICE_FEED], // 构造函数参数，没有就空数组
    { 
      from: deployer, 
      initializer: 'initialize', // 对应你的初始化函数名，和你代码里的一致
      kind: 'uups' // 或者 'transparent'，看你要用的代理类型，ERC1967 常用 uups
    }
  );
  await auctionProxy.waitForDeployment();
  console.log('Auction deployed at:', auctionProxy.target);

  await save("Auction", {
        abi: Auction.interface.format("json"),
        address: auctionProxy.target,
    })
    console.log("代理拍卖合约部署完成，信息保存成功");

  const logicAddress = await upgrades.erc1967.getImplementationAddress(auctionProxy.target); 
  console.log('逻辑合约地址:', logicAddress);
  //console.log(await Auction.deployed());
//console.log(await upgrades.erc1967.getBeaconAddress(auctionImplementation.address));
  // 部署Auction合约的实现版本
  // const auctionImplementation = await deploy("Auction", {
  //   from: deployer, // 部署者的账户地址
  //   args: [], // 部署实现合约时的构造函数参数（此处无参数）
  //   log: true, // 启用部署日志，会在控制台输出部署相关信息
  //   // 等待的区块确认数，优先使用网络配置中的blockConfirmations，默认1
  //   waitConfirmations: network.config.blockConfirmations || 1,
  //   // 配置代理相关信息，使合约支持可升级
  //   proxy: {
  //     // 指定代理合约为OpenZeppelin的ERC1967Proxy（版本4.9.3）
  //     proxyContract: "OpenZeppelin/openzeppelin-contracts@4.9.3:ERC1967Proxy",
  //     // 配置代理的初始化操作
  //     execute: {
  //       methodName: "initialize", // 初始化方法名称
  //       args: [ETH_USD_PRICE_FEED], // 初始化方法的参数（传入价格预言机地址）
  //     },
  //   },
  // });

  // 输出可升级拍卖合约的部署地址
  //log(`可升级拍卖合约已部署至：${auctionImplementation.address}`);
};

// 为该部署脚本添加标签，方便通过标签执行部署（如执行部署所有标签为all或auction的脚本）
module.exports.tags = ["all", "auction"];