// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

abstract contract task2ERC20 {
    mapping(address account => uint256) private _balances;
    mapping(address account => mapping(address spender => uint256)) private _allowances;

    uint256 private _totalSupply;
    string private _name;
    string private _symbol;
    // 代币合约构造函数，明确代币名称和符号
    constructor(string memory name_, string memory symbol_) {
        _name = name_;
        _symbol = symbol_;
    }
    // 代币名称
    function name() public view virtual returns (string memory) {
        return _name;
    }
    // 代币符号
    function symbol() public view virtual returns (string memory) {
        return _symbol;
    }
    // 小数位数
    function decimals() public view virtual returns (uint8) {
        return 18;
    }
    // 返回代币总供给
    function totalSupply() public view virtual returns (uint256) {
        return _totalSupply;
    }
    // 返回账户代币余额
    function balanceOf(address account) public view virtual returns (uint256) {
        return _balances[account];
    }
    // 转账函数
    function transfer(address to, uint256 value) public virtual returns (bool) {
        address owner = msg.sender;
        _transfer(owner, to, value);
        return true;
    }
    // 返回param1给param2的授权额度
    function allowance(address owner, address spender) public view virtual returns (uint256) {
        return _allowances[owner][spender];
    }
    // 授权额度
    function approve(address spender, uint256 value) public virtual returns (bool) {
        address owner = msg.sender;
        _approve(owner, spender, value);
        return true;
    }
    // 授权转账,有param1账户授权额度的账户，操作param1账户转出param3给账户param2
    function transferFrom(address from, address to, uint256 value) public virtual returns (bool) {
        address spender = msg.sender;
        _spendAllowance(from, spender, value);
        _transfer(from, to, value);
        return true;
    }
    // 内部转账函数
    error ERC20InvalidSender(address sender);
    error ERC20InvalidReceiver(address receiver);
    function _transfer(address from, address to, uint256 value) internal {
        if (from == address(0)) {
            revert ERC20InvalidSender(address(0));
        }
        if (to == address(0)) {
            revert ERC20InvalidReceiver(address(0));
        }
        _update(from, to, value);
    }
    // 转账明细函数-支持铸币和销毁（通过零地址实现）
    error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed);
    event Transfer(address indexed from, address indexed to, uint256 value);
    function _update(address from, address to, uint256 value) internal virtual {
        if (from == address(0)) {
            // Overflow check required: The rest of the code assumes that totalSupply never overflows
            _totalSupply += value;
        } else {
            uint256 fromBalance = _balances[from];
            if (fromBalance < value) {
                revert ERC20InsufficientBalance(from, fromBalance, value);
            }
            unchecked {
                // Overflow not possible: value <= fromBalance <= totalSupply.
                _balances[from] = fromBalance - value;
            }
        }

        if (to == address(0)) {
            unchecked {
                // Overflow not possible: value <= totalSupply or value <= fromBalance <= totalSupply.
                _totalSupply -= value;
            }
        } else {
            unchecked {
                // Overflow not possible: balance + value is at most totalSupply, which we know fits into a uint256.
                _balances[to] += value;
            }
        }

        emit Transfer(from, to, value);
    }
    function _mint(address account, uint256 value) internal {
        if (account == address(0)) {
            revert ERC20InvalidReceiver(address(0));
        }
        _update(address(0), account, value);
    }
    function _burn(address account, uint256 value) internal {
        if (account == address(0)) {
            revert ERC20InvalidSender(address(0));
        }
        _update(account, address(0), value);
    }
    function _approve(address owner, address spender, uint256 value) internal {
        _approve(owner, spender, value, true);
    }
    error ERC20InvalidApprover(address approver);
    error ERC20InvalidSpender(address spender);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual {
        if (owner == address(0)) {
            revert ERC20InvalidApprover(address(0));
        }
        if (spender == address(0)) {
            revert ERC20InvalidSpender(address(0));
        }
        _allowances[owner][spender] = value;
        if (emitEvent) {
            emit Approval(owner, spender, value);
        }
    }
    error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed);
    function _spendAllowance(address owner, address spender, uint256 value) internal virtual {
        uint256 currentAllowance = allowance(owner, spender);
        // 无限授权无需扣减额度，无限授权额度为type(uint256).max
        if (currentAllowance < type(uint256).max) {
            if (currentAllowance < value) {
                revert ERC20InsufficientAllowance(spender, currentAllowance, value);
            }
            unchecked {
                _approve(owner, spender, currentAllowance - value, false);
            }
        }
    }
}
