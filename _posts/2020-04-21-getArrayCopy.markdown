---
title: PHP SPL（PHP 标准库）
layout: post
category: php
author: 夏泽民
---
一、什么是spl库？
SPL是用于解决典型问题(standard problems)的一组接口与类的集合。

此扩展只能在php 5.0以后使用，从PHP 5.3.0 不再被关闭,会一直有效.成为php内核组件一部份。

SPL提供了一组标准数据结构。

二、SPL如何使用？
1.构建此扩展不需要其他扩展。

更详细的情况可参考 http://php.net/manual/zh/spl.datastructures.php
<!-- more -->
双向链表
双链表是一种重要的线性存储结构，对于双链表中的每个节点，不仅仅存储自己的信息，还要保存前驱和后继节点的地址。

SplDoublyLinkedList

SplStack（栈）

SplQueue（队列）

SplDoublyLinkedList implements Iterator , ArrayAccess , Countable {    /* 方法 */
    public __construct ( void )    
    public void add ( mixed $index , mixed $newval )    
    public mixed bottom ( void )//双链表的尾部节点
    public int count ( void )//双联表元素的个数
    public mixed current ( void )//当前记录
    public int getIteratorMode ( void ) //获取迭代模式
    public bool isEmpty ( void )//检测双链表是否为空
    public mixed key ( void )//当前节点索引
    public void next ( void )//移到下条记录
    public bool offsetExists ( mixed $index )//指定index处节点是否存在
    public mixed offsetGet ( mixed $index )//获取指定index处节点值
    public void offsetSet ( mixed $index , mixed $newval )//设置指定index处值
    public void offsetUnset ( mixed $index )//删除指定index处节点
    public mixed pop ( void )//从双链表的尾部弹出元素
    public void prev ( void )//移到上条记录
    public void push ( mixed $value )//添加元素到双链表的尾部
    public void rewind ( void )//将指针指向迭代开始处
    public string serialize ( void )//序列化存储
    public void setIteratorMode ( int $mode )//设置迭代模式
    public mixed shift ( void )//双链表的头部移除元素
    public mixed top ( void )//双链表的头部节点
    public void unserialize ( string $serialized )//反序列化
    public void unshift ( mixed $value )//双链表的头部添加元素
    public bool valid ( void )//检查双链表是否还有节点
}


 接下来是使用方法：

$list = new SplDoublyLinkedList();
$list->push('a');
$list->push('b');
$list->push('c');
$list->push('d');
 
$list->unshift('top');
$list->shift();
 
$list->rewind();//rewind操作用于把节点指针指向Bottom所在的节点
echo 'curren node:'.$list->current()."<br />";//获取当前节点
 
$list->next();//指针指向下一个节点
echo 'next node:'.$list->current()."<br />";
 
$list->next();
$list->next();
$list->prev();//指针指向上一个节点
echo 'next node:'.$list->current()."<br />";
 
if($list->current())
    echo 'current node is valid<br />';
else
    echo 'current node is invalid<br />';
     
if($list->valid())//如果当前节点是有效节点，valid返回true
    echo "valid list<br />";
else
  echo "invalid list <br />";
 
var_dump(array(
    'pop' => $list->pop(),
    'count' => $list->count(),    
    'isEmpty' => $list->isEmpty(),    
    'bottom' => $list->bottom(),    
    'top' => $list->top()
));
 
$list->setIteratorMode(SplDoublyLinkedList::IT_MODE_FIFO);
var_dump($list->getIteratorMode());
 
for($list->rewind(); $list->valid(); $list->next()){
    echo $list->current().PHP_EOL;
}
 
var_dump($a = $list->serialize());
//print_r($list->unserialize($a));
 
$list->offsetSet(0,'new one');
$list->offsetUnset(0);
var_dump(array(
    'offsetExists' => $list->offsetExists(4),    
    'offsetGet' => $list->offsetGet(0),
));
var_dump($list);
 
//堆栈，先进后出
$stack = new SplStack();//继承自SplDoublyLinkedList类
 
$stack->push("a<br />");
$stack->push("b<br />");
 
echo $stack->pop();
echo $stack->pop();
echo $stack->offsetSet(0,'B');//堆栈的offset=0是Top所在的位置，offset=1是Top位置节点靠近bottom位置的相邻节点，以此类推
 
$stack->rewind();//双向链表的rewind和堆栈的rewind相反，堆栈的rewind使得当前指针指向Top所在的位置，而双向链表调用之后指向bottom所在位置
echo 'current:'.$stack->current().'<br />';
$stack->next();//堆栈的next操作使指针指向靠近bottom位置的下一个节点，而双向链表是靠近top的下一个节点
echo 'current:'.$stack->current().'<br />';
echo '<br /><br />';
 
//队列，先进先出
$queue = new SplQueue();//继承自SplDoublyLinkedList类
$queue->enqueue("a<br />");//插入一个节点到队列里面的Top位置
$queue->enqueue("b<br />");
$queue->offsetSet(0,'A');//堆栈的offset=0是Top所在的位置，offset=1是Top位置节点靠近bottom位置的相邻节点，以此类推
echo $queue->dequeue();
echo $queue->dequeue();
echo "<br /><br />";


堆
堆(Heap)就是为了实现优先队列而设计的一种数据结构，它是通过构造二叉堆(二叉树的一种)实现。根节点最大的堆叫做最大堆或大根堆（SplMaxHeap），根节点最小的堆叫做最小堆或小根堆（SplMinHeap）。二叉堆还常用于排序(堆排序)

SplHeap

SplMaxHeap

SplMinHeap

SplPriorityQueue

abstract SplHeap implements Iterator , Countable {    
    /* 方法 用法同双向链表一致 */
    public __construct ( void )    
    abstract protected int compare ( mixed $value1 , mixed $value2 )    
    public int count ( void )    
    public mixed current ( void )    
    public mixed extract ( void )    
    public void insert ( mixed $value )    
    public bool isEmpty ( void )    
    public mixed key ( void )    
    public void next ( void )    
    public void recoverFromCorruption ( void )    
    public void rewind ( void )    
    public mixed top ( void )    
    public bool valid ( void )
}


使用方法：

//堆
class MySplHeap extends SplHeap{
    //compare()方法用来比较两个元素的大小，绝对他们在堆中的位置
    public function compare( $value1, $value2 ) {
        return ( $value1 - $value2 );
        }
}
 
