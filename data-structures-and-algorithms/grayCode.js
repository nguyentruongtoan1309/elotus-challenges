/**
 * @param {number} n
 * @return {number[]}
 */
var grayCode = function (n) {
  let result = [0];
  // Iterate through each bit position
  for (let i = 0; i < n; i++) {
    // Calculate the most significant bit (2^i)
    const msb = 1 << i;

    // Add reversed sequence with msb set
    for (let j = result.length - 1; j >= 0; j--) {
      result.push(result[j] + msb);
    }
  }

  return result;
};

// Time Complexity: O(2^n)
// Space Complexity: O(2^n)

console.log("n = 1:", grayCode(1));
console.log("n = 2:", grayCode(2));
console.log("n = 3:", grayCode(3));
