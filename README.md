# Diff

A library for diffing, plus a command to produce unified diffs (mimics diff -u)

# CLI usage
```
go run ./cmd/diff FILE1 FILE2
```

For example:
```
go run ./cmd/diff <(printf "a\nb\nHello") <(printf "b\nc\nHello\n")
```

Produces:
```
--- /dev/fd/63  2025-01-15 08:47:53.573736556 -0800
+++ /dev/fd/62  2025-01-15 08:47:53.573736556 -0800
@@ -1,3 +0,4 @@
-a
 b
-Hello
\ No newline at end of file
+c
+Hello
```

## Notes

The diff algorithm is O(nm) where n and m are the lines in file1 and file2, respectively.

The output is intended to be as close as possible with GNU diff's unified output, but
there are probably still cases where output differs. The most common source of differences
from GNU diff is putting deletes and inserts in a different order. The current implementation
will prefer to delete before performing additions.