$obj = new MySplHeap();
$obj->insert(0);
$obj->insert(1);
$obj->insert(2);
$obj->insert(3);
$obj->insert(4);
echo $obj->top();//4
echo $obj->count();//5
 
foreach ($obj as $item) {
    echo $item."<br />";
}


阵列
优先队列也是非常实用的一种数据结构，可以通过加权对值进行排序，由于排序在php内部实现，业务代码中将精简不少而且更高效。通过SplPriorityQueue::setExtractFlags(int  $flag)设置提取方式可以提取数据（等同最大堆）、优先级、和两者都提取的方式。

SplFixedArray

SplFixedArray implements Iterator , ArrayAccess , Countable {
　　/* 方法 */　　
　　public __construct ([ int $size = 0 ] )
　　public int count ( void )
　　public mixed current ( void )
　　public static SplFixedArray fromArray ( array $array [, bool $save_indexes = true ] )
　　public int getSize ( void )
　　public int key ( void )
　　public void next ( void )
　　public bool offsetExists ( int $index )
　　public mixed offsetGet ( int $index )
　　public void offsetSet ( int $index , mixed $newval )
　　public void offsetUnset ( int $index )
　　public void rewind ( void )
　　public int setSize ( int $size )
　　public array toArray ( void )
　　public bool valid ( void )
　　public void __wakeup ( void )
}


使用方法：

$arr = new SplFixedArray(4);
$arr[0] = 'php';
$arr[1] = 1;
$arr[3] = 'python';//遍历， $arr[2] 为null
foreach($arr as $v) {
    echo $v . PHP_EOL;
}
 
//获取数组长度
echo $arr->getSize(); //4
 
//增加数组长度
$arr->setSize(5);
$arr[4] = 'new one';
 
//捕获异常
try{
    echo $arr[10];
} catch (RuntimeException $e) {
    echo $e->getMessage();
}


映射
用来存储一组对象的，特别是当你需要唯一标识对象的时候。

SplObjectStorage

SplObjectStorage implements Countable , Iterator , Serializable , ArrayAccess {
　　/* 方法 */　　
　　public void addAll ( SplObjectStorage $storage )
　　public void attach ( object $object [, mixed $data = NULL ] )
　　public bool contains ( object $object )
　　public int count ( void )
　　public object current ( void )
　　public void detach ( object $object )
　　public string getHash ( object $object )
　　public mixed getInfo ( void )
　　public int key ( void )
　　public void next ( void )
　　public bool offsetExists ( object $object )
　　public mixed offsetGet ( object $object )
　　public void offsetSet ( object $object [, mixed $data = NULL ] )
　　public void offsetUnset ( object $object )
　　public void removeAll ( SplObjectStorage $storage )
　　public void removeAllExcept ( SplObjectStorage $storage )
　　public void rewind ( void )
　　public string serialize ( void )
　　public void setInfo ( mixed $data )
　　public void unserialize ( string $serialized )
　　public bool valid ( void )
}


使用方法：

class A {
    public $i;    
    public function __construct($i) {
        $this->i = $i;
    }
} 
 
$a1 = new A(1);
$a2 = new A(2);
$a3 = new A(3);
$a4 = new A(4); 
 
$container = new SplObjectStorage(); 
 
//SplObjectStorage::attach 添加对象到Storage中
$container->attach($a1);
$container->attach($a2);
$container->attach($a3); 
 
//SplObjectStorage::detach 将对象从Storage中移除
$container->detach($a2); 
 
//SplObjectStorage::contains用于检查对象是否存在Storage中
var_dump($container->contains($a1)); //true
var_dump($container->contains($a4)); //false
 
//遍历
$container->rewind();
while($container->valid()) {
    var_dump($container->current());    
    $container->next();
}


https://www.w3cschool.cn/doc_php/php-arrayobject-getarraycopy.html
Exports the ArrayObject to an array.

class OrderInfo extends \ArrayObject implements \JsonSerializable {}

https://www.w3cschool.cn/doc_php/php-jsonserializable-jsonserialize.html?lang=en

https://www.php.net/manual/en/arrayiterator.getarraycopy.php

https://www.php.net/manual/de/jsonserializable.jsonserialize.php

https://www.php.net/manual/zh/class.arrayobject.php

https://www.php.net/manual/en/arrayobject.offsetunset.php

https://www.php.net/manual/en/arrayobject.getarraycopy.php


第一部分 简介

1. 什么是SPL？

2. 什么是Iterator？

第二部分 SPL Interfaces

3. Iterator界面

4. ArrayAccess界面

5. IteratorAggregate界面

6. RecursiveIterator界面

7. SeekableIterator界面

8. Countable界面

第三部分 SPL Classes

9. SPL的内置类

10. DirectoryIterator类

11. ArrayObject类

12. ArrayIterator类

13. RecursiveArrayIterator类和RecursiveIteratorIterator类

14. FilterIterator类

15. SimpleXMLIterator类

16. CachingIterator类

17. LimitIterator类

18. SplFileObject类

第一部 简介

1. 什么是SPL？
SPL是Standard PHP Library（PHP标准库）的缩写。

根据官方定义，它是"a collection of interfaces and classes that are meant to solve standard problems"。但是，目前在使用中，SPL更多地被看作是一种使object（物体）模仿array（数组）行为的interfaces和classes。

2. 什么是Iterator？

SPL的核心概念就是Iterator。这指的是一种Design Pattern，根据《Design Patterns》一书的定义，Iterator的作用是"provide an object which traverses some aggregate structure, abstracting away assumptions about the implementation of that structure."

wikipedia中说，"an iterator is an object which allows a programmer to traverse through all the elements of a collection, regardless of its specific implementation"......."the iterator pattern is a design pattern in which iterators are used to access the elements of an aggregate object sequentially without exposing its underlying representation".

通俗地说，Iterator能够使许多不同的数据结构，都能有统一的操作界面，比如一个数据库的结果集、同一个目录中的文件集、或者一个文本中每一行构成的集合。

如果按照普通情况，遍历一个MySQL的结果集，程序需要这样写：


// Fetch the "aggregate structure"
$result = mysql_query("SELECT * FROM users");

// Iterate over the structure
while ( $row = mysql_fetch_array($result) ) {
   // do stuff with the row here
}

读出一个目录中的内容，需要这样写：


// Fetch the "aggregate structure"
$dh = opendir('/home/harryf/files');

