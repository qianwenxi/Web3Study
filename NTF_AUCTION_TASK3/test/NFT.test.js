// 导入chai库中的expect断言函数，用于编写测试断言
const { expect } = require("chai");
// 导入hardhat中的ethers模块，用于与以太坊交互
const { ethers } = require("hardhat");

// 定义测试套件，描述为"NFT合约测试"
describe("NFT合约测试", function () {
  // 声明变量，用于存储合约实例和账户
  let nft; // NFT合约实例
  let deployer, user1; // 部署者和用户1的账户

  // 每个测试用例执行前都会运行的函数，用于初始化环境
  beforeEach(async function () {
    // 获取签名者列表，第一个为部署者，第二个为用户1
    [deployer, user1] = await ethers.getSigners();
    
    // 部署NFT合约
    const NFT = await ethers.getContractFactory("NFT"); // 获取NFT合约工厂
    // 部署合约，传入构造函数参数：名称、符号、基础URI
    // hardhat 默认使用第一个签名者部署合约
    // 如果需要指定账户用户部署合约 nft = await NFT.connect(user1).deploy("Test NFT", "TNFT", "https://test-metadata/"); // 指定user1账户部署 
    nft = await NFT.deploy("Test NFT", "TNFT", "https://test-metadata/");
    await nft.waitForDeployment(); // 等待合约部署完成
  });

  // 测试用例：应该正确初始化合约参数
  it("应该正确初始化合约参数", async function () {
    // 断言合约名称是否正确
    expect(await nft.name()).to.equal("Test NFT");
    // 断言合约符号是否正确
    expect(await nft.symbol()).to.equal("TNFT");
    // 断言基础TokenURI是否正确
    expect(await nft.baseTokenURI()).to.equal("https://test-metadata/");
    // 断言合约所有者是否为部署者
    expect(await nft.owner()).to.equal(deployer.address);
  });

  // 测试用例：应该允许所有者铸造NFT
  it("应该允许所有者铸造NFT", async function () {
    // 部署者为user1铸造NFT
    await nft.mint(user1.address);
    // 断言user1的NFT余额是否为1
    expect(await nft.balanceOf(user1.address)).to.equal(1);
    // 断言ID为0的NFT所有者是否为user1
    expect(await nft.ownerOf(0)).to.equal(user1.address);
  });

  // 测试用例：非所有者不能铸造NFT
//   it("非所有者不能铸造NFT", async function () {
//   // 断言非所有者调用mint会被revert，并匹配错误提示字符串
//   await expect(nft.connect(user1).mint(user1.address))
//     .to.be.revertedWith("Ownable: caller is not the owner");
// });
  it("非所有者不能铸造NFT", async function () {
    // 断言user1（非所有者）调用mint函数会被 revert
    // 并检查错误类型和参数是否符合预期
    await expect(nft.connect(user1).mint(user1.address))
      .to.be.revertedWithCustomError(nft, "OwnableUnauthorizedAccount")
      .withArgs(user1.address);
  });

  // 测试用例：应该正确返回Token URI
  it("应该正确返回Token URI", async function () {
    // 铸造一个NFT给user1
    await nft.mint(user1.address);
    // 断言ID为0的TokenURI是否正确（基础URI + ID）
    expect(await nft.tokenURI(0)).to.equal("https://test-metadata/0");
    
    // 更新基础TokenURI
    await nft.setBaseTokenURI("https://new-metadata/");
    // 断言更新后ID为0的TokenURI是否正确
    expect(await nft.tokenURI(0)).to.equal("https://new-metadata/0");
  });
});