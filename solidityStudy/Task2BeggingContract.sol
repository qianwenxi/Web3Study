// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract {
    // 1. 核心基础变量
    address public immutable owner; // 合约所有者（收款汇总方）
    mapping(address => uint256) public donations; // 记录每个地址的捐赠金额
    uint256 public totalDonations; // 总捐赠金额（辅助核对）

    // 2. 捐赠时间限制相关变量（可自定义时间段）
    uint256 public immutable donationStartTime; // 捐赠开始时间（时间戳）
    uint256 public immutable donationEndTime; // 捐赠结束时间（时间戳）

    // 3. 优化：捐赠排行榜核心变量（存储TOP3地址和金额）
    address[] private topDonors; // 存储前3捐赠者地址（最多3个）
    uint256[] private topDonations; // 存储前3捐赠者对应的金额

    // 部署合约时初始化：绑定所有者+设置捐赠时间段
    // 参数说明：_startTime（开始时间戳），_endTime（结束时间戳）
    error InvalidStarttimeEndtime(uint256 _startTime,uint256 _endTime);
    constructor(uint256 _startTime, uint256 _endTime) {
        owner = msg.sender;
        // 校验：开始时间必须早于结束时间，避免无效设置
        if(_startTime > _endTime){
            revert InvalidStarttimeEndtime(_startTime, _endTime);
        }
        donationStartTime = _startTime;
        donationEndTime = _endTime;
    }

    // 核心功能1：捐赠（新增时间限制校验）
    error OutofDonateTime(uint256 donationStartTime,uint256 donationEndTime);
    error DonateZero();
    function donate() public payable {
        // 校验1：当前时间在捐赠时间段内（不在则禁止捐赠）
        if(!isDonationAvailable()){
            revert OutofDonateTime(donationStartTime, donationEndTime);
        }
        // 校验2：禁止空转账（至少捐一点测试网ETH）
        if(msg.value == 0){
            revert DonateZero();
        }

        // 更新当前地址捐赠金额
        donations[msg.sender] += msg.value;
        // 更新总捐赠金额
        totalDonations += msg.value;
        // 同步更新捐赠排行榜（核心新增逻辑）
        updateTopDonors(msg.sender, donations[msg.sender]);
    }

    // 核心新增功能2：更新捐赠排行榜（内部函数，自动触发，无需手动调用）
    function updateTopDonors(address donor, uint256 amount) private {
        // 1. 先判断捐赠者是否已在TOP3中，若在则更新金额
        for (uint256 i = 0; i < topDonors.length; i++) {
            if (topDonors[i] == donor) {
                topDonations[i] = amount;
                sortTopDonors(); // 更新后重新排序，保证排名准确
                return;
            }
        }

        // 2. 若不在TOP3，且TOP3未满（少于3人），直接加入
        if (topDonors.length < 3) {
            topDonors.push(donor);
            topDonations.push(amount);
            sortTopDonors(); // 加入后排序
            return;
        }

        // 3. 若TOP3已满，且当前捐赠金额超过第三名，替换第三名并排序
        if (amount > topDonations[2]) {
            topDonors[2] = donor;
            topDonations[2] = amount;
            sortTopDonors(); // 替换后排序
        }
    }

    // 辅助函数：对TOP3排行榜排序（从高到低，金额降序）
    function sortTopDonors() private {
        // 简单排序逻辑（适配3个元素，高效不冗余）
        for (uint256 i = 0; i < topDonors.length; i++) {
            for (uint256 j = i + 1; j < topDonors.length; j++) {
                if (topDonations[j] > topDonations[i]) {
                    // 交换金额
                    (topDonations[i], topDonations[j]) = (topDonations[j], topDonations[i]);
                    // 同步交换地址（保持地址和金额对应）
                    (topDonors[i], topDonors[j]) = (topDonors[j], topDonors[i]);
                }
            }
        }
    }

    // 核心新增功能3：查询捐赠TOP3排行榜（外部可调用，任何人都能查）
    // 返回值：前3捐赠者地址数组、对应金额数组（顺序：第1名→第2名→第3名）
    function getTop3Donors() public view returns (address[] memory, uint256[] memory) {
        return (topDonors, topDonations);
    }

    // 原有功能：所有者提取全部捐赠资金（无修改，保留）
    error WithdrawAllNotWoner(address owner);
    function withdrawAll() public {
        if(msg.sender != owner){
            revert WithdrawAllNotWoner(msg.sender);
        }
        // 把合约所有余额转给所有者
        payable(owner).transfer(address(this).balance);
    }

    // 辅助功能：查询当前是否可捐赠（方便测试核对）
    function isDonationAvailable() public view returns (bool) {
        return block.timestamp >= donationStartTime && block.timestamp <= donationEndTime;
    }
}