// Iterate over the structure
while ( $file = readdir($dh) ) {
   // do stuff with the file here
}

读出一个文本文件的内容，需要这样写：


// Fetch the "aggregate structure"
$fh = fopen("/home/hfuecks/files/results.txt", "r");

// Iterate over the structure
while (!feof($fh)) {

   $line = fgets($fh);
   // do stuff with the line here

}

上面三段代码，虽然处理的是不同的resource（资源），但是功能都是遍历结果集（loop over contents），因此Iterator的基本思想，就是将这三种不同的操作统一起来，用同样的命令界面，处理不同的资源。

第二部分 SPL Interfaces

3. Iterator界面

SPL规定，所有部署了Iterator界面的class，都可以用在foreach Loop中。Iterator界面中包含5个必须部署的方法：


    * current()

      This method returns the current index's value. You are solely
      responsible for tracking what the current index is as the 
     interface does not do this for you.

    * key()

      This method returns the value of the current index's key. For 
      foreach loops this is extremely important so that the key 
      value can be populated.

    * next()

      This method moves the internal index forward one entry.

    * rewind()

      This method should reset the internal index to the first element.

    * valid()

      This method should return true or false if there is a current 
      element. It is called after rewind() or next().

下面就是一个部署了Iterator界面的class示例：


/**
* An iterator for native PHP arrays, re-inventing the wheel
*
* Notice the "implements Iterator" - important!
*/
class ArrayReloaded implements Iterator {

   /**
   * A native PHP array to iterate over
   */
 private $array = array();

   /**
   * A switch to keep track of the end of the array
   */
 private $valid = FALSE;

   /**
   * Constructor
   * @param array native PHP array to iterate over
   */
 function __construct($array) {
   $this->array = $array;
 }

   /**
   * Return the array "pointer" to the first element
   * PHP's reset() returns false if the array has no elements
   */
 function rewind(){
   $this->valid = (FALSE !== reset($this->array));
 }

   /**
   * Return the current array element
   */
 function current(){
   return current($this->array);
 }

   /**
   * Return the key of the current array element
   */
 function key(){
   return key($this->array);
 }

   /**
   * Move forward by one
   * PHP's next() returns false if there are no more elements
   */
 function next(){
   $this->valid = (FALSE !== next($this->array));
 }

   /**
   * Is the current element valid?
   */
 function valid(){
   return $this->valid;
 }
}

使用方法如下：


// Create iterator object
$colors = new ArrayReloaded(array ('red','green','blue',));

// Iterate away!
foreach ( $colors as $color ) {
 echo $color."<br>";
}

你也可以在foreach循环中使用key()方法：


// Display the keys as well
foreach ( $colors as $key => $color ) {
 echo "$key: $color<br>";
}

除了foreach循环外，也可以使用while循环，


// Reset the iterator - foreach does this automatically
$colors->rewind();

// Loop while valid
while ( $colors->valid() ) {

   echo $colors->key().": ".$colors->current()."
";
   $colors->next();

}

根据测试，while循环要稍快于foreach循环，因为运行时少了一层中间调用。

4. ArrayAccess界面
部署ArrayAccess界面，可以使得object像array那样操作。ArrayAccess界面包含四个必须部署的方法：


    * offsetExists($offset)

      This method is used to tell php if there is a value
      for the key specified by offset. It should return 
      true or false.

    * offsetGet($offset)

      This method is used to return the value specified 
      by the key offset.

    * offsetSet($offset, $value)

      This method is used to set a value within the object, 
      you can throw an exception from this function for a 
      read-only collection.

    * offsetUnset($offset)

      This method is used when a value is removed from 
      an array either through unset() or assigning the key 
      a value of null. In the case of numerical arrays, this 
      offset should not be deleted and the array should 
      not be reindexed unless that is specifically the 
      behavior you want.

下面就是一个部署ArrayAccess界面的实例：


/**
* A class that can be used like an array
*/
class Article implements ArrayAccess {

 public $title;

 public $author;

 public $category;  

 function __construct($title,$author,$category) {
   $this->title = $title;
   $this->author = $author;
   $this->category = $category;
 }

 /**
 * Defined by ArrayAccess interface
 * Set a value given it's key e.g. $A['title'] = 'foo';
 * @param mixed key (string or integer)
 * @param mixed value
 * @return void
 */
 function offsetSet($key, $value) {
   if ( array_key_exists($key,get_object_vars($this)) ) {
     $this->{$key} = $value;
   }
 }

 /**
 * Defined by ArrayAccess interface
 * Return a value given it's key e.g. echo $A['title'];
 * @param mixed key (string or integer)
 * @return mixed value
 */
 function offsetGet($key) {
   if ( array_key_exists($key,get_object_vars($this)) ) {
     return $this->{$key};
   }
 }

 /**
 * Defined by ArrayAccess interface
 * Unset a value by it's key e.g. unset($A['title']);
 * @param mixed key (string or integer)
 * @return void
 */
 function offsetUnset($key) {
   if ( array_key_exists($key,get_object_vars($this)) ) {
     unset($this->{$key});
   }
 }

 /**
 * Defined by ArrayAccess interface
 * Check value exists, given it's key e.g. isset($A['title'])
 * @param mixed key (string or integer)
 * @return boolean
 */
 function offsetExists($offset) {
   return array_key_exists($offset,get_object_vars($this));
 }

}

使用方法如下：


// Create the object
$A = new Article('SPL Rocks','Joe Bloggs', 'PHP');

// Check what it looks like
echo 'Initial State:<div>';
print_r($A);
echo '</div>';

// Change the title using array syntax
$A['title'] = 'SPL _really_ rocks';

// Try setting a non existent property (ignored)
$A['not found'] = 1;

// Unset the author field
unset($A['author']);

// Check what it looks like again
echo 'Final State:<div>';
print_r($A);
echo '</div>';

运行结果如下：


Initial State:

Article Object
(
   [title] => SPL Rocks
   [author] => Joe Bloggs
   [category] => PHP
)

Final State:

Article Object
(
   [title] => SPL _really_ rocks
   [category] => PHP
)

可以看到，$A虽然是一个object，但是完全可以像array那样操作。

你还可以在读取数据时，增加程序内部的逻辑：


function offsetGet($key) {
   if ( array_key_exists($key,get_object_vars($this)) ) {
     return strtolower($this->{$key});
   }
 }

5. IteratorAggregate界面

