// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

contract task2MyERC20 {
        // ========== ERC20标准状态变量 ==========
    string public constant name = unicode"LLERC20学习币"; // 代币名称
    string public constant symbol = "LL";   // 代币符号
    uint8 public constant decimals = 18;     // 代币精度（默认18位）
    uint256 public totalSupply;             // 总供应量

    // 地址→余额映射
    mapping(address => uint256) public balanceOf;
    // 授权映射：owner→spender→授权额度
    mapping(address => mapping(address => uint256)) public allowance;


    // ========== 权限控制（合约所有者） ==========
    address public owner; // 合约所有者（部署者）

    // 构造函数：部署时初始化所有者
    constructor() {
        owner = msg.sender;
    }

    // 仅所有者可调用的修饰符
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }


    // ========== ERC20标准函数 ==========
    /**
     * @dev 转账：从调用者地址转代币给_to
     * @param _to 接收地址
     * @param _value 转账数量（已包含精度）
     */
    function transfer(address _to, uint256 _value) public returns (bool success) {
        require(_to != address(0), "Invalid recipient address"); // 禁止转0地址
        require(balanceOf[msg.sender] >= _value, "Insufficient balance"); // 余额足够

        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;
        emit Transfer(msg.sender, _to, _value); // 触发转账事件
        return true;
    }

    /**
     * @dev 授权：允许_spender从调用者地址转走最多_value的代币
     * @param _spender 被授权地址
     * @param _value 授权额度（已包含精度）
     */
    function approve(address _spender, uint256 _value) public returns (bool success) {
        require(_spender != address(0), "Invalid spender address");

        allowance[msg.sender][_spender] = _value;
        emit Approval(msg.sender, _spender, _value); // 触发授权事件
        return true;
    }

    /**
     * @dev 代转账：从_from地址转代币给_to（需先授权）
     * @param _from 转出地址
     * @param _to 接收地址
     * @param _value 转账数量（已包含精度）
     */
    function transferFrom(
        address _from,
        address _to,
        uint256 _value
    ) public returns (bool success) {
        require(_to != address(0), "Invalid recipient address");
        require(balanceOf[_from] >= _value, "Insufficient balance in _from");
        require(allowance[_from][msg.sender] >= _value, "Allowance not enough"); // 授权额度足够

        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;
        allowance[_from][msg.sender] -= _value; // 扣减授权额度
        emit Transfer(_from, _to, _value);
        return true;
    }


    // ========== 增发函数（仅所有者） ==========
    /**
     * @dev 增发代币：给_to地址增发_value数量的代币
     * @param _to 接收增发代币的地址
     * @param _value 增发数量（已包含精度）
     */
    function mint(address _to, uint256 _value) public onlyOwner returns (bool success) {
        require(_to != address(0), "Invalid recipient address");
        require(_value > 0, "Mint amount must be greater than 0");

        totalSupply += _value;
        balanceOf[_to] += _value;
        emit Transfer(address(0), _to, _value); // 增发事件：从0地址转至_to
        return true;
    }


    // ========== ERC20标准事件 ==========
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}
