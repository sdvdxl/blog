---
title: go读取文件
tags:
  - go
  - golang
category: golang
abbrlink: 44633
date: 2016-03-13 11:42:18
---
# 使用File
    不多说，直接上代码

```go
func readUseFile() {
	file, err := os.Open("f:/file.txt")
	handleError(err)

	defer file.Close()

	buf := make([]byte, 512)
	for {
		n, err := file.Read(buf)

		//1
		//		if err != nil && err == io.EOF {
		//			break
		//		}

		//2
		//		if n == 0 {
		//			break
		//		}

		//3
		if n == 0 && err != nil && err == io.EOF {
			break
		}

		fmt.Print(n)
		fmt.Print(string(buf))
	}

	handleError(err)

}
 ```
    可以看到，file本身具有读取文件内容的函数，入参事一个切片，是数据的缓冲区，出参第一个是实际读取的大小，第二个是读取过程中发生的错误。如果有数据且读取成功，则n>0,如果恰好读到文件末尾，则n=0。如果读取过程中有错误发生，则err不为nil，如果读取正常且读到了文件末尾，则err为io.EOF。
    读取过程中有三种方法可以跳出死循环。第一种方法是判断err状态，如果不为nil且是io.EOF，则已经读取完毕；第二种方法是判断实际读取的数量，如果读取的量为0，则认为已经读取结束。第三种方式是上面两种的结合，这种判断要比上面两种中仁和一种都要保险，缺点就是罗嗦点。