但是，虽然$A可以像数组那样操作，却无法使用foreach遍历，除非部署了前面提到的Iterator界面。

另一个解决方法是，有时会需要将数据和遍历部分分开，这时就可以部署IteratorAggregate界面。它规定了一个getIterator()方法，返回一个使用Iterator界面的object。

还是以上一节的Article类为例：


class Article implements ArrayAccess, IteratorAggregate {

/**
 * Defined by IteratorAggregate interface
 * Returns an iterator for for this object, for use with foreach
 * @return ArrayIterator
 */
 function getIterator() {
   return new ArrayIterator($this);
 }

使用方法如下：


$A = new Article('SPL Rocks','Joe Bloggs', 'PHP');

// Loop (getIterator will be called automatically)
echo 'Looping with foreach:<div>';
foreach ( $A as $field => $value ) {
 echo "$field : $value<br>";
}
echo '</div>';

// Get the size of the iterator (see how many properties are left)
echo "Object has ".sizeof($A->getIterator())." elements";

显示结果如下：


Looping with foreach:

title : SPL Rocks
author : Joe Bloggs
category : PHP

Object has 3 elements

6. RecursiveIterator界面

这个界面用于遍历多层数据，它继承了Iterator界面，因而也具有标准的current()、key()、next()、 rewind()和valid()方法。同时，它自己还规定了getChildren()和hasChildren()方法。The getChildren() method must return an object that implements RecursiveIterator.

7. SeekableIterator界面

SeekableIterator界面也是Iterator界面的延伸，除了Iterator的5个方法以外，还规定了seek()方法，参数是元素的位置，返回该元素。如果该位置不存在，则抛出OutOfBoundsException。

下面是一个是实例：


<?php

class PartyMemberIterator implements SeekableIterator
{
    public function __construct(PartyMember $member)
    {
        // Store $member locally for iteration
    }

    public function seek($index)
    {
        $this->rewind();
        $position = 0;

        while ($position < $index && $this->valid()) {
            $this->next();
            $position++;
        }

        if (!$this->valid()) {
            throw new OutOfBoundsException('Invalid position');
        }
    }

    // Implement current(), key(), next(), rewind()
    // and valid() to iterate over data in $member
}

?>

8. Countable界面

这个界面规定了一个count()方法，返回结果集的数量。

第三部分 SPL Classes

9. SPL的内置类

SPL除了定义一系列Interfaces以外，还提供一系列的内置类，它们对应不同的任务，大大简化了编程。

查看所有的内置类，可以使用下面的代码：


<?php
// a simple foreach() to traverse the SPL class names
foreach(spl_classes() as $key=>$value)
        {
        echo $key.' -&gt; '.$value.'<br />';
        }
?>

10. DirectoryIterator类

这个类用来查看一个目录中的所有文件和子目录：


<?php

try{
  /*** class create new DirectoryIterator Object ***/
    foreach ( new DirectoryIterator('./') as $Item )
        {
        echo $Item.'<br />';
        }
    }
/*** if an exception is thrown, catch it here ***/
catch(Exception $e){
    echo 'No files Found!<br />';
}
?>

查看文件的详细信息：


<table>
<?php

foreach(new DirectoryIterator('./' ) as $file )
    {
    if( $file->getFilename()  == 'foo.txt' )
        {
        echo '<tr><td>getFilename()</td><td> '; var_dump($file->getFilename()); echo '</td></tr>';
    echo '<tr><td>getBasename()</td><td> '; var_dump($file->getBasename()); echo '</td></tr>';
        echo '<tr><td>isDot()</td><td> '; var_dump($file->isDot()); echo '</td></tr>';
        echo '<tr><td>__toString()</td><td> '; var_dump($file->__toString()); echo '</td></tr>';
        echo '<tr><td>getPath()</td><td> '; var_dump($file->getPath()); echo '</td></tr>';
        echo '<tr><td>getPathname()</td><td> '; var_dump($file->getPathname()); echo '</td></tr>';
        echo '<tr><td>getPerms()</td><td> '; var_dump($file->getPerms()); echo '</td></tr>';
        echo '<tr><td>getInode()</td><td> '; var_dump($file->getInode()); echo '</td></tr>';
        echo '<tr><td>getSize()</td><td> '; var_dump($file->getSize()); echo '</td></tr>';
        echo '<tr><td>getOwner()</td><td> '; var_dump($file->getOwner()); echo '</td></tr>';
        echo '<tr><td>$file->getGroup()</td><td> '; var_dump($file->getGroup()); echo '</td></tr>';
        echo '<tr><td>getATime()</td><td> '; var_dump($file->getATime()); echo '</td></tr>';
        echo '<tr><td>getMTime()</td><td> '; var_dump($file->getMTime()); echo '</td></tr>';
        echo '<tr><td>getCTime()</td><td> '; var_dump($file->getCTime()); echo '</td></tr>';
        echo '<tr><td>getType()</td><td> '; var_dump($file->getType()); echo '</td></tr>';
        echo '<tr><td>isWritable()</td><td> '; var_dump($file->isWritable()); echo '</td></tr>';
        echo '<tr><td>isReadable()</td><td> '; var_dump($file->isReadable()); echo '</td></tr>';
        echo '<tr><td>isExecutable(</td><td> '; var_dump($file->isExecutable()); echo '</td></tr>';
        echo '<tr><td>isFile()</td><td> '; var_dump($file->isFile()); echo '</td></tr>';
        echo '<tr><td>isDir()</td><td> '; var_dump($file->isDir()); echo '</td></tr>';
        echo '<tr><td>isLink()</td><td> '; var_dump($file->isLink()); echo '</td></tr>';
        echo '<tr><td>getFileInfo()</td><td> '; var_dump($file->getFileInfo()); echo '</td></tr>';
        echo '<tr><td>getPathInfo()</td><td> '; var_dump($file->getPathInfo()); echo '</td></tr>';
        echo '<tr><td>openFile()</td><td> '; var_dump($file->openFile()); echo '</td></tr>';
        echo '<tr><td>setFileClass()</td><td> '; var_dump($file->setFileClass()); echo '</td></tr>';
        echo '<tr><td>setInfoClass()</td><td> '; var_dump($file->setInfoClass()); echo '</td></tr>';
        }
}
?>
</table>

除了foreach循环外，还可以使用while循环：


<?php
/*** create a new iterator object ***/
$it = new DirectoryIterator('./');

/*** loop directly over the object ***/
while($it->valid())
    {
    echo $it->key().' -- '.$it->current().'<br />';
    /*** move to the next iteration ***/
    $it->next();
    }
?>

如果要过滤所有子目录，可以在valid()方法中过滤：


<?php
/*** create a new iterator object ***/
$it = new DirectoryIterator('./');

/*** loop directly over the object ***/
while($it->valid())
        {
        /*** check if value is a directory ***/
        if($it->isDir())
                {
                /*** echo the key and current value ***/
                echo $it->key().' -- '.$it->current().'<br />';
                }
        /*** move to the next iteration ***/
        $it->next();
        }
?>

11. ArrayObject类

这个类可以将Array转化为object。


<?php

/*** a simple array ***/
$array = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'kiwi', 'kookaburra', 'platypus');

