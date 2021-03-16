import sys
a = []

for l in sys.stdin:
    a.append(l.strip("\\n"))

def listParser(s):
    return s.replace("[", "").replace("]", "").split(",")

def main(nums, num):
    return 9

is_success = main(*a[:-1]) == a[-1]
success_flag = "success" if is_success else "fail"

print(f"{success_flag}")
