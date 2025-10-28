// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/* 
    合并两个有序数组 (Merge Sorted Array)
        题目描述：将两个有序数组合并为一个有序数组。
*/

library utils {
    function merge_array(
        uint256[] memory nums1,
        uint256[] memory nums2
    ) public pure returns (uint256[] memory) {
        uint256 len1 = nums1.length;
        uint256 index1 = 0;
        uint256 len2 = nums2.length;
        uint256 index2 = 0;
        uint256[] memory ret = new uint256[](len1 + len2);
        for (uint256 i = 0; i < len1 + len2; ++i) {
            if (index1 < len1 && index2 < len2) {
                if (nums1[index1] < nums2[index2]) {
                    ret[i] = nums1[index1];
                    index1++;
                } else {
                    ret[i] = nums2[index2];
                    index2++;
                }
            } else if (index1 < len1) {
                ret[i] = nums1[index1];
                index1++;
            } else if (index2 < len2) {
                ret[i] = nums2[index2];
                index2++;
            }
        }

        return ret;
    }
}
