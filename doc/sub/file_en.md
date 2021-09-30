## File Permission Keywords

File permissions vary with security mechanisms. To shield differences for users, secPaver extracts and encapsulates some common file permissions.

| Permission  |              Description              |             Note              |
| :---------: | :-----------------------------------: | :---------------------------: |
|   create    |           Create new file。           |                               |
|    read     |          Read file contents.          |                               |
|    write    |           Write to a file.            |                               |
|   append    | Write to a file opened with O_APPEND. |                               |
| rename/link |   Rename/Create hard link to file.    |                               |
|   remove    |            Delete a file.             |                               |
|    lock     |       Set and unset file locks.       |                               |
|     map     |          Memory map a file.           |                               |
|    exec     |               Execute.                |                               |
|   search    |            Search access.             |   Only valid for directory.   |
|    ioctl    |   IO control system call requests.    |                               |
|    mount    |              Mount file.              | Only valid for SELinux policy |
|   mounton   |          Use as mount point.          | Only valid for SELinux policy |

### Note

- The file permissions of the SELinux mechanism are complex and cannot be completely overwritten. If you need to use the file permissions, you can directly enter the permissions in the actions column. The secPaver checks the extended permissions. If the permissions are valid, they are added to SELinux rules.
- For ease of use, the main application has the search permission on directory files (except files set to private) by default during policy generation.