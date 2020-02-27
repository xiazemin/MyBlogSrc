---
title: LimitReader
layout: post
category: golang
author: 夏泽民
---
https://github.com/dragonflyoss/Dragonfly/blob/v0.4.3/pkg/ratelimiter/ratelimiter.go
<!-- more -->
package main

import (
   "fmt"
   "os"
   "io/ioutil"
   "io"
)

func main() {
   file, err := os.OpenFile("test.txt", os.O_RDWR, 777) //make一个reader
   if err != nil {
      panic(err)
   }
   r := io.NewSectionReader(file, 4, 10) //规定范围
   buf, err := ioutil.ReadAll(r)
   fmt.Println(string(buf))

   lr := io.LimitReader(file, 12) //规定最大能读取到哪里
   buf, err = ioutil.ReadAll(lr)
   fmt.Println(string(buf))
}
最终结果输出

 is io-tes
this is io-t

io包的使用
功能：io包提供了对I/O原语的基本接口。本包的基本任务是包装这些原语已有的实现（如os包里的原语），使之成为共享的公共接口，这些公共接口抽象出了泛用的函数并附加了一些相关的原语的操作

１、Reader 接口

type Reader interface {
     Read(p []byte) (n int, err error)
}
Read方法读取len(p)字节数据写入p。它返回写入的字节数和遇到的任何错误。即使Read方法返回值n < len(p)，本方法在被调用时仍可能使用p的全部长度作为暂存空间。如果有部分可用数据，但不够len(p)字节，Read按惯例会返回可以读取到的数据，而不是等待更多数据。

当Read在读取n > 0个字节后遭遇错误或者到达文件结尾时，会返回读取的字节数。它可能会在该次调用返回一个非nil的错误，或者在下一次调用时返回0和该错误。一个常见的例子，Reader接口会在输入流的结尾返回非0的字节数，返回值err == EOF或err == nil。但不管怎样，下一次Read调用必然返回(0, EOF)。调用者应该总是先处理读取的n > 0字节再处理错误值。这么做可以正确的处理发生在读取部分数据后的I/O错误，也能正确处理EOF事件。

如果Read的某个实现返回0字节数和nil错误值，表示被阻碍；调用者应该将这种情况视为未进行操作

２、Writer 接口

type Writer interface {
    Write(p []byte) (n int, err error)
}
Writer接口用于包装基本的写入方法。

Write方法len(p) 字节数据从p写入底层的数据流。它会返回写入的字节数(0 <= n <= len(p))和遇到的任何导致写入提取结束的错误。Write必须返回非nil的错误，如果它返回的 n < len(p)。Write不能修改切片p中的数据，即使临时修改也不行

３、Seeker 接口

type Seeker interface {
        Seek(offset int64, whence int) (int64, error)
}
Seek方法设定下一次读写的位置：偏移量为offset，校准点由whence确定：0表示相对于文件起始；1表示相对于当前位置；2表示相对于文件结尾。Seek方法返回新的位置以及可能遇到的错误

４、Closer接口

type Closer interface {
    Close() error
}
Closer接口用于包装基本的关闭方法。
在第一次调用之后再次被调用时，Close方法的的行为是未定义的。某些实现可能会说明他们自己的行为。

５、ReaderAt 和 WriterAt 接口

type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
}
ReadAt从底层输入流的偏移量off位置读取len(p)字节数据写入p， 它返回读取的字节数(0 <= n <= len(p))和遇到的任何错误。当ReadAt方法返回值n < len(p)时，它会返回一个非nil的错误来说明为啥没有读取更多的字节。在这方面，ReadAt是比Read要严格的。即使ReadAt方法返回值 n < len(p)，它在被调用时仍可能使用p的全部长度作为暂存空间。如果有部分可用数据，但不够len(p)字节，ReadAt会阻塞直到获取len(p)个字节数据或者遇到错误。在这方面，ReadAt和Read是不同的。如果ReadAt返回时到达输入流的结尾，而返回值n == len(p)，其返回值err既可以是EOF也可以是nil

type WriterAt interface {
    WriteAt(p []byte, off int64) (n int, err error)
}
WriteAt将p全部len(p)字节数据写入底层数据流的偏移量off位置。它返回写入的字节数(0 <= n <= len(p))和遇到的任何导致写入提前中止的错误。当其返回值n < len(p)时，WriteAt必须放哪会一个非nil的错误。
如果WriteAt写入的对象是某个有偏移量的底层输出流（的Writer包装），WriteAt方法既不应影响底层的偏移量，也不应被底层的偏移量影响

６、ReaderFrom 和 WriterTo 接口

type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
}

type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
}
７、ByteReader 和 ByteWriter

type ByteReader interface {
    ReadByte() (c byte, err error)
}

type ByteWriter interface {
    WriteByte(c byte) error
}
８、ByteScanner、RuneReader 和 RuneScanner

type ByteScanner interface {
    ByteReader
    UnreadByte() error
}

type RuneReader interface {
    ReadRune() (r rune, size int, err error)
}

type RuneScanner interface {
    RuneReader
    UnreadRune() error
}
９、嵌套接口　ReadCloser、ReadSeeker、ReadWriteCloser、ReadWriteSeeker、ReadWriter、WriteCloser 和 WriteSeeker 接口

１０、SectionReader类型

func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader
func (s *SectionReader) Size() int64
func (s *SectionReader) Read(p []byte) (n int, err error)
func (s *SectionReader) ReadAt(p []byte, off int64) (n int, err error)
func (s *SectionReader) Seek(offset int64, whence int) (int64, error)
１１、 LimitedReader类型

func LimitReader(r Reader, n int64) Reader
func (l *LimitedReader) Read(p []byte) (n int, err error)
１２、PipeReader 和 PipeWriter 类型

func Pipe() (*PipeReader, *PipeWriter)
func (r *PipeReader) Read(data []byte) (n int, err error)
func (r *PipeReader) Close() error
func (r *PipeReader) CloseWithError(err error) error

func (w *PipeWriter) Write(data []byte) (n int, err error)
func (w *PipeWriter) Close() error
func (w *PipeWriter) CloseWithError(err error) error
１３、Copy 和 CopyN 函数

func Copy(dst Writer, src Reader) (written int64, err error)
func CopyN(dst Writer, src Reader, n int64) (written int64, err error)
１４、ReadAtLeast 和 ReadFull 函数

func ReadFull(r Reader, buf []byte) (n int, err error)
func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error)
１５、WriteString 函数

func WriteString(w Writer, s string) (n int, err error)
１６、MultiReader 和 MultiWriter 函数

func MultiWriter(writers ...Writer) Writer
func MultiReader(readers ...Reader) Reader
１７、TeeReader函数


func TeeReader(r Reader, w Writer) Reader
