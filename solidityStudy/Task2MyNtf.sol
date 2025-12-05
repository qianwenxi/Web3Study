// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

// 导入OpenZeppelin的ERC721核心库
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
// 导入自动递增tokenId的库
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
// 导入元数据（tokenURI）的库
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
// 导入权限控制（可选，这里简化为所有人可mint）
import "@openzeppelin/contracts/access/Ownable.sol";
// 导入安全数学库
import "@openzeppelin/contracts/utils/Counters.sol";

contract MyNFT is ERC721, ERC721Enumerable, ERC721URIStorage, Ownable {
    // 用于自动递增tokenId
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;

    // 构造函数：设置NFT名称、符号，同时初始化Owner为部署者
    constructor() ERC721("MyTestNFT", "LLNTF") Ownable(msg.sender) {}

    // mintNFT函数：接收收件人地址、元数据IPFS链接
    function mintNFT(address to, string memory uri) public {
        // 获取当前tokenId（从1开始）
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();
        // 铸造NFT到to地址
        _safeMint(to, tokenId);
        // 关联元数据链接
        _setTokenURI(tokenId, uri);
    }

    // 以下是ERC721Enumerable、ERC721URIStorage的必要重写函数
    function _update(address to, uint256 tokenId, address auth)
        internal
        override(ERC721, ERC721Enumerable)
        returns (address)
    {
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(address account, uint128 value)
        internal
        override(ERC721, ERC721Enumerable)
    {
        super._increaseBalance(account, value);
    }

    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721Enumerable, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}