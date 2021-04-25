using System.Collections.Generic;
using System.Text;

namespace q38 {
    public class Solution {
    public string CountAndSay(int n) {
        var prev = new List<int>();
        prev.Add(1);
        for(int i=1; i<n; i++) {
            var next = say(prev);
            prev = next;
        }
        var output = new StringBuilder();
        foreach(int x in prev) {
            output.Append(x);
        }
        return output.ToString();
    }
    private List<int> say(in List<int> input) {
        var output = new List<int>();
        for(int i = 0; i < input.Count; ) {
            int j = i+1;
            for(; j < input.Count && input[j] == input[i]; j++) {
            }
            int cnt = j-i;
            output.Add(cnt);
            output.Add(input[i]);
            i = j;
        }
        return output;
    }
}
}