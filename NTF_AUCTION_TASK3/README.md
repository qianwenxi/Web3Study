# solidity 任务三 NFT拍卖合约项目文档
## 文档版本
v1.0

## 一、项目概述
### 1.1 项目背景
本项目基于以太坊生态开发去中心化NFT拍卖市场，采用Hardhat框架实现智能合约开发，集成Chainlink价格预言机实现多币种出价的美元价值统一计算，并通过UUPS代理模式支持合约无缝升级，满足NFT拍卖的核心业务需求。

### 1.2 核心特性
- **NFT管理**：基于ERC721标准实现NFT的铸造、转移、元数据管理；
- **多币种拍卖**：支持ETH/ERC20代币两种出价方式；
- **价格标准化**：集成Chainlink预言机，将出价金额转换为美元价值，方便用户比价；
- **合约可升级**：采用UUPS代理模式，支持合约逻辑迭代且不改变合约地址；
- **安全保障**：内置权限控制、重入攻击防护、参数校验等安全机制。

### 1.3 技术栈
| 分类         | 技术/工具                          | 版本       |
|--------------|------------------------------------|------------|
| 合约开发     | Solidity                           | 0.8.28     |
| 开发框架     | Hardhat                            | ^2.28.0    |
| 合约基础     | OpenZeppelin Contracts             | ^4.9.3     |
| 预言机       | Chainlink Price Feeds              | ^1.5.0     |
| 测试         | Chai + Hardhat Chai Matchers       | ^2.1.0     |
| 部署         | hardhat-deploy + Ethers.js         | ^1.0.4     |
| 网络         | Sepolia测试网 / Hardhat本地节点     | -          |

## 二、功能说明
### 2.1 NFT合约（NFT.sol）
实现ERC721标准NFT的核心功能，支持所有权管理和元数据配置：
| 功能         | 方法名                | 权限       | 说明                                   |
|--------------|-----------------------|------------|----------------------------------------|
| 铸造NFT      | mint(address to)      | 合约所有者 | 为指定地址铸造新NFT，返回Token ID      |
| 设置元数据URI | setBaseTokenURI(string) | 合约所有者 | 配置NFT元数据的基础URI前缀             |
| 获取元数据URI | tokenURI(uint256)     | 公开       | 根据Token ID返回完整的NFT元数据地址    |
| 权限管理     | Ownable内置方法       | 合约所有者 | 转让合约所有权、验证所有者身份         |

### 2.2 拍卖合约（Auction.sol）
核心业务合约，支持NFT拍卖全流程管理，集成Chainlink预言机和UUPS升级：
#### 2.2.1 核心枚举类型
| 枚举名       | 取值                | 说明                     |
|--------------|---------------------|--------------------------|
| BidType      | ETH / ERC20         | 出价币种类型             |
| AuctionStatus| PENDING / ACTIVE / ENDED / CANCELLED | 拍卖状态 |

#### 2.2.2 核心功能
| 功能         | 方法名                                                                 | 权限       | 说明                                                                 |
|--------------|------------------------------------------------------------------------|------------|----------------------------------------------------------------------|
| 创建拍卖     | createAuction(nftContract, tokenId, startTime, endTime, startingPrice, bidType, erc20Token) | NFT所有者 | 上架NFT拍卖，需授权合约转移NFT，返回拍卖ID                           |
| 提交出价     | placeBid(uint256 auctionId)                                            | 公开       | 按拍卖类型提交ETH/ERC20出价，自动验证加价幅度，退回前最高出价者资金   |
| 结束拍卖     | endAuction(uint256 auctionId)                                          | 公开       | 拍卖结束后执行，转移NFT给最高出价者，转移资金给卖家                   |
| 取消拍卖     | cancelAuction(uint256 auctionId)                                        | 拍卖卖家   | 仅无出价时可取消拍卖                                                 |
| 美元价格计算 | getBidInUSD(uint256 auctionId, uint256 amount)                          | 公开       | 通过Chainlink预言机将出价金额转换为美元价值（单位：美分）             |
| 预言机配置   | setERC20PriceFeed(address erc20Token, address priceFeed)                | 合约所有者 | 配置ERC20代币对应的Chainlink价格预言机地址                           |
| 合约升级     | upgradeTo(address newImplementation)                                    | 合约所有者 | 升级合约实现逻辑（UUPS模式）                                         |

