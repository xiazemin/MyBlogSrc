---
title: rapidjson
layout: post
category: lang
author: 夏泽民
---
https://github.com/Tencent/rapidjson
http://miloyip.github.io/rapidjson/md_doc_internals.html#IterativeParser
https://github.com/miloyip/nativejson-benchmark

<!-- more -->
https://www.zhihu.com/question/24640264
https://zhuanlan.zhihu.com/p/22457315
{% raw %}
描述JSON的结构JSON的结构有两种，一种是由花括号（{和}）包裹，并用逗号（,）分隔的值键对（key-value-pair）列表；另一种是由中括号（[和]）包裹，并用逗号分隔的值列表。我们将用同一个类JsonValue来表示所有的字面量和这两种结构，并用C#中的string，int和bool来分别描述上述的三种字面量类型。用一个枚举来表示JsonValue内储存的值的类型：    public enum JsonType
    {
        Array,
        Object,
        Int32,
        String,
        Double,
        Boolean,
        Null,
        Symbol
    }
对于JsonValue，我们可以很用一个dictionary和一个list来分别表示给出这样的表示：    public struct JsonValue
    {
            private readonly Dictionary<string, JsonValue> dictionary;
            private readonly List<JsonValue> list;
            public JsonType Type { get; set; }
            public string Value { get; set; }

            public JsonValue(JsonType type) {
                dictionary = null;
                list = null;
                Value = string.Empty;
                    switch (type) {
                        case JsonType.Array:
                            list = new List<JsonValue>();
                            Type = JsonType.Array;
                            break;
                        case JsonType.Object:
                            dictionary = new Dictionary<string, JsonValue>();
                            Type = JsonType.Object;
                            break;
                }
                Type = type;
          }
    }
解析过程这个JSON解析器由接受一个文本流作为输入，输出一个该文本流对应的JsonObject。 解析的过程如下（这是最普通、最通用的结构）：构造一个词法分析器，把文本流转换成Token流构造一个语法分析器，把Token流转换成JsonValue当然，我们可以观察到，一个通常的JSON字面量是不含叶子的。而不含叶子的JsonValue跟一个通常的Token几乎等同。 所以，我们可以在词法分析的时候，就构建JsonValue，并且在语法分析的时候再构建一棵树。优点自然是减少了新对象的生成，使得运行效率大大提高。因为对GC的负担几乎减少了一倍。词法分析（Tokenization）词法分析负责把输入字符流解析成一个个词法单元（Token），以便之后的处理。首先要定义Token。而在上面的讨论中，我们发现可以直接使用JsonValue来代替Token结构。一般的Lexer是一个DFA（有限状态机，下面有实例帮助你迅速理解DFA）。我们可以用很多种方法来表示。比如显式地使用一个表示状态的变量来表示当前的状态。我更喜欢用控制流的变化来表示隐式地表示状态的变化。于是我们可以得到这样的一个NextToken函数，调用一次这个函数，便会返回一个Token。再调用一次，返回下一个Token。于是我们就把字符流变成了Token流。NextToken() return Token {
    c = lookahead()
    if c is whitespace
        Consume()
        return NextToken()
    if c is a char like '{',':','['
        Consume()
        return new Token(c)
    if c is digit
        return GetIntToken()
    if c is '"'
        return GetStringToken()
    if c is 't' or 'f'
        return GetBoolToken()

    return null
}
（伪代码表示，其中lookahead()表示当前字符，但不移动流指针的位置；Consume()也表示当前的字符，但在调用后，会将表示字符流当前位置的指针推进一位）为了让大家更好地理解DFA的概念，我用最简单也最肮脏的GetBoolToken()函数来举例子。我们为了解析true的DFA长这个样子（用于解析false的DFA同理）：转化成代码是这个样子：        private Token? GetBoolToken() {
            var c = PeekChar();
            if (c == 't') {
                GetChar();
                c = PeekChar();
                if (c == 'r') {
                    GetChar();
                    c = PeekChar();
                    if (c == 'u') {
                        GetChar();
                        c = PeekChar();
                        if (c == 'e') {
                            GetChar();
                            return new Token("true", TokenType.BoolType);
                        }
                    }
                }
            }
            if (c == 'f') {
                GetChar();
                c = PeekChar();
                if (c == 'a') {
                    GetChar();
                    c = PeekChar();
                    if (c == 'l') {
                        GetChar();
                        c = PeekChar();
                        if (c == 's') {
                            GetChar();
                            c = PeekChar();
                            if (c == 'e') {
                                GetChar();
                                return new Token("false", TokenType.BoolType);
                            }
                        }
                    }
                }
            }
            throw new FormatException();
        }
