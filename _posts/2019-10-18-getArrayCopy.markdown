---
title: ArrayObject getArrayCopy
layout: post
category: lang
author: 夏泽民
---
ArrayObject的使用是说明
ArrayObject是将数组转换为数组对象。
$array =array('1'=>'one', '2'=>'two', '3'=>'three');

$arrayobject = new ArrayObject($array);//构造一个ArrayObject对象
for($iterator= $arrayobject->getIterator();//构造一个迭代器    
$iterator->valid();//检查是否还含有元素    
$iterator->next()){ //指向下个元素    
echo $iterator->key() . ' => ' . $iterator->current() . "\n";//打印数组元素
}
<!-- more -->

public ArrayObject::getArrayCopy ( void ) : array
Exports the ArrayObject to an array.

Parameters ¶
This function has no parameters.

Return Values ¶
Returns a copy of the array. When the ArrayObject refers to an object, an array of the public properties of that object will be returned.

https://www.php.net/manual/en/arrayobject.getarraycopy.php