/*** create the array object ***/
$arrayObj = new ArrayObject($array);

/*** iterate over the array ***/
for($iterator = $arrayObj->getIterator();
   /*** check if valid ***/
   $iterator->valid();
   /*** move to the next array member ***/
   $iterator->next())
    {
    /*** output the key and current array value ***/
    echo $iterator->key() . ' => ' . $iterator->current() . '<br />';
    }
?>

增加一个元素：


$arrayObj->append('dingo');

对元素排序：


$arrayObj->natcasesort();

显示元素的数量：


echo $arrayObj->count();

删除一个元素：


$arrayObj->offsetUnset(5);

某一个元素是否存在：


 if ($arrayObj->offsetExists(3))
    {
       echo 'Offset Exists<br />';
    }

更改某个位置的元素值：


 $arrayObj->offsetSet(5, "galah");

显示某个位置的元素值：


echo $arrayObj->offsetGet(4);

12. ArrayIterator类
这个类实际上是对ArrayObject类的补充，为后者提供遍历功能。

示例如下：


<?php
/*** a simple array ***/
$array = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'kiwi', 'kookaburra', 'platypus');

try {
    $object = new ArrayIterator($array);
    foreach($object as $key=>$value)
        {
        echo $key.' => '.$value.'<br />';
        }
    }
catch (Exception $e)
    {
    echo $e->getMessage();
    }
?>

ArrayIterator类也支持offset类方法和count()方法：


<ul>
<?php
/*** a simple array ***/
$array = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'kiwi', 'kookaburra', 'platypus');

try {
    $object = new ArrayIterator($array);
    /*** check for the existence of the offset 2 ***/
    if($object->offSetExists(2))
    {
    /*** set the offset of 2 to a new value ***/
    $object->offSetSet(2, 'Goanna');
    }
   /*** unset the kiwi ***/
   foreach($object as $key=>$value)
        {
        /*** check the value of the key ***/
        if($object->offSetGet($key) === 'kiwi')
            {
            /*** unset the current key ***/
            $object->offSetUnset($key);
            }
        echo '<li>'.$key.' - '.$value.'</li>'."\n";
        }
    }
catch (Exception $e)
    {
    echo $e->getMessage();
    }
?>
</ul>

13. RecursiveArrayIterator类和RecursiveIteratorIterator类

ArrayIterator类和ArrayObject类，只支持遍历一维数组。如果要遍历多维数组，必须先用RecursiveIteratorIterator生成一个Iterator，然后再对这个Iterator使用RecursiveIteratorIterator。


<?php
$array = array(
    array('name'=>'butch', 'sex'=>'m', 'breed'=>'boxer'),
    array('name'=>'fido', 'sex'=>'m', 'breed'=>'doberman'),
    array('name'=>'girly','sex'=>'f', 'breed'=>'poodle')
);

foreach(new RecursiveIteratorIterator(new RecursiveArrayIterator($array)) as $key=>$value)
    {
    echo $key.' -- '.$value.'<br />';
    }
?>

14. FilterIterator类

FilterIterator类可以对元素进行过滤，只要在accept()方法中设置过滤条件就可以了。

示例如下：


<?php
/*** a simple array ***/
$animals = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'NZ'=>'kiwi', 'kookaburra', 'platypus');

class CullingIterator extends FilterIterator{

/*** The filteriterator takes  a iterator as param: ***/
public function __construct( Iterator $it ){
  parent::__construct( $it );
}

/*** check if key is numeric ***/
function accept(){
  return is_numeric($this->key());
}

}/*** end of class ***/
$cull = new CullingIterator(new ArrayIterator($animals));

foreach($cull as $key=>$value)
    {
    echo $key.' == '.$value.'<br />';
    }
?>

下面是另一个返回质数的例子：


<?php

class PrimeFilter extends FilterIterator{

/*** The filteriterator takes  a iterator as param: ***/
public function __construct(Iterator $it){
  parent::__construct($it);
}

/*** check if current value is prime ***/
function accept(){
if($this->current() % 2 != 1)
    {
    return false;
    }
$d = 3;
$x = sqrt($this->current());
while ($this->current() % $d != 0 && $d < $x)
    {
    $d += 2;
    }
 return (($this->current() % $d == 0 && $this->current() != $d) * 1) == 0 ? true : false;
}

}/*** end of class ***/

/*** an array of numbers ***/
$numbers = range(212345,212456);

/*** create a new FilterIterator object ***/
$primes = new primeFilter(new ArrayIterator($numbers));

foreach($primes as $value)
    {
    echo $value.' is prime.<br />';
    }
?>

15. SimpleXMLIterator类

这个类用来遍历xml文件。

示例如下：


<?php

/*** a simple xml tree ***/
 $xmlstring = <<<XML
<?xml version = "1.0" encoding="UTF-8" standalone="yes"?>
<document>
  <animal>
    <category id="26">
      <species>Phascolarctidae</species>
      <type>koala</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="27">
      <species>macropod</species>
      <type>kangaroo</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="28">
      <species>diprotodon</species>
      <type>wombat</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="31">
      <species>macropod</species>
      <type>wallaby</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="21">
      <species>dromaius</species>
      <type>emu</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="22">
      <species>Apteryx</species>
      <type>kiwi</type>
      <name>Troy</name>
    </category>
  </animal>
  <animal>
    <category id="23">
      <species>kingfisher</species>
      <type>kookaburra</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="48">
      <species>monotremes</species>
      <type>platypus</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="4">
      <species>arachnid</species>
      <type>funnel web</type>
      <name>Bruce</name>
      <legs>8</legs>
    </category>
  </animal>