当然我们大可不必这么写。我们还可以计算各个token的FIRST集合。然后就可以得到一个丑陋但是快速的Tokenizer。 代码就是这样：        public JsonValue GetNextToken() {
            while (true) {
                var c = cache.Lookahead();
                switch (c) {
                    case ' ':
                    case '\t':
                    case '\n':
                    case '\r':
                        cache.Next();
                        continue;
                    case ':':
                        cache.Next();
                        return JsonValue.Colon;
                    case ',':
                        cache.Next();
                        return JsonValue.Comma;
                    case '{':
                        cache.Next();
                        return JsonValue.LeftBrace;
                    case '}':
                        cache.Next();
                        return JsonValue.RightBrace;
                    case '[':
                        cache.Next();
                        return JsonValue.LeftBracket;
                    case ']':
                        cache.Next();
                        return JsonValue.RightBracket;
                    case 't':
                    case 'n':
                    case 'f':
                        return ParseKeywordToken();
                    case '\"':
                        return ParseStringToken();
                    case '-':
                        return ParseNumberToken();
                    case '0':
                    case '1':
                    case '2':
                    case '3':
                    case '4':
                    case '5':
                    case '6':
                    case '7':
                    case '8':
                    case '9':
                        return ParseNumberToken();
                    default:
                        return JsonValue.Null;
                }
            }
        }
