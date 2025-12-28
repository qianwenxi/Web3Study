// 从hardhat库中导入network模块，用于获取当前部署的网络信息
const { network } = require("hardhat");

// 导出一个异步函数，该函数会在Hardhat部署脚本执行时被调用
// 函数接收一个包含getNamedAccounts和deployments属性的对象作为参数
// getNamedAccounts用于获取配置中定义的命名账户，deployments用于部署合约等操作
module.exports = async ({ getNamedAccounts, deployments }) => {
  // 从deployments对象中解构出deploy和log方法
  // deploy用于部署合约，log用于在控制台输出日志信息
  const { deploy, log } = deployments;
  
  // 调用getNamedAccounts方法获取命名账户，并从中解构出deployer账户（部署者账户）
  const { deployer } = await getNamedAccounts();

  // 输出分隔线，用于在控制台中区分不同的部署步骤或信息，增强可读性
  log("------------------------------");
  
  // 定义部署NFT合约时需要传入的构造函数参数数组
  const nftArgs = [
    "NFT Auction Demo", // 第一个参数：NFT合约的名称
    "NFTAD",            // 第二个参数：NFT合约的符号（缩写）
    "https://ipfs.io/ipfs/your-metadata-uri/" // 第三个参数：NFT元数据的URI前缀，后续会与tokenId拼接形成完整的元数据地址
  ];

  // 调用deploy方法部署名为"NFT"的合约
  const nftContract = await deploy("NFT", {
    from: deployer, // 指定部署者账户为deployer
    args: nftArgs, // 传入构造函数参数数组
    log: true, // 设为true表示部署过程中会输出相关日志
    // 等待区块确认的数量，优先使用当前网络配置的blockConfirmations，若未配置则默认等待1个区块确认
    waitConfirmations: network.config.blockConfirmations || 1,
  });

  // 输出NFT合约部署后的地址信息
  log(`NFT合约已部署至：${nftContract.address}`);
};

// 为当前部署脚本添加标签，标签用于在执行部署命令时指定要运行的脚本
// 这里标签为["all", "nft"]，表示可以通过执行"hardhat deploy --tags all"或"hardhat deploy --tags nft"来运行该脚本
module.exports.tags = ["all", "nft"];