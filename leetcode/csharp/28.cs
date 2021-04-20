namespace q28
{
public class Solution
{
    // 使用Rabin-Karp算法滚动匹配字符串，移植自go bytealg.go
    public int StrStr(string haystack, string needle) {
        if (needle.Length == 0) {
            return 0;
        }
        if (haystack.Length <= needle.Length) {
            return haystack == needle ? 0 : -1;
        }
        var (hashseq, pow) = HashStr(needle);
        var n = needle.Length;
        uint h = 0;
        for (int i = 0; i < n; i++) {
            h = h * PrimeRK + (uint)haystack[i];
        }

        if (hashseq == h && isSubStr(haystack, needle, 0)) {
            return 0;
        }
        for (int i = n; i < haystack.Length;) {
            h *= PrimeRK;
            h += (uint)haystack[i];
            h -= pow * (uint)haystack[i - n];
            i++;
            if (h == hashseq && isSubStr(haystack, needle, i - n)) {
                return i - n;
            }
        }
        return -1;
    }
    public static bool isSubStr(string s, string sub, int start) {
        for (int i = 0; i < sub.Length; i++) {
            if (s[start++] != sub[i]) {
                return false;
            }
        }
        return true;
    }
    public static uint PrimeRK = 16777619;
    public static (uint, uint) HashStr(string s) {
        var hash = (uint)0;
        for (int i = 0; i < s.Length; i++) {
            hash = hash * PrimeRK + (uint)s[i];
        }
        var pow = (uint)1;
        var sq = PrimeRK;
        for (int i = s.Length; i > 0; i >>= 1) {
            if ((i & 1) != 0) {
                pow *= sq;
            }
            sq *= sq;
        }
        return (hash, pow);
    }
}

}