# Run this by piping some processed logs to it:
"""
To get stats on outcome latency per processor:

```bash
grep -h "tracking processor latency" logs/Test_CCIPBatching_MaxBatchSizeEVM_17424836*.log | \
    grep -h "\"method\":\"outcome\"" | \
    grep -h "\"plugin\":\"Commit\"" | \
    jq -r '"\(.processor):\(.latency)"' | python3 ~/stats.py
```

Make sure to replace the logs file with the correct one on your local. This is printed out
at the beginning of the test run.

```bash
cd integration-tests/smoke/ccip
go test -v -timeout 45m -run "Test_CCIPBatching_MaxBatchSizeEVM" .
```
"""
import sys
import re
import numpy as np
from collections import defaultdict

def parse_duration(duration: str) -> float:
    match = re.match(r"([0-9.]+)(s|ms|µs)", duration)
    if not match:
        raise ValueError(f"Invalid duration format: {duration}")
    
    value, unit = match.groups()
    value = float(value)
    
    if unit == "s":
        return value
    elif unit == "ms":
        return value / 1000
    elif unit == "µs":
        return value / 1_000_000
    
    raise ValueError(f"Unknown unit: {unit}")

# Read input from stdin
duration_map = defaultdict(list)

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    
    component, duration = line.split(":")
    duration_map[component].append(parse_duration(duration))

# Compute statistics for each component
for component, durations in duration_map.items():
    durations.sort()
    min_value = durations[0]
    max_value = durations[-1]
    median_value = np.median(durations)
    p95_value = np.percentile(durations, 95)
    
    print(f"{component}:")
    print(f"  Min: {min_value:.9f} s")
    print(f"  Max: {max_value:.9f} s")
    print(f"  Median: {median_value:.9f} s")
    print(f"  P95: {p95_value:.9f} s")
    print()
