// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import "hardhat/console.sol";
// 任务一功能合集
contract task1Funcs {
    /*
    1. 反转字符串 (Reverse String)
- 题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
    */
    function reverseString(
        string memory _str
    ) public pure returns (string memory) {
        bytes memory strBytes = bytes(_str);
        bytes memory resultBytes = new bytes(strBytes.length);

        for (uint256 i = 0; i < strBytes.length; i++) {
            resultBytes[i] = strBytes[strBytes.length - 1 - i];
        }
        return string(resultBytes);
    }
    /*
2.   用 solidity 实现整数转罗马数字
- 题目描述在 https://leetcode.cn/problems/roman-to-integer/description/3.
    */
    error NumberOverRomanRange(uint256 _num); // 自定义的带参数的error
    function intToRoman(uint256 _num) public pure returns (string memory) {
        // 基于标准罗马数字 限制参数范围
        if (_num < 1 || _num > 3999) {
            revert NumberOverRomanRange(_num);
        }
        //require(_num>=1 && _num<=3999,"number must be between 1 and 3999");
        uint256[] memory values = new uint256[](13);
        values[0] = 1000;
        values[1] = 900;
        values[2] = 500;
        values[3] = 400;
        values[4] = 100;
        values[5] = 90;
        values[6] = 50;
        values[7] = 40;
        values[8] = 10;
        values[9] = 9;
        values[10] = 5;
        values[11] = 4;
        values[12] = 1;

        string[] memory symbols = new string[](13);
        symbols[0] = "M";
        symbols[1] = "CM";
        symbols[2] = "D";
        symbols[3] = "CD";
        symbols[4] = "C";
        symbols[5] = "XC";
        symbols[6] = "L";
        symbols[7] = "XL";
        symbols[8] = "X";
        symbols[9] = "IX";
        symbols[10] = "V";
        symbols[11] = "IV";
        symbols[12] = "I";

        string memory result = "";
        for (uint256 i = 0; i < values.length; i++) {
            while (_num >= values[i]) {
                result = string.concat(result, symbols[i]);
                _num -= values[i];
            }
            if (_num == 0) break;
        }
        return result;
    }

    /*
3.   用 solidity 实现罗马数字转数整数
- 题目描述在 https://leetcode.cn/problems/integer-to-roman/description/
    */
    struct RomanPair {
        bytes1 symbol;
        uint256 value;
    }

    function getValue(
        RomanPair[] memory romanPairs,
        bytes1 symbol
    ) public pure returns (uint256) {
        console.log("0-1");
        for (uint256 i = 0; i < romanPairs.length; i++) {
            if (romanPairs[i].symbol == symbol) {
                return romanPairs[i].value;
            }
        }
        return 0;
    }
    function romanToInt(string memory _roman) public pure returns (uint256) {
        RomanPair[] memory romanPairs = new RomanPair[](7);
        romanPairs[0] = RomanPair(bytes1("I"), 1);
        romanPairs[1] = RomanPair(bytes1("V"), 5);
        romanPairs[2] = RomanPair(bytes1("X"), 10);
        romanPairs[3] = RomanPair(bytes1("L"), 50);
        romanPairs[4] = RomanPair(bytes1("C"), 100);
        romanPairs[5] = RomanPair(bytes1("D"), 500);
        romanPairs[6] = RomanPair(bytes1("M"), 1000);
        bytes memory romanBytes = bytes(_roman);
        uint256 result = 1000;
        for (uint256 i = 0; i < romanBytes.length; i++) {
            uint256 curValue = getValue(romanPairs,romanBytes[i]);

            if (i < romanBytes.length - 1) {
                uint256 aftValue = getValue(romanPairs,romanBytes[i + 1]);
                if (curValue < aftValue) {
                    console.log("1-curValue:",curValue);
                    console.log("1-aftValue:",aftValue);
                    console.log("1-result:",result);
                    result -= curValue;
                    console.log("1-result:",result);
                }else{
                    console.log("2-curValue:",curValue);
                    console.log("2-aftValue:",aftValue);
                    console.log("2-result:",result);
                result += curValue;
                console.log("2-result:",result);
                }
            } else {
                console.log("3-curValue:",curValue);
                    console.log("3-result:",result);
                result += curValue;
                console.log("3-result:",result);
            }
        }
        console.log("4-result:",result);
        result -= 1000;
        console.log("5-result:",result);
        return result;
    }
    /*
4.   合并两个有序数组 (Merge Sorted Array)
- 题目描述：将两个有序数组合并为一个有序数组。
    */
    function mergeSortedArrays(
        uint256[] memory _arr1,
        uint256[] memory _arr2
    ) public pure returns (uint256[] memory) {
        uint256[] memory result = new uint256[](_arr1.length + _arr2.length);
        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;

        while (i < _arr1.length && j < _arr2.length) {
            if (_arr1[i] <= _arr2[j]) {
                result[k] = _arr1[i];
                i++;
            } else {
                result[k] = _arr2[j];
                j++;
            }
            k++;
        }

        // 处理剩余元素
        while (i < _arr1.length) {
            result[k] = _arr1[i];
            i++;
            k++;
        }
        while (j < _arr2.length) {
            result[k] = _arr2[j];
            j++;
            k++;
        }

        return result;
    }
    /*
5.   二分查找 (Binary Search)
- 题目描述：在一个有序数组中查找目标值。
    */
    function binarySearch(
        uint256[] memory _sortedArr,
        uint256 _target
    ) public pure returns (int256) {
        uint256 left = 0;
        uint256 right = _sortedArr.length - 1;

        while (left <= right) {
            uint256 mid = left + (right - left) / 2; // 避免溢出
            if (_sortedArr[mid] == _target) {
                return int256(mid); // 找到，返回索引
            } else if (_sortedArr[mid] < _target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        return -1; // 未找到
    }
}
