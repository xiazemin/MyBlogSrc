---
title: php Reflection
layout: post
category: lang
author: 夏泽民
---
PHP的反射机制提供了一套反射API，用来访问和使用类、方法、属性、参数和注释等，比如可以通过一个对象知道这个对象所属的类，这个类包含哪些方法，这些方法需要传入什么参数，每个参数是什么类型等等，不用创建类的实例也可以访问类的成员和方法，就算类成员定义为 private 也可以在外部访问。
官方文档提供了诸如 ReflectionClass、ReflectionMethod、ReflectionObject、ReflectionExtension 等反射类及相应的API

PHP 5 具有完整的反射 API，添加了对类、接口、函数、方法和扩展进行反向工程的能力。 此外，反射 API 提供了方法来取出函数、类和方法中的文档注释。

请注意部分内部 API 丢失了反射扩展工作所需的代码。
https://www.php.net/manual/zh/book.reflection.php
<!-- more -->
Reflection {
/* 方法 */
public static export ( Reflector $reflector [, bool $return = false ] ) : string
public static getModifierNames ( int $modifiers ) : array
}

Reflection::export — Exports
Reflection::getModifierNames — 获取修饰符的名称

interface Reflector  {
	static function export ();
	function __toString ();
}

abstract class ReflectionFunctionAbstract implements Reflector {

}

class ReflectionFunction extends ReflectionFunctionAbstract implements Reflector {

}

class ReflectionParameter implements Reflector {
}

class ReflectionMethod extends ReflectionFunctionAbstract implements Reflector {
}


class ReflectionClass implements Reflector {
}

class ReflectionObject extends ReflectionClass implements Reflector {
}

class ReflectionProperty implements Reflector {
}

class ReflectionExtension implements Reflector {
}

class ReflectionZendExtension implements Reflector {
}

使用：
class User {
    public $username;
    private $password;
    ｝

创建反射类实例

$refClass = new ReflectionClass(new User('liulu', '123456'));
// 也可以写成 
$refClass = new ReflectionClass('User'); // 将类名User作为参数，建立User类的反射类

反射属性

$properties = $refClass->getProperties(); // 获取User类的所有属性，返回ReflectionProperty的数组
$property = $refClass->getProperty('password'); // 获取User类的password属性
//$properties 结果如下：
Array

反射方法

$methods = $refClass->getMethods(); // 获取User类的所有方法，返回ReflectionMethod数组
$method = $refClass->getMethod('getUsername');  // 获取User类的getUsername方法

//$methods 结果如下：
Array

反射注释

$classComment = $refClass->getDocComment();  // 获取User类的注释文档，即定义在类之前的注释
$methodComment = $refClass->getMethod('setPassowrd')->getDocComment();  // 获取User类中setPassowrd方法的注释
//$classComment 结果如下：
/** * 用户相关类 */
//$methodComment 结果如下：
/** * 设置密码 * @param string $password */
反射实例化

$instance = $refClass->newInstance('admin', 123, '***');  // 从指定的参数创建一个新的类实例
//$instance 结果如下：
User Object ( [username] => admin [password:User:private] => 123 )
注：虽然构造函数中是两个参数，但是newInstance方法接受可变数目的参数，用于传递到类的构造函数。 

$params = ['xiaoming', 'asdfg'];
$instance = $refClass->newInstanceArgs($params); // 从给出的参数创建一个新的类实例

访问、执行类的公有方法——public

$instance->setUsername('admin_1'); // 调用User类的实例调用setUsername方法设置用户名
$username = $instance->getUsername(); // 用过User类的实例调用getUsername方法获取用户名

// 也可以写成
$refClass->getProperty('username')->setValue($instance, 'admin_2'); // 通过反射类ReflectionProperty设置指定实例的username属性值
$username = $refClass->getProperty('username')->getValue($instance); // 通过反射类ReflectionProperty获取username的属性值


// 还可以写成
$refClass->getMethod('setUsername')->invoke($instance, 'admin_3'); // 通过反射类ReflectionMethod调用指定实例的方法，并且传送参数
$value = $refClass->getMethod('getUsername')->invoke($instance); // 通过反射类ReflectionMethod调用指定实例的方法

访问、执行类的非公有方法——private、protected

try {
    // 正确写法
    $property = $refClass->getProperty('password'); // ReflectionProperty Object ( [name] => password [class] => User )
    $property->setAccessible(true); // 修改 $property 对象的可访问性
    $property->setValue($instance, '987654321'); // 可以执行
    $value = $property->getValue($instance); // 可以执行
    echo $value . "\n";   // 输出 987654321

    // 错误写法
    $refClass->getProperty('password')->setAccessible(true); // 临时修改ReflectionProperty对象的可访问性
    $refClass->getProperty('password')->setValue($instance, 'password'); // 不能执行，抛出不能访问异常
    $refClass = $refClass->getProperty('password')->getValue($instance); // 不能执行，抛出不能访问异常
    $refClass = $instance->password;   // 不能执行，类本身的属性没有被修改，仍然是private
} catch (Exception $e){
    echo $e;
}

// 错误写法 结果如下：
ReflectionException: Cannot access non-public member User::password in xxx.php
小结
不管反射类中定义的属性、方法是否为 public，都可以获取到。
直接访问 protected 或则 private 的属性、方法，会抛出异常。
访问非公有成员需要调用指定的 ReflectionProperty 或 ReflectionMethod 对象 setAccessible(true)方法。
