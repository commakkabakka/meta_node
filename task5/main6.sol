// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/*
    二分查找 (Binary Search)
        题目描述：在一个有序数组中查找目标值。
*/

library utils {
    function BinarySearch(
        uint256[] memory arr,
        uint256 item
    ) public pure returns (int256) {
        uint256 left = 0;
        uint256 right = arr.length - 1;
        while (left <= right) {
            uint256 mid = (left + right) / 2;
            if (item == arr[mid]) {
                return int256(mid);
            } else if (item < arr[mid]) {
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        }
        return -1;
    }
}
