# IPFS-Watch-Files-Tool  

## Need To Do  

- 实时监控指定的目录；  
- 目录中的任何改动，包含对文件的增删改操作，都需要同步到ipfs中；  
- 有单元测试代码；  
-  
## Usage  

```
$ go run main.go 
ipfs watch dir is:  /home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test
INFO[0000] start                                        
INFO[0010] Current time: 2021-10-31 20:38:16.090045053 +0800 CST m=+10.005397790 
watch file name is:  /home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test/t1.txt
INFO[0010] /home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test/t1.txt cid is: QmUgAgTVxq7UeY3Tbumz72fBsSvkUnveEgEkWvVquEvJVV 
watch file name is:  /home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test/t2.txt
INFO[0010] /home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test/t2.txt cid is: QmaRGe7bVmVaLmxbrMiVNXqW4pRNNp3xq7hFtyRKA3mtJL 
Done!


$ ipfs cat QmUgAgTVxq7UeY3Tbumz72fBsSvkUnveEgEkWvVquEvJVV
hello1
$ ipfs cat QmaRGe7bVmVaLmxbrMiVNXqW4pRNNp3xq7hFtyRKA3mtJL
world
```