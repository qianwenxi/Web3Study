// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "hardhat/console.sol"; 
import "./interfaces/IAuction.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
//import "@chainlink/contracts/src/interfaces/AggregatorV3Interface.sol";

/**
 * @title Auction
 * @dev 可升级的NFT拍卖合约，支持ETH/ERC20出价，集成Chainlink预言机计算美元价格
 */
contract Auction is IAuction, OwnableUpgradeable, UUPSUpgradeable, ReentrancyGuardUpgradeable {
    // 拍卖ID计数器
    uint256 private _auctionIdCounter;
    
    // 拍卖信息映射（auctionId => AuctionInfo）
    mapping(uint256 => AuctionInfo) private _auctions;

    // Chainlink价格预言机地址（Sepolia测试网）
    AggregatorV3Interface public ethUsdPriceFeed;
    mapping(address => AggregatorV3Interface) public erc20UsdPriceFeeds;

    // 初始化函数（替代构造函数）
    function initialize(
        address _ethUsdPriceFeed
    ) external initializer {
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        ethUsdPriceFeed = AggregatorV3Interface(_ethUsdPriceFeed);
        _auctionIdCounter = 1; // 拍卖ID从1开始
    }

    /**
     * @dev UUPS升级权限控制（仅所有者可升级）
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    /**
     * @dev 设置ERC20代币的价格预言机
     * @param erc20Token ERC20代币地址
     * @param priceFeed 对应的Chainlink预言机地址
     */
    function setERC20PriceFeed(address erc20Token, address priceFeed) external onlyOwner {
        erc20UsdPriceFeeds[erc20Token] = AggregatorV3Interface(priceFeed);
    }

    /**
     * @dev 创建拍卖
     * @param nftContract NFT合约地址
     * @param tokenId NFT Token ID
     * @param startTime 开始时间戳
     * @param endTime 结束时间戳
     * @param startingPrice 起拍价
     * @param bidType 出价类型（ETH/ERC20）
     * @param erc20Token ERC20代币地址（bidType为ERC20时必填）
     * @return 新创建的拍卖ID
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startTime,
        uint256 endTime,
        uint256 startingPrice,
        BidType bidType,
        address erc20Token
    ) external override returns (uint256) {
        // 验证参数
        console.log("Auction status updated to ACTIVE | Current Time:", block.timestamp, "| Auction Start Time:", startTime);
        require(startTime > block.timestamp, "Auction: start time must be in future");
        require(endTime > startTime, "Auction: end time must be after start time");
        require(startingPrice > 0, "Auction: starting price must be > 0");
        require(bidType == BidType.ETH || erc20Token != address(0), "Auction: ERC20 token address required");

        // 验证卖家拥有NFT且授权合约转移
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Auction: not NFT owner");
        require(nft.isApprovedForAll(msg.sender, address(this)) || nft.getApproved(tokenId) == address(this), 
            "Auction: contract not approved to transfer NFT");

        // 创建拍卖
        uint256 auctionId = _auctionIdCounter;
        _auctions[auctionId] = AuctionInfo({
            auctionId: auctionId,
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            startTime: startTime,
            endTime: endTime,
            startingPrice: startingPrice,
            highestBid: 0,
            highestBidder: address(0),
            bidType: bidType,
            erc20Token: erc20Token,
            status: AuctionStatus.PENDING
        });

        _auctionIdCounter++;
        emit AuctionCreated(auctionId, msg.sender, tokenId);
        return auctionId;
    }

    /**
     * @dev 提交出价
     * @param auctionId 拍卖ID
     */
    function placeBid(uint256 auctionId) external payable override nonReentrant {
        console.log("auction.upgraded:new palce bid log");
        AuctionInfo storage auction = _auctions[auctionId];
        require(auction.auctionId != 0, "Auction: auction does not exist");
        require(block.timestamp >= auction.startTime && block.timestamp <= auction.endTime, "Auction: not active");
        require(auction.status == AuctionStatus.PENDING || auction.status == AuctionStatus.ACTIVE, "Auction: not active");
        require(msg.sender != auction.seller, "Auction: seller cannot bid");

        // 更新拍卖状态为ACTIVE
        if (auction.status == AuctionStatus.PENDING) {
            auction.status = AuctionStatus.ACTIVE;
        }

        // 计算出价金额
        uint256 bidAmount;
        if (auction.bidType == BidType.ETH) {
            require(msg.value > 0, "Auction: ETH bid amount must be > 0");
            bidAmount = msg.value;
        } else {
            // ERC20出价：验证授权并转账
            IERC20 token = IERC20(auction.erc20Token);
            bidAmount = token.allowance(msg.sender, address(this));
            require(bidAmount > 0, "Auction: ERC20 allowance required");
            token.transferFrom(msg.sender, address(this), bidAmount);
        }

        // 验证出价高于当前最高出价（首次出价需高于起拍价）
        uint256 minBid = auction.highestBid > 0 ? auction.highestBid : auction.startingPrice;
        require(bidAmount > minBid, "Auction: bid must be higher than current highest bid");

        // 退回之前最高出价者的资金
        if (auction.highestBidder != address(0)) {
            _refund(auction.highestBidder, auction.highestBid, auction.bidType, auction.erc20Token);
        }

        // 更新最高出价
        auction.highestBid = bidAmount;
        auction.highestBidder = msg.sender;

        emit BidPlaced(auctionId, msg.sender, bidAmount);
    }

    /**
     * @dev 结束拍卖
     * @param auctionId 拍卖ID
     */
    function endAuction(uint256 auctionId) external override nonReentrant {
        AuctionInfo storage auction = _auctions[auctionId];
        require(auction.auctionId != 0, "Auction: auction does not exist");
        require(block.timestamp > auction.endTime, "Auction: not ended yet");
        require(auction.status == AuctionStatus.ACTIVE, "Auction: not active");

        // 更新拍卖状态
        auction.status = AuctionStatus.ENDED;
console.log("auction.highestBidder:", auction.highestBidder);
        if (auction.highestBidder != address(0)) {
            // 转移NFT给出价最高者
            IERC721(auction.nftContract).safeTransferFrom(auction.seller, auction.highestBidder, auction.tokenId);
            
            // 转移资金给卖家
            _transferFunds(auction.seller, auction.highestBid, auction.bidType, auction.erc20Token);
        } else {
            // 无出价，NFT归还给卖家（无需操作，卖家仍为所有者）
        }

        emit AuctionEnded(auctionId, auction.highestBidder, auction.highestBid);
    }

    /**
     * @dev 取消拍卖（仅卖家可取消，且无出价时）
     * @param auctionId 拍卖ID
     */
    function cancelAuction(uint256 auctionId) external override {
        AuctionInfo storage auction = _auctions[auctionId];
        require(auction.auctionId != 0, "Auction: auction does not exist");
        require(msg.sender == auction.seller, "Auction: only seller can cancel");
        require(auction.highestBid == 0, "Auction: cannot cancel with bids");
        require(auction.status != AuctionStatus.ENDED, "Auction: already ended");

        auction.status = AuctionStatus.CANCELLED;
        emit AuctionCancelled(auctionId);
    }

    /**
     * @dev 获取拍卖信息
     * @param auctionId 拍卖ID
     * @return 拍卖详情
     */
    function getAuctionInfo(uint256 auctionId) external view override returns (AuctionInfo memory) {
        return _auctions[auctionId];
    }

    /**
     * @dev 获取出价金额对应的美元价值
     * @param auctionId 拍卖ID
     * @param amount 出价金额（对应币种单位）
     * @return 美元价值（单位：美分，避免小数）
     */
    function getBidInUSD(uint256 auctionId, uint256 amount) external view override returns (uint256) {
        AuctionInfo memory auction = _auctions[auctionId];
        require(auction.auctionId != 0, "Auction: auction does not exist");

        uint256 priceInUsd;
        if (auction.bidType == BidType.ETH) {
            // 获取ETH/USD价格（Chainlink预言机）
            (, int256 price, , , ) = ethUsdPriceFeed.latestRoundData();
            // 价格格式：1 ETH = price * 10^8 USD（Chainlink返回8位小数）
            priceInUsd = (uint256(price) * amount) / 1e18; // amount是wei单位，转换为ETH后乘价格
        } else {
            // 获取ERC20/USD价格
            AggregatorV3Interface priceFeed = erc20UsdPriceFeeds[auction.erc20Token];
            require(address(priceFeed) != address(0), "Auction: no price feed for ERC20");
            (, int256 price, , , ) = priceFeed.latestRoundData();
            // 假设ERC20是18位小数，价格返回8位小数
            priceInUsd = (uint256(price) * amount) / 1e18;
        }

        return priceInUsd; // 单位：美分（例如100 = $1.00）
    }

    /**
     * @dev 内部方法：转移资金给卖家
     */
    function _transferFunds(address recipient, uint256 amount, BidType bidType, address erc20Token) private {
        if (bidType == BidType.ETH) {
            (bool success, ) = recipient.call{value: amount}("");
            require(success, "Auction: ETH transfer failed");
        } else {
            IERC20(erc20Token).transfer(recipient, amount);
        }
    }

    /**
     * @dev 内部方法：退回资金给出价者
     */
    function _refund(address recipient, uint256 amount, BidType bidType, address erc20Token) private {
        if (bidType == BidType.ETH) {
            (bool success, ) = recipient.call{value: amount}("");
            require(success, "Auction: ETH refund failed");
        } else {
            IERC20(erc20Token).transfer(recipient, amount);
        }
    }

    // 接收ETH
    receive() external payable {}
}