// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/*
    反转字符串 (Reverse String)
        题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
*/

library utils {
    // 仅处理 ASCII 编码字符串。
    function reverse_ascii_string(
        string memory istr
    ) public pure returns (string memory ostr) {
        bytes memory tmp = bytes(istr);
        uint256 len = tmp.length;
        for (uint256 i = 0; i < len / 2; i++) {
            bytes1 t = tmp[len - 1 - i];
            tmp[len - 1 - i] = tmp[i];
            tmp[i] = t;
        }
        return string(tmp);
    }
}