</document>
XML;

/*** a new simpleXML iterator object ***/
try    {
       /*** a new simple xml iterator ***/
       $it = new SimpleXMLIterator($xmlstring);
       /*** a new limitIterator object ***/
       foreach(new RecursiveIteratorIterator($it,1) as $name => $data)
          {
          echo $name.' -- '.$data.'<br />';
          }
    }
catch(Exception $e)
    {
    echo $e->getMessage();
    }
?>

new RecursiveIteratorIterator($it,1)表示显示所有包括父元素在内的子元素。

显示某一个特定的元素值，可以这样写：


<?php
try {
    /*** a new simpleXML iterator object ***/
    $sxi =  new SimpleXMLIterator($xmlstring);

    foreach ( $sxi as $node )
        {
        foreach($node as $k=>$v)
            {
            echo $v->species.'<br />';
            }
        }
    }
catch(Exception $e)
    {
    echo $e->getMessage();
    }
?>

相对应的while循环写法为：


<?php

try {
$sxe = simplexml_load_string($xmlstring, 'SimpleXMLIterator');

for ($sxe->rewind(); $sxe->valid(); $sxe->next())
    {
    if($sxe->hasChildren())
        {
        foreach($sxe->getChildren() as $element=>$value)
          {
          echo $value->species.'<br />';
          }
        }
     }
   }
catch(Exception $e)
   {
   echo $e->getMessage();
   }
?>

最方便的写法，还是使用xpath：


<?php
try {
    /*** a new simpleXML iterator object ***/
    $sxi =  new SimpleXMLIterator($xmlstring);

    /*** set the xpath ***/
    $foo = $sxi->xpath('animal/category/species');

    /*** iterate over the xpath ***/
    foreach ($foo as $k=>$v)
        {
        echo $v.'<br />';
        }
    }
catch(Exception $e)
    {
    echo $e->getMessage();
    }
?>

下面的例子，显示有namespace的情况：


<?php

/*** a simple xml tree ***/
 $xmlstring = <<<XML
<?xml version = "1.0" encoding="UTF-8" standalone="yes"?>
<document xmlns:spec="http://example.org/animal-species">
  <animal>
    <category id="26">
      <species>Phascolarctidae</species>
      <spec:name>Speed Hump</spec:name>
      <type>koala</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="27">
      <species>macropod</species>
      <spec:name>Boonga</spec:name>
      <type>kangaroo</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="28">
      <species>diprotodon</species>
      <spec:name>pot holer</spec:name>
      <type>wombat</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="31">
      <species>macropod</species>
      <spec:name>Target</spec:name>
      <type>wallaby</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="21">
      <species>dromaius</species>
      <spec:name>Road Runner</spec:name>
      <type>emu</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="22">
      <species>Apteryx</species>
      <spec:name>Football</spec:name>
      <type>kiwi</type>
      <name>Troy</name>
    </category>
  </animal>
  <animal>
    <category id="23">
      <species>kingfisher</species>
      <spec:name>snaker</spec:name>
      <type>kookaburra</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="48">
      <species>monotremes</species>
      <spec:name>Swamp Rat</spec:name>
      <type>platypus</type>
      <name>Bruce</name>
    </category>
  </animal>
  <animal>
    <category id="4">
      <species>arachnid</species>
      <spec:name>Killer</spec:name>
      <type>funnel web</type>
      <name>Bruce</name>
      <legs>8</legs>
    </category>
  </animal>
</document>
XML;

/*** a new simpleXML iterator object ***/
try {
    /*** a new simpleXML iterator object ***/
    $sxi =  new SimpleXMLIterator($xmlstring);

    $sxi-> registerXPathNamespace('spec', 'http://www.exampe.org/species-title');

    /*** set the xpath ***/
    $result = $sxi->xpath('//spec:name');

    /*** get all declared namespaces ***/
   foreach($sxi->getDocNamespaces('animal') as $ns)
        {
        echo $ns.'<br />';
        }

    /*** iterate over the xpath ***/
    foreach ($result as $k=>$v)
        {
        echo $v.'<br />';
        }
    }
catch(Exception $e)
    {
    echo $e->getMessage();
    }
?>

增加一个节点：


<?php 
 $xmlstring = <<<XML
<?xml version = "1.0" encoding="UTF-8" standalone="yes"?>
<document>
  <animal>koala</animal>
  <animal>kangaroo</animal>
  <animal>wombat</animal>
  <animal>wallaby</animal>
  <animal>emu</animal>
  <animal>kiwi</animal>
  <animal>kookaburra</animal>
  <animal>platypus</animal>
  <animal>funnel web</animal>
</document>
XML;

try {
    /*** a new simpleXML iterator object ***/
    $sxi =  new SimpleXMLIterator($xmlstring);

    /*** add a child ***/
    $sxi->addChild('animal', 'Tiger');

    /*** a new simpleXML iterator object ***/
    $new = new SimpleXmlIterator($sxi->saveXML());

    /*** iterate over the new tree ***/
    foreach($new as $val)
        {
        echo $val.'<br />';
        }
    }
catch(Exception $e)
    {
    echo $e->getMessage();
    }
?>

增加属性：


<?php 
$xmlstring =<<<XML
<?xml version = "1.0" encoding="UTF-8" standalone="yes"?>
<document>
  <animal>koala</animal>
  <animal>kangaroo</animal>
  <animal>wombat</animal>
  <animal>wallaby</animal>
  <animal>emu</animal>
  <animal>kiwi</animal>
  <animal>kookaburra</animal>
  <animal>platypus</animal>
  <animal>funnel web</animal>
</document>
XML;

try {
    /*** a new simpleXML iterator object ***/
    $sxi =  new SimpleXMLIterator($xmlstring);

    /*** add an attribute with a namespace ***/
    $sxi->addAttribute('id:att1', 'good things', 'urn::test-foo');

    /*** add an attribute without a  namespace ***/
    $sxi->addAttribute('att2', 'no-ns');

    echo htmlentities($sxi->saveXML());
    }
catch(Exception $e)
    {
    echo $e->getMessage();
    }
?>

16. CachingIterator类

这个类有一个hasNext()方法，用来判断是否还有下一个元素。

示例如下：


<?php
/*** a simple array ***/
$array = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'kiwi', 'kookaburra', 'platypus');

