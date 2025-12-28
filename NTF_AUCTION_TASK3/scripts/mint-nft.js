// scripts/mint-nft.js
const { ethers } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();
  const nftAddress = "0xA62e291C2AAd58EC9fad6f8ED2FF2c4648B20aB6";
  
  const nft = await ethers.getContractAt("NFT", nftAddress);
  const tx = await nft.mint(deployer.address);
  await tx.wait();
  
  console.log(`NFT铸造成功，Token ID: 0`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});