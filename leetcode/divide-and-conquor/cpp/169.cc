#include <vector>
#include <iterator>
using std::vector;

class Solution {
public:
    int majorityElement(vector<int>& nums) {
        int count = 0;
        int result = 0;
        for(int v : nums) {
            if (count == 0) {
                count = 1;
                result = v;
            } else {
                v == result ? count++ : count--;
            }
        }
        return result;
    }
};