try {
    /*** create a new object ***/
    $object = new CachingIterator(new ArrayIterator($array));
    foreach($object as $value)
        {
        echo $value;
        if($object->hasNext())
            {
            echo ',';
            }
        }
    }
catch (Exception $e)
    {
    echo $e->getMessage();
    }
?>

17. LimitIterator类

这个类用来限定返回结果集的数量和位置，必须提供offset和limit两个参数，与SQL命令中limit语句类似。

示例如下：


<?php
/*** the offset value ***/
$offset = 3;

/*** the limit of records to show ***/
$limit = 2;

$array = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'kiwi', 'kookaburra', 'platypus');

$it = new LimitIterator(new ArrayIterator($array), $offset, $limit);

foreach($it as $k=>$v)
    {
    echo $it->getPosition().'<br />';
    }
?>

另一个例子是：


<?php

/*** a simple array ***/
$array = array('koala', 'kangaroo', 'wombat', 'wallaby', 'emu', 'kiwi', 'kookaburra', 'platypus');

$it = new LimitIterator(new ArrayIterator($array));

try
    {
    $it->seek(5);
    echo $it->current();
    }
catch(OutOfBoundsException $e)
    {
    echo $e->getMessage() . "<br />";
    }
?>

18. SplFileObject类

这个类用来对文本文件进行遍历。

示例如下：


<?php

try{
    // iterate directly over the object
    foreach( new SplFileObject(&quot;/usr/local/apache/logs/access_log&quot;) as $line)
    // and echo each line of the file
    echo $line.'<br />';
}
catch (Exception $e)
    {
    echo $e->getMessage();
    }
?>

返回文本文件的第三行，可以这样写：


<?php

try{
    $file = new SplFileObject("/usr/local/apache/logs/access_log");

    $file->seek(3);

    echo $file->current();
        }
catch (Exception $e)
    {
    echo $e->getMessage();
    }
?>

[参考文献]

1. Introduction to Standard PHP Library (SPL), By Kevin Waterson

2. Introducing PHP 5's Standard Library, By Harry Fuecks

3. The Standard PHP Library (SPL), By Ben Ramsey

4. SPL - Standard PHP Library Documentation


http://www.ruanyifeng.com/blog/2008/07/php_spl_notes.html


SPL，PHP 标准库（Standard PHP Library） ，从 PHP 5.0 起内置的组件和接口，并且从 PHP5.3 已逐渐的成熟。SPL 其实在所有的 PHP5 开发环境中被内置，同时无需任何设置。

似乎众多的 PHP 开发人员基本没有使用它，甚至闻所未闻。究其原因，可以追述到它那阳春白雪般的说明文档，使你忽略了「它的存在」。SPL 这块宝石犹如铁达尼的「海洋之心」般，被沉入海底。而现在它应该被我们捞起，并将它穿戴在应有的位置 ，而这也是这篇文章所要表述的观点。

那么，SPL 提供了什么？

SPL 对 PHP 引擎进行了扩展，例如 ArrayAccess、Countable 和 SeekableIterator 等接口，它们用于以数组形式操作对象。同时，你还可以使用 RecursiveIterator、ArrayObejcts 等其他迭代器进行数据的迭代操作。

它还内置几个的对象例如 Exceptions、SplObserver、Spltorage 以及 splautoloadregister、splclasses、iteratorapply 等的帮助函数（helper functions），用于重载对应的功能。

这些工具聚合在一起就好比是把多功能的瑞士军刀，善用它们可以从质上提升 PHP 的代码效率。那么，我们如何发挥它的威力？

重载 autoloader

如果你是位「教科书式的程序员」，那么你保证了解如何使用 __autoload 去代替 includes/requires 操作惰性载入对应的类，对不？

但久之，你会发现你已经陷入了困境，首先是你要保证你的类文件必须在指定的文件路径中，例如在 Zend 框架中你必须使用「_」来分割类、方法名称（你如何解决这一问题？）。

另外的一个问题，就是当项目变得越来越复杂， __autoload内的逻辑也会变得相应的复杂。到最后，甚至你会加入异常判断，以及将所有的载入类的逻辑如数写到其中。

大家都知道「鸡蛋不能放到一个篮子中」，利用 SPL 可以分离 __autoload的载入逻辑。只需要写个你自己的 autoload 函数，然后利用 SPL 提供的函数重载它。

例如上述 Zend 框架的问题，你可以重载 Zend loader 对应的方法，如果它没有找到对应的类，那么就使用你先前定义的函数。

PHP SPL标准库-接口
PHP SPL标准库有一下接口：

Countable
OuterIterator 
RecursiveIterator 
SeekableIterator
SplObserver 
SplSubject
ArrayObject
其中OuterIterator、RecursiveIterator、SeekableIterator都是继承Iterator类的。

Coutable接口：
实现Coutable接口的对象可用于 count() 函数计数。

class Mycount implements Countable
{
    public function count()
    {
        static $count = 0;
        $count++;
        return $count;
    }
}
  
$count = new Mycount();
$count->count();
$count->count();
  
echo count($count); //3
echo count($count); //4
说明：

调用 count() 函数时，Mycount::count() 方法被调用，count() 函数的第二个参数将不会产生影响。

OuterIterator接口：
它是自定义或修改迭代过程。

// IteratorIterator是OuterIterator的一个实现类
class MyOuterIterator extends  IteratorIterator {
    public function current()
    {
        return parent::current() . 'TEST';
    }
}
  
foreach(new MyOuterIterator(new ArrayIterator(['b','a','c'])) as $key => $value) {
    echo "$key->$value".PHP_EOL;
}
 
/*
结果：
0->bTEST
1->aTEST
2->cTEST
*/
在实际应用中，OuterIterator非常有用：

$db = new PDO('mysql:host=localhost;dbname=test', 'root', 'mckee');
$db->query('set names utf8');
$pdoStatement = $db->query('SELECT * FROM test1', PDO::FETCH_ASSOC);
$iterator = new IteratorIterator($pdoStatement);
$tenRecordArray = iterator_to_array($iterator);
print_r($tenRecordArray);
RecursiveIterator接口：
用于循环迭代多层结构的数据，RecursiveIterator另外提供了两个方法：

RecursiveIterator::getChildren：获取当前元素下子迭代器

RecursiveIterator::hasChildren：判断当前元素下是否有迭代器
class MyRecursiveIterator implements RecursiveIterator
{
    private $_data;
    private $_position = 0;
  
