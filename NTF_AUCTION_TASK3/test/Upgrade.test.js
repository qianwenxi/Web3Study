const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("合约升级测试", function () {
  let auction;
  let auctionV2;
  let deployer;
  const ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";

  beforeEach(async function () {
    [deployer] = await ethers.getSigners();

    // 部署初始版本
    const Auction = await ethers.getContractFactory("Auction");
    auction = await upgrades.deployProxy(Auction, [ETH_USD_PRICE_FEED]);
    await auction.waitForDeployment();
  });

  it("应该成功升级合约", async function () {
    // 部署新版本实现（此处复用Auction合约，实际项目中可修改代码）
    const AuctionV2 = await ethers.getContractFactory("Auction");
    auctionV2 = await upgrades.upgradeProxy(auction.target, AuctionV2);

    // 验证升级后地址不变
    expect(auctionV2.target).to.equal(auction.target);

    // 验证初始化参数仍有效
    expect(await auctionV2.ethUsdPriceFeed()).to.equal(ETH_USD_PRICE_FEED);
  });
});