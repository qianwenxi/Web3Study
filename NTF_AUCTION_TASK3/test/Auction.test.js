// 引入chai断言库的expect模块，用于测试断言
const { expect } = require("chai");
// 引入ethers库，用于与以太坊交互
require("hardhat");

// 定义测试套件：拍卖合约测试
describe("拍卖合约测试", function () {
  // 声明变量，用于存储合约实例
  let nft;               // NFT合约实例
  let auction;           // 拍卖合约实例
  // 声明变量，用于存储测试账户
  let deployer, seller, bidder1, bidder2;
  // 定义ETH/USD价格预言机地址（使用Chainlink的测试网地址）
  const ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";

  // 每个测试用例执行前的准备工作
  beforeEach(async function () {
    // 获取 获取测试账户列表，分别赋值给部署者、卖家、竞拍者1、竞拍者2
    [deployer, seller, bidder1, bidder2] = await ethers.getSigners();

    // 部署NFT合约
    // 获取NFT合约工厂
    const NFT = await ethers.getContractFactory("NFT");
    // 部署NFT合约，参数为名称、符号、元数据基础URL
    nft = await NFT.deploy("Test NFT", "TNFT", "https://test-metadata/");
    // 等待合约部署完成
    await nft.waitForDeployment();

    // 铸造NFT给卖家
    await nft.mint(seller.address);
    // 卖家取消对零地址的批量授权（清理状态）
    // await nft.connect(seller).setApprovalForAll(ethers.ZeroAddress, false);
    // await nft.connect(seller).setApprovalForAll(ethers.ZeroAddress, false);
    // 卖家取消单个NFT授权（清理状态）  当 to 为零地址且 tokenId 有效时，实际效果是取消该 tokenId 的授权
    await nft.connect(seller).approve(ethers.ZeroAddress, 0);
    // 卖家授权部署者地址管理其所有NFT（可能用于测试中的其他操作）
    await nft.connect(seller).setApprovalForAll(deployer.address, true);

    // 部署拍卖合约（使用升级合约模式）
    // 获取Auction合约工厂
    const Auction = await ethers.getContractFactory("Auction");
    // 部署代理合约，构造函数参数为ETH/USD价格预言机地址
    const auctionProxy = await upgrades.deployProxy(Auction, [ETH_USD_PRICE_FEED]);
    // 等待拍卖合约部署完成
    auction = await auctionProxy.waitForDeployment();

    // 卖家授权拍卖合约转移其NFT（为后续拍卖转移NFT做准备）
    await nft.connect(seller).setApprovalForAll(auction.target, true);
  });

  // 测试用例：应该成功创建拍卖
  it("应该创建拍卖", async function () {
    // 计算拍卖开始时间：当前时间+60秒（1分钟后开始）
    const blockNumBefore = await ethers.provider.getBlockNumber();
    const blockBefore = await ethers.provider.getBlock(blockNumBefore);
    const timestampBefore = blockBefore.timestamp;
    //const startTime = Math.floor(Date.now() / 1000) + 60;
    const startTime = timestampBefore + 2;
    // 计算拍卖结束时间：开始时间+3600秒（1小时后结束）
    const endTime = startTime + 60;

    // 执行创建拍卖操作，并验证事件是否正确触发
    await expect(auction.connect(seller).createAuction(
      nft.target,          // NFT合约地址
      0,                    // NFT的tokenId
      startTime,            // 开始时间
      endTime,              // 结束时间
      ethers.parseEther("0.1"),  // 起拍价（0.1 ETH）
      0,                    // 出价类型（0表示ETH）
      ethers.ZeroAddress     // 代币地址（ETH时为零地址）
    ))
      .to.emit(auction, "AuctionCreated")  // 验证AuctionCreated事件触发
      .withArgs(1, seller.address, 0);     // 验证事件参数：拍卖ID=1，卖家地址，tokenId=0

    // 获取创建的拍卖信息
    const auctionInfo = await auction.getAuctionInfo(1);
    // 验证拍卖信息是否正确
    expect(auctionInfo.seller).to.equal(seller.address);       // 卖家地址正确
    expect(auctionInfo.nftContract).to.equal(nft.target);     // NFT合约地址正确
    expect(auctionInfo.tokenId).to.equal(0);                   // tokenId正确
    expect(auctionInfo.startingPrice).to.equal(ethers.parseEther("0.1"));  // 起拍价正确
    expect(auctionInfo.status).to.equal(0);                    // 状态为PENDING（0）
  });

  // 测试用例：应该成功提交ETH出价
  it("应该提交ETH出价", async function () {
    // const blockNumBefore = await ethers.provider.getBlockNumber();
    // const blockBefore = await ethers.provider.getBlock(blockNumBefore);
    // const timestampBefore = blockBefore.timestamp;
    // 计算拍卖开始时间：当前时间-60秒（已开始）
    const blockNumBefore = await ethers.provider.getBlockNumber();
    const blockBefore = await ethers.provider.getBlock(blockNumBefore);
    const timestampBefore = blockBefore.timestamp;
    //const startTime = Math.floor(Date.now() / 1000) + 60;
    const startTime = timestampBefore + 2;
    // 计算拍卖结束时间：当前时间+3600秒（1小时后结束）
    const endTime = Math.floor(Date.now() / 1000) + 60;
    // 创建拍卖
    await auction.connect(seller).createAuction(
      nft.target,
      0,
      startTime,
      endTime,
      ethers.parseEther("0.1"),
      0,
      ethers.ZeroAddress
    );
      // 等待60秒（新增代码）
  await new Promise(resolve => setTimeout(resolve, 10000));
    // bidder1执行出价操作，并验证事件是否正确触发
    await expect(auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("0.2") }))
      .to.emit(auction, "BidPlaced")  // 验证BidPlaced事件触发
      .withArgs(1, bidder1.address, ethers.parseEther("0.2"));  // 验证事件参数：拍卖ID=1，出价者地址，出价金额=0.2ETH
    // 获取拍卖信息
    const auctionInfo = await auction.getAuctionInfo(1);
    // 验证出价后拍卖信息是否正确
    expect(auctionInfo.highestBid).to.equal(ethers.parseEther("0.2"));  // 最高出价正确
    expect(auctionInfo.highestBidder).to.equal(bidder1.address);              // 最高出价者正确
    expect(auctionInfo.status).to.equal(1);                                   // 状态为ACTIVE（1）

     // bidder2执行出价操作，并验证事件是否正确触发
    await expect(auction.connect(bidder2).placeBid(1, { value: ethers.parseEther("0.3") }))
      .to.emit(auction, "BidPlaced")  // 验证BidPlaced事件触发
      .withArgs(1, bidder2.address, ethers.parseEther("0.3"));  // 验证事件参数：拍卖ID=1，出价者地址，出价金额=0.2ETH
    // 获取拍卖信息
    const auctionInfo2 = await auction.getAuctionInfo(1);
    // 验证出价后拍卖信息是否正确
    expect(auctionInfo2.highestBid).to.equal(ethers.parseEther("0.3"));  // 最高出价正确
    expect(auctionInfo2.highestBidder).to.equal(bidder2.address);              // 最高出价者正确
    expect(auctionInfo2.status).to.equal(1);                                   // 状态为ACTIVE（1）
  });

  // 测试用例：应该成功结束拍卖并转移NFT和资金
  it("应该结束拍卖并转移NFT和资金", async function () {
    // 计算拍卖开始时间：当前时间-60秒（已开始）
        const blockNumBefore = await ethers.provider.getBlockNumber();
    const blockBefore = await ethers.provider.getBlock(blockNumBefore);
    const timestampBefore = blockBefore.timestamp;
    //const startTime = Math.floor(Date.now() / 1000) + 60;
    const startTime = timestampBefore + 2;
    // 计算拍卖结束时间：当前时间-10秒（已结束）
    const endTime = startTime + 20;

    // 创建拍卖
    await auction.connect(seller).createAuction(
      nft.target,
      0,
      startTime,
      endTime,
      ethers.parseEther("0.1"),
      0,
      ethers.ZeroAddress
    );

          // 等待60秒（新增代码）
  await new Promise(resolve => setTimeout(resolve, 10000));

    // 提交出价
    await auction.connect(bidder1).placeBid(1, { value: ethers.parseEther("0.2") });

    await new Promise(resolve => setTimeout(resolve, 15000));
    // 执行结束拍卖操作，并验证事件是否正确触发
    await expect(auction.endAuction(1))
      .to.emit(auction, "AuctionEnded")  // 验证AuctionEnded事件触发
      .withArgs(1, bidder1.address, ethers.parseEther("0.2"));  // 验证事件参数：拍卖ID=1，获胜者地址，成交价格

    // 验证NFT所有权已转移给最高出价者
    expect(await nft.ownerOf(0)).to.equal(bidder1.address);

    // 验证卖家收到资金（此处为简化验证，实际项目中需更精确计算gas等因素）
    const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);
    // 实际项目中需更精确验证，此处简化
  });

  // 测试用例：应该正确计算出价的美元价值  测试网验证 TODO
  // it("应该计算出价的美元价值", async function () {
  //   // 计算拍卖开始时间：当前时间-60秒（已开始）
  //           const blockNumBefore = await ethers.provider.getBlockNumber();
  //   const blockBefore = await ethers.provider.getBlock(blockNumBefore);
  //   const timestampBefore = blockBefore.timestamp;
  //   //const startTime = Math.floor(Date.now() / 1000) + 60;
  //   const startTime = timestampBefore + 2;
  //   // 计算拍卖结束时间：当前时间+3600秒（1小时后结束）
  //   const endTime = startTime + 3600;

  //   // 创建拍卖
  //   await auction.connect(seller).createAuction(
  //     nft.target,
  //     0,
  //     startTime,
  //     endTime,
  //     ethers.parseEther("0.1"),
  //     0,
  //     ethers.ZeroAddress
  //   );

  //   // 调用合约方法计算1 ETH对应的美元价值（通过Chainlink预言机）
  //   const usdValue = await auction.getBidInUSD(1, ethers.parseEther("1"));
  //   // 验证返回的美元价值大于0（确保预言机调用有效）
  //   expect(usdValue).to.be.gt(0);
  // });
});