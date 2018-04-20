---
title: urlencode
layout: post
category: cryptology
author: 夏泽民
---
<!-- more -->
为什么要 urlencode()
1.是因为当字符串数据以url的形式传递给web服务器时,字符串中是不允许出现空格和特殊字符的

2.因为 url 对字符有限制，比如把一个邮箱放入 url，就需要使用 urlencode 函数，因为 url 中不能包含 @ 字符。     3.url转义其实也只是为了符合url的规范而已。因为在标准的url规范中中文和很多的字符是不允许出现在url中的。    
看一下php的urlencode的说明:
urlencode — 编码 URL 字符串
string urlencode ( string $str )
返回字符串，此字符串中除了 -_. 之外的所有非字母数字字符都将被替换成百分号（%）后跟两位十六进制数，空格则编码为加号（+）。此编码与 WWW 表单 POST 数据的编码方式是一样的，同时与 application/x-www-form-urlencoded 的媒体类型编码方式一样。由于历史原因，此编码在将空格编码为加号（+）方面与 RFC1738 编码（参见 rawurlencode()）不同。此函数便于将字符串编码并将其用于 URL 的请求部分，同时它还便于将变量传递给下一页。
哪些字符是需要转化的呢？
1. ASCII 的控制字符
这些字符都是不可打印的，自然需要进行转化。
2. 一些非ASCII字符
这些字符自然是非法的字符范围。转化也是理所当然的了。
3. 一些保留字符
很明显最常见的就是“&”了，这个如果出现在url中了，那你认为是url中的一个字符呢，还是特殊的参数分割用的呢？
4. 就是一些不安全的字符了。
例如：空格。为了防止引起歧义，需要被转化为“+”。
明白了这些，也就知道了为什么需要转化了，而转化的规则也是很简单的。

按照每个字符对应的字符编码，不是符合我们范围的，统统的转化为%的形式也就是了。自然也是16进制的形式。

和字符编码无关
通过urlencode的转化规则和目的，我们也很容易的看出，urleocode是基于字符编码的。同样的一个汉字，不同的编码类型，肯定对应不同的urleocode的串。gbk编码的有gbk的encode结果。
apache等服务器，接受到字符串后，可以进行decode，但是还是无法解决编码的问题。编码问题，还是需要靠约定或者字符编码的判断解决。
因此，urleocode只是为了url中一些非ascii字符，可以正确无误的被传输，至于使用哪种编码，就不是encode所关心和解决的问题了。
编码问题，不是urlencode所要解决的。

{% highlight c++ linenos %}
bool UrlEncode(const string& src, string& dst)
{
	if(src.size() == 0)
		return false;
	char hex[] = "0123456789ABCDEF";
    size_t size = src.size();
    for (size_t i = 0; i < size; ++i)
	{
		unsigned char cc = src[i];
		if (isascii(cc))
		{
			//代码这里没写全，这里只转码了空格 / & 这三个字符，实际应用的时候需要补全
			if (cc == ' ')
			{
				dst += "%20";
			}
			else if (cc == '/')
			{
				dst += "%2f";
			}
			else if (cc == '&')
			{
				dst += "%26";
			}
			else
				dst += cc;
		}
		else
		{
			unsigned char c = static_cast<unsigned char>(src[i]);
			dst += '%';
			dst += hex[c / 16];
			dst += hex[c % 16];
		}
	}
	return true;
}
bool UrlDecode(const string& src, string& dst)
{
	if(src.size() == 0)
		return false;
	int hex = 0;
	for (size_t i = 0; i < src.length(); ++i)
	{
		switch (src[i])
		{
			case '+':
				dst += ' ';
				break;
			case '%':
				{
					if (isxdigit(src[i + 1]) && isxdigit(src[i + 2]))
					{
						string hexStr = src.substr(i + 1, 2);
						hex = strtol(hexStr.c_str(), 0, 16);
						//字母和数字[0-9a-zA-Z]、一些特殊符号[$-_.+!*'(),] 、以及某些保留字[$&+,/:;=?@]
						//可以不经过编码直接用于URL
						if (!((hex >= 48 && hex <= 57) || //0-9
						(hex >=97 && hex <= 122) ||   //a-z
						(hex >=65 && hex <= 90) ||    //A-Z
						//一些特殊符号及保留字[$-_.+!*'(),]  [$&+,/:;=?@]
						hex == 0x21 || hex == 0x24 || hex == 0x26 || hex == 0x27 || hex == 0x28 || hex == 0x29
						|| hex == 0x2a || hex == 0x2b|| hex == 0x2c || hex == 0x2d || hex == 0x2e || hex == 0x2f
						|| hex == 0x3A || hex == 0x3B|| hex == 0x3D || hex == 0x3f || hex == 0x40 || hex == 0x5f
						))
						{
							dst += char(hex);
							i += 2;
						}
						else
							dst += '%';
					}
					else
						dst += '%';
				}
				break;
			default:
				dst += src[i];
				break;
		}
	}
	return true;
}
{% endhighlight %}