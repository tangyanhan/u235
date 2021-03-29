/**
 * Union Find
 */
class Solution {
    public boolean possibleBipartition(int N, int[][] dislikes) {
        UnionFind uf = new UnionFind(N * 2 + 1);
        for (int i = 0; i < dislikes.length; i++) {
            int x = uf.find(dislikes[i][0]); //查找自己的根节点
            int y = uf.find(dislikes[i][1]); //不喜欢的人的根节点
            int a = uf.find(dislikes[i][0] + N); //查找自己不喜欢的人群根节点
            int b = uf.find(dislikes[i][1] + N);  // 自己不喜欢的人不喜欢的人群节点
            if(x == y) {
                return false; //发现这俩人已经在一组，GG
            }
            else{
                uf.union(y, a); // Union persons that are disliked
                uf.union(b, x);
            }
        }
        return true;
    }
    
    private class UnionFind {
        int roots;
        int[] parent;
        
        public UnionFind(int size) {
            this.roots = size;
            this.parent = new int[size];
            for (int i = 0; i < size; i++) {
                parent[i] = i;
            }
        }
        
        void union(int p, int q) {
            int rootP = parent[p];
            int rootQ = parent[q];
            
            if (rootP != rootQ) {
                parent[rootP] = rootQ;
                roots--;
            }
            
            return;
        }
        
        int find(int p) {
            while (p != parent[p]) {
                p = parent[parent[p]];
            }
            
            return p;
        }
        
        boolean isConnected(int p, int q) {
            return find(p) == find(q);
        }
    }
}

// 链接：https://leetcode-cn.com/problems/possible-bipartition/solution/javaunion-findgai-jin-guo-de-bing-cha-ji-zuo-fa-by/