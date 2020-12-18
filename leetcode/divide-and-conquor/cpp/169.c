int majorityElement(int* nums, int numsSize){
    int count = 0;
    int result = 0;
    for(int i=0; i<numsSize; i++) {
        int v = nums[i];
        if (count == 0) {
            result = v;
            count = 1;
        } else if (result == v) {
            count++;
        } else {
            count--;
        }
    }
    return result;
}