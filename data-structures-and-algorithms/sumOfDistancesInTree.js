/**
 * @param {number} n
 * @param {number[][]} edges
 * @return {number[]}
 */
var sumOfDistancesInTree = function (n, edges) {
  // Handle edge case
  if (n === 1) return [0];

  // Build adjacency list
  const graph = Array(n)
    .fill(null)
    .map(() => []);
  for (const [u, v] of edges) {
    graph[u].push(v);
    graph[v].push(u);
  }

  const subtreeSize = Array(n).fill(0);
  const answer = Array(n).fill(0);

  // First DFS: Calculate subtree sizes and answer for root (node 0)
  function dfs1(node, parent) {
    subtreeSize[node] = 1;

    for (const child of graph[node]) {
      if (child === parent) continue;

      dfs1(child, node);
      subtreeSize[node] += subtreeSize[child];
      // Add distances: each node in child's subtree is 1 step farther from current node
      answer[0] += answer[child] + subtreeSize[child];
    }
  }

  // Second DFS: Re-root the tree to calculate answers for all nodes
  function dfs2(node, parent) {
    for (const child of graph[node]) {
      if (child === parent) continue;

      // When we move root from 'node' to 'child':
      // - Nodes in child's subtree get 1 step closer (subtract subtreeSize[child])
      // - Nodes outside child's subtree get 1 step farther (add (n - subtreeSize[child]))
      answer[child] =
        answer[node] - subtreeSize[child] + (n - subtreeSize[child]);

      dfs2(child, node);
    }
  }

  dfs1(0, -1);
  dfs2(0, -1);

  return answer;
};

// Time Complexity: O(n)
// Space Complexity: O(n)

console.log(
  sumOfDistancesInTree(6, [
    [0, 1],
    [0, 2],
    [2, 3],
    [2, 4],
    [2, 5],
  ])
);