#### 2.2.3 核心事件
| 事件名          | 触发场景               | 关键参数                          |
|-----------------|------------------------|-----------------------------------|
| AuctionCreated  | 创建拍卖成功           | auctionId, seller, tokenId        |
| BidPlaced       | 提交出价成功           | auctionId, bidder, amount         |
| AuctionEnded    | 拍卖结束成功           | auctionId, winner, amount         |
| AuctionCancelled| 拍卖取消成功           | auctionId                         |

### 2.3 合约接口（IAuction.sol）
定义拍卖合约的接口规范，包含数据结构、枚举、方法和事件声明，保证合约交互的规范性和可扩展性。

## 三、部署步骤
### 3.1 环境准备
#### 3.1.1 前置条件
- 已安装Node.js（v16+）和npm（v8+）；
- 拥有MetaMask钱包（Sepolia测试网账户），并充值少量Sepolia ETH（可通过Chainlink/Infura水龙头获取）；
- 注册Infura/Alchemy账号（获取Sepolia RPC URL）；
- 注册Etherscan账号（获取API Key，用于合约验证）。

#### 3.1.2 环境变量配置
在项目根目录创建`.env`文件，填写以下内容：
```env
# Sepolia测试网RPC地址（Infura/Alchemy）
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/你的Infura API Key
# 钱包私钥（仅测试网，切勿使用主网私钥）
PRIVATE_KEY=你的MetaMask测试网私钥
# Etherscan API Key（合约验证用）
ETHERSCAN_API_KEY=你的Etherscan API Key
```

### 3.2 项目初始化（略）

### 3.3 合约编译
```bash
# 编译Solidity合约（生成ABI和字节码）
npx hardhat compile

# （可选）清除编译缓存（遇到编译错误时执行）
npx hardhat clean
```
编译成功后会生成`artifacts`（合约ABI/字节码）和`cache`（编译缓存）文件夹，ABI是合约与外部交互的核心接口。

### 3.4 本地测试
#### 3.4.1 运行测试用例
```bash
# 运行所有测试用例（NFT + 拍卖 + 合约升级）
npx hardhat test

# 仅运行指定测试文件
npx hardhat test test/NFT.test.js       # 仅测试NFT合约
npx hardhat test test/Auction.test.js   # 仅测试拍卖合约
```

#### 3.4.2 本地节点部署（模拟测试）
```bash
# 1. 启动Hardhat本地节点（新开终端，保持运行）
npx hardhat node

# 2. 部署合约到本地节点
npx hardhat deploy --network hardhat
```
本地节点会生成10个测试账户（带私钥），部署日志会输出合约地址，可用于本地交互测试。

### 3.5 Sepolia测试网部署
#### 3.5.1 部署所有合约
```bash
# 部署NFT + 拍卖合约到Sepolia
npx hardhat deploy --network sepolia
```

#### 3.5.2 单独部署/升级合约
```bash
# 仅部署NFT合约
npx hardhat deploy --tags nft --network sepolia

# 仅部署拍卖合约
npx hardhat deploy --tags auction --network sepolia

# 升级拍卖合约（更新逻辑时）
npx hardhat deploy --tags upgrade --network sepolia
```
部署成功后，终端会输出合约地址（需记录NFT合约地址、拍卖代理合约地址），等待6个区块确认后完成部署。

### 3.6 合约验证（Etherscan）
部署后需验证合约，使代码在Etherscan上开源可见：
```bash
# 验证NFT合约（替换为实际部署地址）
npx hardhat verify --network sepolia 0xNFT合约地址 "NFT Auction Demo" "NFTAD" "https://ipfs.io/ipfs/your-metadata-uri/"

# 验证拍卖实现合约（替换为实际部署地址）
npx hardhat verify --network sepolia 0x拍卖实现合约地址
```
验证成功后，可在Etherscan上查看合约代码、调用合约方法。并且在etherscan上完成代理合约的验证

