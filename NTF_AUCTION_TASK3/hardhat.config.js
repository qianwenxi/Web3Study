// require("@nomicfoundation/hardhat-chai-matchers");
// require("@nomiclabs/hardhat-ethers");
// require("hardhat-deploy");
// require("dotenv").config();

require("dotenv").config();
require("@nomicfoundation/hardhat-toolbox"); // v2 版本，适配 ethers v5
require("@nomicfoundation/hardhat-verify");
//require("@nomiclabs/hardhat-etherscan");
require("hardhat-deploy");
require("@openzeppelin/hardhat-upgrades");
require("solidity-coverage"); 
// 配置代理（关键）
const proxy = "http://127.0.0.1:7897";
process.env.HTTP_PROXY = proxy;
process.env.HTTPS_PROXY = proxy;
// 让Etherscan请求走代理，且忽略证书校验（避免代理HTTPS证书问题）
//process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";

const SEPOLIA_RPC_URL = process.env.SEPOLIA_RPC_URL || "https://sepolia.infura.io/v3/your-infura-key";
const PRIVATE_KEY = process.env.PRIVATE_KEY || "your-private-key";
const ETHERSCAN_API_KEY = process.env.ETHERSCAN_API_KEY || "your-etherscan-key";

module.exports = {
  solidity: {
    version: "0.8.28", // 兼容 Chainlink 和 OpenZeppelin
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  networks: {
    hardhat: {
      chainId: 31337,
    },
    sepolia: {
      url: SEPOLIA_RPC_URL,
      accounts: [PRIVATE_KEY],
      chainId: 11155111,
      blockConfirmations: 6, // 部署确认区块数
    },
  },
  namedAccounts: {
    deployer: {
      default: 0, // 默认第一个账户为部署者
      11155111: 0, // Sepolia 测试网
    },
  },
  // etherscan: {
  //         apiKey: {
  //           sepolia: ETHERSCAN_API_KEY
  //         },
  //         //url: "https://api-sepolia.etherscan.io/v2/api"
  //     },
  // etherscan: {
  //   apiKey: ETHERSCAN_API_KEY,
  // url: "https://api-sepolia.etherscan.io/v2/api"
  // },
  // etherscan: {
  //   apiKey: "CNQT7T8QDQA9VFGRZRJUJ8AFQA477SX8NW",
  //   apiUrl: "https://api-sepolia.etherscan.io/v2/api"
  // },
  // verify: {
  //   etherscan: {
  //     apiKey: "CNQT7T8QDQA9VFGRZRJUJ8AFQA477SX8NW",
  //     apiUrl: "https://api-sepolia.etherscan.io/v2/api"
  //   },
  // },
  // 关键：配置hardhat-deploy的Etherscan V2端点
  etherscan: {
    apiKey: ETHERSCAN_API_KEY
    // customChains: [{
    //   network: "sepolia",
    //   chainId: 11155111,
    //   urls: { apiURL: "https://api-sepolia.etherscan.io/api", 
    //     browserURL: "https://sepolia.etherscan.io" }
    // }]
  },
  // },
  mocha: {
    timeout: 500000, // 测试超时时间
  },
};