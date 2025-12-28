// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IAuction {
    // 出价类型（ETH/ERC20）
    enum BidType { ETH, ERC20 }

    // 拍卖状态
    enum AuctionStatus { PENDING, ACTIVE, ENDED, CANCELLED }

    // 拍卖结构体
    struct AuctionInfo {
        uint256 auctionId;          // 拍卖ID
        address seller;             // 卖家地址
        address nftContract;        // NFT合约地址
        uint256 tokenId;            // NFT Token ID
        uint256 startTime;          // 开始时间（时间戳）
        uint256 endTime;            // 结束时间（时间戳）
        uint256 startingPrice;      // 起拍价（对应币种单位）
        uint256 highestBid;         // 最高出价（对应币种单位）
        address highestBidder;      // 最高出价者
        BidType bidType;            // 出价类型（ETH/ERC20）
        address erc20Token;         // ERC20代币地址（bidType为ERC20时有效）
        AuctionStatus status;       // 拍卖状态
    }

    // 事件定义
    event AuctionCreated(uint256 indexed auctionId, address indexed seller, uint256 indexed tokenId);
    event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount);
    event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 amount);
    event AuctionCancelled(uint256 indexed auctionId);

    // 核心方法
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startTime,
        uint256 endTime,
        uint256 startingPrice,
        BidType bidType,
        address erc20Token
    ) external returns (uint256);

    function placeBid(uint256 auctionId) external payable;
    function endAuction(uint256 auctionId) external;
    function cancelAuction(uint256 auctionId) external;
    function getAuctionInfo(uint256 auctionId) external view returns (AuctionInfo memory);
    function getBidInUSD(uint256 auctionId, uint256 amount) external view returns (uint256);
}