### 3.7 测试NFT铸造
```bash
# 1. 修改scripts/mint-nft.js中的nftAddress为实际部署的NFT合约地址
# 2. 执行铸造脚本
npx hardhat run scripts/mint-nft.js --network sepolia
```
执行成功后，终端会输出NFT铸造的交易哈希，可在Etherscan上查看NFT持有地址。

## 四、合约交互示例
### 4.1 控制台交互（Sepolia测试网）
```bash
# 启动Hardhat控制台（连接Sepolia）
npx hardhat console --network sepolia

# 示例1：铸造NFT
const NFT = await ethers.getContractFactory("NFT");
const nft = await NFT.attach("0xNFT合约地址");
await nft.mint("你的钱包地址");

# 示例2：创建拍卖
const Auction = await ethers.getContractFactory("Auction");
const auction = await Auction.attach("0x拍卖代理合约地址");
const startTime = Math.floor(Date.now()/1000) + 60; // 1分钟后开始
const endTime = startTime + 3600; // 1小时后结束
await auction.createAuction(
  "0xNFT合约地址",
  0, // Token ID
  startTime,
  endTime,
  ethers.utils.parseEther("0.1"), // 起拍价0.1 ETH
  0, // BidType.ETH
  ethers.constants.AddressZero
);

# 示例3：提交ETH出价
await auction.placeBid(1, { value: ethers.utils.parseEther("0.2") });

# 示例4：结束拍卖
await auction.endAuction(1);
```

### 4.2 核心参数说明
- `auctionId`：拍卖创建时返回的唯一ID；
- `startingPrice`：起拍价（ETH/ERC20代币单位，需大于0）；
- `startTime/endTime`：拍卖开始/结束时间戳（秒级，endTime需晚于startTime）；
- `bidType`：0=ETH，1=ERC20（ERC20需提前配置预言机地址）。

## 五、常见问题解决
### 5.1 部署失败
- 检查`.env`文件中RPC URL、私钥是否正确；
- 确认钱包账户有足够的Sepolia ETH（支付Gas费）；
- 更换RPC地址（如Infura切换为Alchemy）。

### 5.2 测试用例失败
- 确保Solidity版本（0.8.20）与合约、Hardhat配置一致；
- 检查Chainlink预言机地址是否适配Sepolia测试网（ETH/USD：0x694AA1769357215DE4FAC081bf1f309aDC325306）；
- 重置本地节点：`npx hardhat node --reset`。

### 5.3 合约验证失败
- 确保部署时的构造函数参数与验证时一致；
- 等待部署交易上链（至少6个区块确认）后再验证；
- 检查Hardhat版本与`@nomiclabs/hardhat-verify`版本兼容性。

## 六、安全注意事项
1. 私钥仅用于测试网，切勿将主网私钥写入`.env`文件；
2. 合约升级前需充分测试新实现逻辑，避免存储结构不兼容；
3. ERC20出价需确保用户已授权合约转移代币（approve）；
4. 拍卖结束后需及时调用`endAuction`，避免资金/NFT锁定。

## 七、附录
### 7.1 Chainlink预言机地址（Sepolia）
| 代币   | 预言机地址                          |
|--------|-------------------------------------|
| ETH/USD| 0x694AA1769357215DE4FAC081bf1f309aDC325306 |
| USDC/USD| 0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E |
| LINK/USD| 0xc59E3633BAAC79493d908e63626716e204A45EdF |

### 7.2 项目目录结构
```
NFT auction/
├── artifacts/            # 编译产物（ABI/字节码）
├── cache/                # 编译缓存
├── contracts/            # 智能合约
│   ├── NFT.sol           # ERC721 NFT合约
│   ├── Auction.sol       # 可升级拍卖合约
│   └── interfaces/       # 接口定义
│       └── IAuction.sol
├── deploy/               # 部署脚本
│   ├── 01-deploy-nft.js
│   ├── 02-deploy-auction.js
│   └── 03-upgrade-auction.js
├── test/                 # 测试用例
│   ├── NFT.test.js
│   ├── Auction.test.js
│   └── Upgrade.test.js
├── scripts/              # 辅助脚本
│   └── mint-nft.js
├── .env                  # 环境变量
├── hardhat.config.js     # Hardhat配置
├── package.json          # 依赖配置
└── README.md             # 项目文档
```

