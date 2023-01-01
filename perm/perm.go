package perm

import "os"

// This small package exists mainly for documentation.
// It's a specific subject with many cases that are hard
// to cover in a single constant.

/*
+-----+---+--------------------------+
| rwx | 7 | Read, write and execute  |
| rw- | 6 | Read, write              |
| r-x | 5 | Read, and execute        |
| r-- | 4 | Read,                    |
| -wx | 3 | Write and execute        |
| -w- | 2 | Write                    |
| --x | 1 | Execute                  |
| --- | 0 | no permissions           |
+------------------------------------+

+------------+------+-------+
| Permission | Octal| Field |
+------------+------+-------+
| rwx------  | 0700 | User  |
| ---rwx---  | 0070 | Group |
| ------rwx  | 0007 | Other |
+------------+------+-------+
*/

// Note: The underlying types for these constants are os.FileMode.
// Might not need these custom types if the FileMode constants in
// fs.go are applicable.

// Read, write, execute permission bitmask.
const RWX os.FileMode = 0750

// Read and write permission bitmask.
const RW os.FileMode = 0660
