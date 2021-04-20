namespace q27 {
public class Solution {
    public int RemoveElement(int[] nums, int val) {
        int cnt = 0;
        for(int i = 0; i < nums.Length; i++) {
            if(nums[i] == val) {
                continue;
            }
            nums[cnt++] = nums[i];
        }
        return cnt;
    }
}

}