解析其他类型的Token时同理。语法分析（Parsing）得到了Token流之后，接下来就是进行语法分析。JSON文法是LL(1)文法，且文法无左递归，这使得JSON很容易使用递归下降来构造解析器。因此，编写JSON解析器可以作为一个很好的编译器前端的学习材料。 我们可以简单地使用递归下降构建一个LL(1) Parser。比如：    Parse() return JsonValue{
        if lookahead is {
            return ParseJsonObject()
        else if lookahead is [
            return ParseJsonArray()
        throw exception
    }

    ParseJsonObject() return JsonValue{
        Consume("{")
        Loop for parsing key-value pair
        Consume("}")    
    }
以此类推，即可得到一个JSON Parser。就这样结束了？写一个JSON Parser就是这么简单。当然，我也重写了多次才写成现在放在github上的那样的代码。那所有的Parser都是这样简单的吗？并不是。文中的parser还没有提到错误恢复，文法也没有左递归。事实上，一个递归下降的parser不能处理左递归的文法。 现实中的Parser要处理很多问题。比如这样的歧义（List<U, a > b>）。它虽然重要，但是很基础。Parsing不是编译的全部，这只是编译的开始。编译还包括很多阶段，词法分析和语法分析之后还有语义分析，代码优化，代码生成等步骤。当然，随着相关工具和库的成熟，编写一个编译器也没有想以前那样难了。现在我们要写一个编译器，可以直接使用Parser Generater。常见的有ANTLR，bison，flex，VBF等。使用这些工具我们可以很快地生成tokenizer和parser。编译器前端已经相当成熟。 而研究尚未有前端成熟的编译器后端，也有了很多工具可以选择。比如LLVM，DLR，JVM，CLR等等。造一门语言不再是一种困难的事情。可以看到每年都有很多语言横空出世。如果想练习写parser的技能，可以尝试写写下列语言的parser或者解释器（难度逐渐增加）：小学 - JSON - Lisp中学 - Lua - Pascal - C的子集大学 - C - Python地狱 - C++


作者：知乎用户
链接：https://www.zhihu.com/question/24640264/answer/80500016
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

比如词法分析、EBNF乃至AST之类的名词。然而如果我们仅仅是要解析JSON这种数据结构的话，其实用不着那么多东西。要知道JSON的流行很大程度上就是因为它是一个轻量级的、简单易读的数据结构。说个无关的，要知道我第一次看编译原理的书，了解上下文无关文法这个概念的时候，很大程度上是靠了脑子里闪现出来的JSON官网 JSON 里描述的完整的JSON文法。object = {}
       | { members }

members = pair
        | pair , members

pair = string : value

array = []
      | [ elements ]

elements = value 
         | value , elements

value = string
      | number
      | object
      | array
      | true
      | false
      | null

string = ""
       | " chars "

chars = char
      | char chars

char = any-Unicode-character-except-"-or-\-or- control-character
     | \"
     | \\
     | \/
     | \b
     | \f
     | \n
     | \r
     | \t
     | \u four-hex-digits

number = int
       | int frac
       | int exp
       | int frac exp

int = digit
    | digit1-9 digits 
    | - digit
    | - digit1-9 digits

frac = . digits

exp = e digits

digits = digit
       | digit digits

e = e
  | e+
  | e-
  | E
  | E+
  | E-
JSON的优势除了文法很简单之外，还有一个很重要的地方是，我们逐个字符地解析JSON，在不考虑错误处理的情况下，甚至都不需要前瞻！举个例子，在JSON解析的一开始，如果发现是字符'n'，可以立马断定是null；是'{'，就是对象。这给我们解析器的编写带来了很大的方便。其实不管是JSON、XML还是YAML，其核心都是要描绘一种树状类型的结构。原因是生活中的大部分信息都是以这种形式保存的。有树，就会有结点。而不同的标记语言，这个结点里包含的东西也就不同。但是这种树状的结构是基本相同的。把源字符串转换成我们定义好的数据结构的过程，也就是所谓的解析。那么解析的难点在哪里呢？如果一个JSON字符串里的内容仅仅是“true”四个字符，一次strcmp即可完成，你还觉得难吗？是任意数量的空格吗？貌似也不是。不知道题主初学编程时有没有写过那种“统计输入中空格的数量”的程序。如果有，其实排除空格的原理差不多。而且请注意，我们写一个JSON解析器的目的是让杂乱的字符串变成C/Java/Python...等语言能够识别的数据结构。我们所做的工作仅仅是“保存”，没有“解释”，更不是“编译”。外加前面说的根本不需要前瞻一个字符的特性，我们编写JSON Parser其实连词法分析都不用。在这里，我们先定义好一个JSON数据结构的结点类型。struct json_node {
    enum {
        NUMBER,
        STRING,
        ARRAY,
        OBJECT,
        TRUE,
        FALSE,
        NUL,
        UNKNOWN
    } type;
    char *name;
    union {
        double number;
        struct json_node *child_head;
        char *val_str;
    };
};定义好了这个结构体，我们的工作其实就已经完成一半了。根据上面完整定义的文法，再来写解释基本类型的代码。struct json_node *
parse_json(void)
{
    struct json_node *res = malloc(sizeof(struct json_node));
    char ch = getchar();
    switch (ch) {
    case 't':
        while (isalpha(ch)) {
            ch = getchar();
        }
        res->type = TRUE;
        break;
    /* 对于null和false也一样 */
    case '\"':
        res->type = STRING;
        /* 寻找到第一个不是转义字符的双引号 " */
        break;
    case '+':
    case '-':
    case '1':
    /* 1-9 */
        res->type = NUMBER;
        /* 缓冲区存储字符直到停止，然后用系统库函数转换 */
    default:
        res->type = UNKNOWN;
        break;
    }
    return res;
}挺容易明白的。那么为什么JSON解析对于一个之前从未写过的人来说会感觉很难呢？因为它是递归的。递归很强大，很直观，很简洁。但对于不懂递归的人来说，递归往往会使人手足无措。我解析一个true可以，但是我要是要解析未知层数的括号包着的内容呢？没关系，我们解析的函数也递归调用就行了。只要系统的输入流一直在向前走，这没关系。解析object和array的过程类似，只不过我们把词法分析和语法分析弄到了一起，所以跳过空格的过程会比较麻烦：struct json_node *
parse_json(void)
{
    /* 开头的声明 */
    switch (ch) {
    /* 数字和字符串等的解析 */
    case '{':
    {
        char *tmp_name = NULL;
        struct json_node *tmp_child = NULL;
        struct json_node *tmp_tail = NULL;
        res->type = OBJECT;
        res->child_head = NULL;
        while (ch != '}') {
            /* 跳过空格 */
            /* 向后解析字符串 */
            tmp_name = parse_string(ch);
            /* 跳过空格 */
            /* 冒号 */
            ch = getchar();
            current_child = parse_json();
            current_child->name = tmp_name;
            /* 跳过空格 */
            /* 逗号 */
            ch = getchar();
            /* 跳过空格 */
            if (res->child_head == NULL) {
                res->child_head = tmp_tail = tmp_child;
            } else {
                tmp_tail->next = tmp_child;
                tmp_tail = tmp_child;
            }
        }
        break;
    }
    case '[':
    {
        /*
         * 跟解析对象的过程类似
         * 只不过没有名字
         */
    }
    /* 未知情况 */
    }
    return res;
}
{% endraw %}

https://zhuanlan.zhihu.com/p/22457315


