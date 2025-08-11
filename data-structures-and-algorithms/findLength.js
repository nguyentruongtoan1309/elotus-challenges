/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 */
var findLength = function (nums1, nums2) {
  const m = nums1.length;
  const n = nums2.length;

  let prev = Array(n + 1).fill(0);
  let maxLength = 0;

  for (let i = 1; i <= m; i++) {
    const curr = Array(n + 1).fill(0);

    for (let j = 1; j <= n; j++) {
      if (nums1[i - 1] === nums2[j - 1]) {
        curr[j] = prev[j - 1] + 1;
        maxLength = Math.max(maxLength, curr[j]);
      }
    }

    prev = curr;
  }

  return maxLength;
};

// Time Complexity: O(m × n) where m, n are lengths of the arrays
// Space Complexity: O(min(m,n))

console.log(findLength([1, 2, 3, 2, 1], [2, 3, 2, 1, 4, 7]));
