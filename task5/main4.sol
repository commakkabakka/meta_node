// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/* 
    用 solidity 实现罗马数字转数整数
        题目描述在 https://leetcode.cn/problems/integer-to-roman/description/

        七个不同的符号代表罗马数字，其值如下：

        符号	值
        I	1
        V	5
        X	10
        L	50
        C	100
        D	500
        M	1000
        罗马数字是通过添加从最高到最低的小数位值的转换而形成的。将小数位值转换为罗马数字有以下规则：

        如果该值不是以 4 或 9 开头，请选择可以从输入中减去的最大值的符号，将该符号附加到结果，减去其值，然后将其余部分转换为罗马数字。
        如果该值以 4 或 9 开头，使用 减法形式，表示从以下符号中减去一个符号，例如 4 是 5 (V) 减 1 (I): IV ，9 是 10 (X) 减 1 (I)：IX。仅使用以下减法形式：4 (IV)，9 (IX)，40 (XL)，90 (XC)，400 (CD) 和 900 (CM)。
        只有 10 的次方（I, X, C, M）最多可以连续附加 3 次以代表 10 的倍数。你不能多次附加 5 (V)，50 (L) 或 500 (D)。如果需要将符号附加4次，请使用 减法形式。
        给定一个整数，将其转换为罗马数字。
*/

library utils {
    function to_roman(uint256 num) public pure returns (string memory) {
        uint256[13] memory values = [
            uint256(1000),
            uint256(900),
            uint256(500),
            uint256(400),
            uint256(100),
            uint256(90),
            uint256(50),
            uint256(40),
            uint256(10),
            uint256(9),
            uint256(5),
            uint256(4),
            uint256(1)
        ];
        string[13] memory symbols = [
            "M",
            "CM",
            "D",
            "CD",
            "C",
            "XC",
            "L",
            "XL",
            "X",
            "IX",
            "V",
            "IV",
            "I"
        ];
        string memory result = "";
        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                result = string(abi.encodePacked(result, symbols[i]));
                num -= values[i];
            }
        }
        return result;
    }
}