    public function __construct(array $data) {
        $this->_data = $data;
    }
  
    public function valid() {
        return isset($this->_data[$this->_position]);
    }
  
    public function hasChildren() {
        return is_array($this->_data[$this->_position]);
    }
  
    public function next() {
        $this->_position++;
    }
  
    public function current() {
        return $this->_data[$this->_position];
    }
  
    public function getChildren() {
        print_r($this->_data[$this->_position]);
    }
  
    public function rewind() {
        $this->_position = 0;
    }
  
    public function key() {
        return $this->_position;
    }
}
  
$arr = array(0, 1=> array(10, 20), 2, 3 => array(1, 2));
$mri = new MyRecursiveIterator($arr);
  
foreach ($mri as $c => $v) {
    if ($mri->hasChildren()) {
        echo "$c has children: " .PHP_EOL;
        $mri->getChildren();
    } else {
        echo "$v" .PHP_EOL;
    }
}
输出结果：
0
1 has children:
Array
(
    [0] => 10
    [1] => 20
)
2
3 has children:
Array
(
    [0] => 1
    [1] => 2
)
SeekableIterator接口：
通过 seek() 方法实现可搜索的迭代器，用于搜索某个位置下的元素。

class  MySeekableIterator  implements  SeekableIterator  {
    private  $position = 0;
 
    private  $array  = array(
        "first element" ,
        "second element" ,
        "third element" ,
        "fourth element"
    );
 
    public function  seek ( $position ) {
        if (!isset( $this -> array [ $position ])) {
            throw new  OutOfBoundsException ( "invalid seek position ( $position )" );
        }
  
       $this -> position  =  $position ;
    }
  
    public function  rewind () {
        $this -> position  =  0 ;
    }
  
    public function  current () {
        return  $this -> array [ $this -> position ];
    }
  
    public function  key () {
        return  $this -> position ;
    }
  
    public function  next () {
        ++ $this -> position ;
    }
  
    public function  valid () {
        return isset( $this -> array [ $this -> position ]);
    }
}
  
try{
    $it  = new  MySeekableIterator ;
    echo  $it -> current (),  "\n" ;
  
    $it -> seek ( 2 );
    echo  $it -> current (),  "\n" ;
  
    $it -> seek ( 1 );
    echo  $it -> current (),  "\n" ;
  
    $it -> seek ( 10 );
}catch( OutOfBoundsException $e ){
    echo  $e -> getMessage ();
}
输出结果：

1
2
3
4
first element
third element
second element
invalid seek position ( 10 )
SplObserver和SplSubject接口：
SplObserver和SplSubject接口用来实现观察者设计模式，观察者设计模式是指当一个类的状态发生变化时，依赖它的对象都会收到通知并更新。使用场景非常广泛，比如说当一个事件发生后，需要更新多个逻辑操作，传统方式是在事件添加后编写逻辑，这种代码耦合并难以维护，观察者模式可实现低耦合的通知和更新机制。

SplObserver和SplSubject的接口结构：
// SplSubject结构 被观察的对象
interface SplSubject{
    public function attach(SplObserver $observer); // 添加观察者
    public function detach(SplObserver $observer); // 剔除观察者
    public function notify(); // 通知观察者
}
  
// SplObserver结构 代表观察者
interface SplObserver{
    public function update(SplSubject $subject); // 更新操作
}
看下面一个实现观察者的例子：
class Subject implements SplSubject
{
    private $observers = array();
  
    public function attach(SplObserver  $observer)
    {
        $this->observers[] = $observer;
    }
  
    public function detach(SplObserver  $observer)
    {
        if($index = array_search($observer, $this->observers, true)) {
            unset($this->observers[$index]);
        }
    }
  
    public function notify()
    {
        foreach($this->observers as $observer) {
            $observer->update($this);
        }
    }
}
 
class Observer1 implements  SplObserver
{
    public function update(SplSubject  $subject)
    {
        echo "逻辑1代码".PHP_EOL;
    }
}
 
class Observer2 implements  SplObserver
{
    public function update(SplSubject  $subject)
    {
        echo "逻辑2代码".PHP_EOL;
    }
}
 
$subject = new Subject();
$subject->attach(new Observer1());
$subject->attach(new Observer2());
 
$subject->notify();
运行结果：

1
2
逻辑1代码
逻辑2代码
ArrayObject接口：
ArrayObject 是将数组转换为数组对象。
// 返回当前数组元素
ArrayIterator::current( void )
 
// 返回当前数组key
ArrayIterator::key(void)
 
// 指向下个数组元素
ArrayIterator::next (void)
 
// 重置数组指针到头
ArrayIterator::rewind(void )
 
// 查找数组中某一位置
ArrayIterator::seek()
 
// 检查数组是否还包含其他元素
ArrayIterator::valid()
 
// 添加新元素
ArrayObject::append()
 
// 构造一个新的数组对象
ArrayObject::__construct()
 
// 返回迭代器中元素个数
ArrayObject::count()
 
// 从一个数组对象构造一个新迭代器
ArrayObject::getIterator()
 
// 判断提交的值是否存在
ArrayObject::offsetExists(mixed index )
 
// 指定 name 获取值
ArrayObject::offsetGet()
 
// 修改指定 name 的值
ArrayObject::offsetSet()
 
// 删除数据
ArrayObject::offsetUnset()
实现例子1：
$array =array('1'=>'one', '2'=>'two', '3'=>'three');
// 构造一个ArrayObject对象
$arrayobject = new ArrayObject($array);
for( $iterator= $arrayobject->getIterator();// 构造一个迭代器   
    $iterator->valid();  // 检查是否还含有元素   
    $iterator->next() ){ // 指向下个元素   
    echo $iterator->key() . ' => ' . $iterator->current() . "\n"; // 打印数组元素
}
输出结果：

1
1 => one 2 => two 3 => three
实现例子2：
$arrayobject =new ArrayObject();
$arrayobject[] = 'zero';
$arrayobject[] = 'one';
$arrayobject[] = 'two';
$iterator= $arrayobject->getIterator();
$iterator->next();
echo $iterator->key()."<br>";
// 重置指针到头部
$iterator->rewind();
echo $iterator->key();
输出结果：
1
0

https://stackoverflow.com/questions/14610307/spl-arrayobject-arrayobjectstd-prop-list/16619183#16619183

http://bobao.360.cn/news/detail/215.html
https://www.anquanke.com/vul/id/1041474

