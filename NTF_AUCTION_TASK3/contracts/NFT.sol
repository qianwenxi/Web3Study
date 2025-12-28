// 声明合约的开源许可证类型为MIT，这是区块链合约常用许可证，允许自由使用、修改和分发
// SPDX-License-Identifier 是Solidity编译器要求的标准化注释，用于明确许可证信息
// SPDX-License-Identifier: MIT

// 声明合约使用的Solidity编译器版本，^0.8.20 表示兼容0.8.20及以上版本（但不包含0.9.0及更高主版本）
pragma solidity ^0.8.20;

// 从OpenZeppelin库导入ERC721合约：ERC721是NFT（非同质化代币）的标准实现，提供了NFT的核心功能（如转账、余额查询、所有者查询等）
// OpenZeppelin是行业公认的安全合约库，其实现经过严格审计，避免重复开发和安全漏洞
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

// 从OpenZeppelin库导入Ownable合约：提供合约所有权管理功能（如指定所有者、仅所有者可调用函数等）
// 用于实现"仅合约所有者可执行特定操作"的权限控制逻辑（如本合约的mint、setBaseTokenURI函数）
import "@openzeppelin/contracts/access/Ownable.sol";

// 从OpenZeppelin库导入Counters合约：提供安全的计数器功能，用于自动生成唯一的NFT Token ID
// 避免手动管理计数器时可能出现的并发问题或逻辑漏洞
//import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/**
 * @title NFT：合约名称，用于标识合约功能（基础ERC721 NFT合约）
 * @dev 合约功能说明：基础ERC721 NFT合约，支持铸造（mint）和转移（继承自ERC721）
 * @dev 这是NatSpec注释格式（Solidity官方推荐），用于生成文档和提示开发者
 */
contract NFT is ERC721, Ownable {
    // 引入Counters库的Counter类型，并创建一个私有计数器变量_tokenIdCounter
    // 私有变量（private）表示仅合约内部可访问，外部无法直接修改
    // 用于记录当前已铸造的NFT数量，每次铸造时自动递增，生成唯一Token ID
    //using Counters for Counters.Counter;
    using Strings for uint256;

    uint256 private _tokenIdCounter;

    // 声明公共（public）状态变量baseTokenURI，用于存储NFT元数据的基础URI前缀
    // public修饰符会自动生成一个同名的getter函数（baseTokenURI()），外部可通过该函数查询值
    // NFT元数据（如图片、描述）的完整URI = baseTokenURI + Token ID（如"https://test-metadata/" + "0"）
    string public baseTokenURI;

    /**
     * @dev 合约构造函数：合约部署时自动执行，用于初始化合约状态
     * @param name：NFT合约的名称（如"Test NFT"），继承自ERC721合约的构造参数
     * @param symbol：NFT合约的符号（如"TNFT"），继承自ERC721合约的构造参数
     * @param _baseTokenURI：NFT元数据的基础URI前缀，赋值给本合约的baseTokenURI变量
     * @dev 构造函数后接的 ERC721(name, symbol) 表示调用父合约ERC721的构造函数，初始化名称和符号
     * @dev Ownable() 表示调用父合约Ownable的构造函数，默认将合约部署者设为所有者
     */
    constructor(
        string memory name,
        string memory symbol,
        string memory _baseTokenURI
    ) ERC721(name, symbol) Ownable(msg.sender) {
        // 将传入的基础URI前缀赋值给状态变量baseTokenURI，完成初始化
        baseTokenURI = _baseTokenURI;
    }

    /**
     * @dev 铸造NFT的核心函数：创建新的NFT并分配给指定接收者
     * @param to NFT接收者的地址（新NFT将归属该地址）
     * @return tokenId 新铸造NFT的唯一Token ID（返回给调用者，方便后续操作）
     * @dev external修饰符：表示该函数仅可被外部账户或其他合约调用，合约内部无法调用
     * @dev onlyOwner修饰符：继承自Ownable合约，限制仅合约所有者可调用该函数（防止任意地址铸造NFT）
     */
    function mint(address to) external onlyOwner returns (uint256) {
        // 获取当前计数器的值，作为新NFT的Token ID（初始值为0，首次铸造时Token ID为0）
        uint256 tokenId = _tokenIdCounter;
        // 计数器自增1，确保下一次铸造时使用新的Token ID（避免重复）
        _tokenIdCounter++;
        // 调用ERC721父合约的_safeMint函数：安全铸造NFT并分配给to地址
        // _safeMint会检查to地址是否支持接收NFT（如为合约地址，需实现ERC721Receiver接口），避免NFT丢失
        _safeMint(to, tokenId);
        // 返回新铸造NFT的Token ID
        return tokenId;
    }

    /**
     * @dev 更新NFT元数据基础URI的函数：修改baseTokenURI的值
     * @param _newBaseURI：新的基础URI前缀（如从"https://test-metadata/"改为"https://new-metadata/"）
     * @dev external修饰符：仅外部可调用
     * @dev onlyOwner修饰符：仅合约所有者可调用（防止任意地址修改元数据URI）
     */
    function setBaseTokenURI(string memory _newBaseURI) external onlyOwner {
        // 将新的基础URI赋值给baseTokenURI，覆盖旧值（后续查询Token URI时将使用新前缀）
        baseTokenURI = _newBaseURI;
    }

    /**
     * @dev 重写（override）ERC721父合约的tokenURI函数：自定义NFT元数据的完整URI生成逻辑
     * @param tokenId 需要查询URI的NFT的Token ID
     * @return uri 该NFT的完整元数据URI（baseTokenURI + Token ID的字符串形式）
     * @dev public view修饰符：public表示外部可调用，view表示函数仅读取状态（不修改合约数据，无需消耗gas）
     * @dev override修饰符：声明该函数重写了父合约ERC721中的同名函数（必须与父函数签名一致）
     */
    function tokenURI(
        uint256 tokenId
    ) public view override returns (string memory) {
        // 调用ERC721父合约的ownerOf函数：检查传入的tokenId是否已铸造（若未铸造，会抛出错误）
        // 作用是验证Token ID的有效性，避免查询不存在的NFT的URI
        ownerOf(tokenId);
        // 获取当前的基础URI前缀（baseTokenURI）
        string memory baseURI = baseTokenURI;
        // 三元运算符逻辑：
        // 1. 若baseURI的字节长度>0（即基础URI已设置），则拼接baseURI和Token ID的字符串形式（如"https://test-metadata/" + "0"）
        // 2. 若baseURI为空，则返回空字符串（表示该NFT暂无元数据URI）
        // abi.encodePacked：Solidity内置函数，用于将多个参数按字节直接拼接（不添加额外填充）
        // Strings.toString(tokenId)：将uint256类型的Token ID转为字符串（需确保项目中导入了Strings库，OpenZeppelin的ERC721会自动依赖）
        return
            bytes(baseURI).length > 0
                ? string(abi.encodePacked(baseURI, tokenId.toString()))
                : "";
    }
}
