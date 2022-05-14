---
title: antlr graphviz
layout: post
category: lang
author: 夏泽民
---
https://www.jianshu.com/p/37573261d3cf
https://blog.csdn.net/marlonyao/article/details/83816299

https://github.com/antlr/grammars-v4/tree/master/dot
https://www.cnblogs.com/csguo/p/7644277.html
DOT是一种用来描述图结构的声明式语言，用它可以描述网络拓扑图，树结构或者是状态机。（之所以说DOT是一种声明式语言，是因为这种语言只描述图是怎么连接的，而不是描述怎样建立图。）这是一个非常普遍而有用的图形工具，尤其是你的程序需要生成图像的时候。例如，ANTLR的-atn选项就是使用DOT来生成可视化的状态机的。

先举个例子感受下这个语言的用途，比如我们需要将一个有四个函数的程序的调用树进行可视化。当然，我们可以用手在纸上将它画出来，但是，我们可以像下面那样用DOT将它们之间的关系指定出来（不管是手画而是自动生成，都需要从程序源文件中计算出函数之间的调用关系）：

examples/t.dot

digraph G{

    rankdir=LR;

    main [shape=box];

    main -> f -> g;           // main calls f which calls g

    f -> f [style=dotted] ; // f isrecursive

    f -> h;                 // f calls h

}
<!-- more -->
{% raw %}
幸运的是，DOT的参考手册中有我们需要的语法规则，我们几乎可以将它们全部直接引用过来，翻译成ANTLR的语法就行了。不幸的是，我们需要自己指定所有的词法规则。我们不得不通读整个文档以及一些例子，从而找到准确的规则。首先，让我们先从语法规则开始。

DOT的语法规则

下面列出了用ANTLR翻译的DOT参考手册中的核心语法：

examples/DOT.g4

graph : STRICT? (GRAPH | DIGRAPH) id? '{'stmt_list '}' ;

stmt_list : ( stmt ';'? )* ;

stmt : node_stmt

    |edge_stmt

    |attr_stmt

    | id '=' id

    |subgraph

    ;

attr_stmt : (GRAPH | NODE | EDGE) attr_list ;

attr_list : ('[' a_list?']')+ ;

a_list : (id ('=' id)?','?)+ ;

edge_stmt : (node_id | subgraph) edgeRHS attr_list? ;

edgeRHS : ( edgeop (node_id | subgraph) )+ ;

edgeop : '->' | '--';

node_stmt : node_id attr_list? ;

node_id : id port? ;

port : ':' id (':'id)? ;

subgraph : (SUBGRAPH id?)? '{' stmt_list '}' ;

id : ID

    |STRING

    |HTML_STRING

    |NUMBER

    ;

其中，唯一一个和参考手册中语法有点不同的就是port规则。参考手册中是这么定义这个规则的。

port: ':' ID [ ':' compass_pt ]

    | ':' compass_pt

compass_pt

    : (n | ne | e | se| s | sw | w | nw)

如果说指南针参数是关键字而不是合法的变量名，那么这些规则这么写是没问题的。但是，手册中的这句话改变了语法的意思。

注意，指南针参数的值并不是关键字，也就是说指南针参数的那些字符串也可以当作是普通的标识符在任何地方使用…

这意味着我们必须接受像“n ->sw”这样的边语句，而这句话中的n和sw都只是标识符，而不是指南针参数。手册后面还这么说道：“…相反的，编译器需要接受任何标识符。”这句话说的并不明确，但是这句话听起来像是编译器需要将指南针参数也接受为标识符。如果真是这样的话，那么我们也不用去考虑语法中的指南针参数；我们可以直接用id来替换规则中的compass_pt就可以了。

port: ':' id (':'id)? ;

为了验证我们的假设，我们不妨用一些DOT的查看器来尝试下这个假设，比如用Graphviz网站上的一些查看器。事实上，DOT也的确接受下面这样的图的定义，所以我们的port规则是没问题的：

digraph G { n -> sw; }

gcc - graphviz如何将由gcc生成的抽象语法树转储到. dot 文件中？
有两种方法，包括两个步骤

使用GCC内部子对象支持

编译你的代码( 说 test.c )，使用

gcc -fdump-tree-vcg -g test.c

使用任何第三方工具从获取点输出

graph-easy test.c.006t.vcg --as_dot

使用原始转储进行编译，然后用一些脚本对它们进行预处理以形成点文件( 在中这是有用的文章 )

这两种方法都有自己的好和坏的边 --，首先只能在，转换之前获得一个转储。 你可以将任何原始转储转换为点格式，但必须支持脚本，即开销。

更喜欢--的是你自己的选择。

更新：时间是变化的。gcc的全新选项使我们可以立即生成点文件。 只需提供：

复制代码

gcc test.c -fdump-tree-all-graph


你会得到一些已经格式化为你的点文件：

复制代码

test.c.008t.lower.dot


test.c.012t.cfg.dot


test.c.016t.ssa.dot


... etc.. .

https://kb.kutu66.com/others/post_13065198
https://www.it1352.com/784411.html

https://www.it1352.com/980592.html
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.tree.ParseTree;

import java.util.*;

public class Main {

  public static void main(String[] args) {

    /*
        // Expression.g4

        grammar Expression;

        expression
         : '-' expression
         | expression ('*' | '/') expression
         | expression ('+' | '-') expression
         | '(' expression ')'
         | NUMBER
         | VARIABLE
         ;

        NUMBER
         : [0-9]+ ( '.' [0-9]+ )?
         ;

        VARIABLE
         : [a-zA-Z] [a-zA-Z0-9]+
         ;

        SPACE
         : [ \t\r\n] -> skip
         ;
    */

    String source = "3 + 42 * (PI - 3.14159)";

    ExpressionLexer lexer = new ExpressionLexer(CharStreams.fromString(source));
    ExpressionParser parser = new ExpressionParser(new CommonTokenStream(lexer));

    SimpleTree tree = new SimpleTree.Builder()
            .withParser(parser)
            .withParseTree(parser.expression())
            .withDisplaySymbolicName(false)
            .build();

    DotOptions options = new DotOptions.Builder()
            .withParameters("  labelloc=\"t\";\n  label=\"Expression Tree\";\n\n")
            .withLexerRuleShape("circle")
            .build();

    System.out.println(new DotTreeRepresentation().display(tree, options));
  }
}


class DotTreeRepresentation {

  public String display(SimpleTree tree) {
    return display(tree, DotOptions.DEFAULT);
  }

  public String display(SimpleTree tree, DotOptions options) {
    return display(new InOrderTraversal().traverse(tree), options);
  }

  public String display(List<SimpleTree.Node> nodes, DotOptions options) {

    StringBuilder builder = new StringBuilder("graph tree {\n\n");
    Map<SimpleTree.Node, String> nodeNameMap = new HashMap<>();
    int nodeCount = 0;

    if (options.parameters != null) {
      builder.append(options.parameters);
    }

    for (SimpleTree.Node node : nodes) {

      nodeCount++;
      String nodeName = String.format("node_%s", nodeCount);
      nodeNameMap.put(node, nodeName);

      builder.append(String.format("  %s [label=\"%s\", shape=%s];\n",
              nodeName,
              node.getLabel().replace("\"", "\\\""),
              node.isTokenNode() ? options.lexerRuleShape : options.parserRuleShape));
    }

    builder.append("\n");

    for (SimpleTree.Node node : nodes) {

      String name = nodeNameMap.get(node);

      for (SimpleTree.Node child : node.getChildren()) {
        String childName = nodeNameMap.get(child);
        builder.append("  ").append(name).append(" -- ").append(childName).append("\n");
      }
    }

    return builder.append("}\n").toString();
  }
}

class InOrderTraversal {

  public List<SimpleTree.Node> traverse(SimpleTree tree) {

    if (tree == null)
      throw new IllegalArgumentException("tree == null");

    List<SimpleTree.Node> nodes = new ArrayList<>();

    traverse(tree.root, nodes);

    return nodes;
  }

  private void traverse(SimpleTree.Node node, List<SimpleTree.Node> nodes) {

    if (node.hasChildren()) {
      traverse(node.getChildren().get(0), nodes);
    }

    nodes.add(node);

    for (int i = 1; i < node.getChildCount(); i++) {
      traverse(node.getChild(i), nodes);
    }
  }
}

class DotOptions {

  public static final DotOptions DEFAULT = new DotOptions.Builder().build();

  public static final String DEFAULT_PARAMETERS = null;
  public static final String DEFAULT_LEXER_RULE_SHAPE = "box";
  public static final String DEFAULT_PARSER_RULE_SHAPE = "ellipse";

  public static class Builder {

    private String parameters = DEFAULT_PARAMETERS;
    private String lexerRuleShape = DEFAULT_LEXER_RULE_SHAPE;
    private String parserRuleShape = DEFAULT_PARSER_RULE_SHAPE;

    public DotOptions.Builder withParameters(String parameters) {
      this.parameters = parameters;
      return this;
    }

    public DotOptions.Builder withLexerRuleShape(String lexerRuleShape) {
      this.lexerRuleShape = lexerRuleShape;
      return this;
    }

    public DotOptions.Builder withParserRuleShape(String parserRuleShape) {
      this.parserRuleShape = parserRuleShape;
      return this;
    }

    public DotOptions build() {

      if (lexerRuleShape == null)
        throw new IllegalStateException("lexerRuleShape == null");

      if (parserRuleShape == null)
        throw new IllegalStateException("parserRuleShape == null");

      return new DotOptions(parameters, lexerRuleShape, parserRuleShape);
    }
  }

  public final String parameters;
  public final String lexerRuleShape;
  public final String parserRuleShape;

  private DotOptions(String parameters, String lexerRuleShape, String parserRuleShape) {
    this.parameters = parameters;
    this.lexerRuleShape = lexerRuleShape;
    this.parserRuleShape = parserRuleShape;
  }
}

class SimpleTree {

  public static class Builder {

    private Parser parser = null;
    private ParseTree parseTree = null;
    private Set<Integer> ignoredTokenTypes = new HashSet<>();
    private boolean displaySymbolicName = true;

    public SimpleTree build() {

      if (parser == null) {
        throw new  IllegalStateException("parser == null");
      }

      if (parseTree == null) {
        throw new  IllegalStateException("parseTree == null");
      }

      return new SimpleTree(parser, parseTree, ignoredTokenTypes, displaySymbolicName);
    }

    public SimpleTree.Builder withParser(Parser parser) {
      this.parser = parser;
      return this;
    }

    public SimpleTree.Builder withParseTree(ParseTree parseTree) {
      this.parseTree = parseTree;
      return this;
    }

    public SimpleTree.Builder withIgnoredTokenTypes(Integer... ignoredTokenTypes) {
      this.ignoredTokenTypes = new HashSet<>(Arrays.asList(ignoredTokenTypes));
      return this;
    }

    public SimpleTree.Builder withDisplaySymbolicName(boolean displaySymbolicName) {
      this.displaySymbolicName = displaySymbolicName;
      return this;
    }
  }

  public final SimpleTree.Node root;

  private SimpleTree(Parser parser, ParseTree parseTree, Set<Integer> ignoredTokenTypes, boolean displaySymbolicName) {
    this.root = new SimpleTree.Node(parser, parseTree, ignoredTokenTypes, displaySymbolicName);
  }

  public SimpleTree(SimpleTree.Node root) {

    if (root == null)
      throw new IllegalArgumentException("root == null");

    this.root = root;
  }

  public SimpleTree copy() {
    return new SimpleTree(root.copy());
  }

  public String toLispTree() {

    StringBuilder builder = new StringBuilder();

    toLispTree(this.root, builder);

    return builder.toString().trim();
  }

  private void toLispTree(SimpleTree.Node node, StringBuilder builder) {

    if (node.isLeaf()) {
      builder.append(node.getLabel()).append(" ");
    }
    else {
      builder.append("(").append(node.label).append(" ");

      for (SimpleTree.Node child : node.children) {
        toLispTree(child, builder);
      }

      builder.append(") ");
    }
  }

  @Override
  public String toString() {
    return String.format("%s", this.root);
  }

  public static class Node {

    protected String label;
    protected int level;
    protected boolean isTokenNode;
    protected List<SimpleTree.Node> children;

    Node(Parser parser, ParseTree parseTree, Set<Integer> ignoredTokenTypes, boolean displaySymbolicName) {
      this(parser.getRuleNames()[((RuleContext)parseTree).getRuleIndex()], 0, false);
      traverse(parseTree, this, parser, ignoredTokenTypes, displaySymbolicName);
    }

    public Node(String label, int level, boolean isTokenNode) {
      this.label = label;
      this.level = level;
      this.isTokenNode = isTokenNode;
      this.children = new ArrayList<>();
    }

    public void replaceWith(SimpleTree.Node node) {
      this.label = node.label;
      this.level = node.level;
      this.isTokenNode = node.isTokenNode;
      this.children.remove(node);
      this.children.addAll(node.children);
    }

    public SimpleTree.Node copy() {

      SimpleTree.Node copy = new SimpleTree.Node(this.label, this.level, this.isTokenNode);

      for (SimpleTree.Node child : this.children) {
        copy.children.add(child.copy());
      }

      return copy;
    }

    public void normalizeLevels(int level) {

      this.level = level;

      for (SimpleTree.Node child : children) {
        child.normalizeLevels(level + 1);
      }
    }

    public boolean hasChildren() {
      return !children.isEmpty();
    }

    public boolean isLeaf() {
      return !hasChildren();
    }

    public int getChildCount() {
      return children.size();
    }

    public SimpleTree.Node getChild(int index) {
      return children.get(index);
    }

    public int getLevel() {
      return level;
    }

    public String getLabel() {
      return label;
    }

    public boolean isTokenNode() {
      return isTokenNode;
    }

    public List<SimpleTree.Node> getChildren() {
      return new ArrayList<>(children);
    }

    private void traverse(ParseTree parseTree, SimpleTree.Node parent, Parser parser, Set<Integer> ignoredTokenTypes, boolean displaySymbolicName) {

      List<SimpleTree.ParseTreeParent> todo = new ArrayList<>();

      for (int i = 0; i < parseTree.getChildCount(); i++) {

        ParseTree child = parseTree.getChild(i);

        if (child.getPayload() instanceof CommonToken) {

          CommonToken token = (CommonToken) child.getPayload();

          if (!ignoredTokenTypes.contains(token.getType())) {

            String tempText = displaySymbolicName ?
                    String.format("%s: '%s'",
                            parser.getVocabulary().getSymbolicName(token.getType()),
                            token.getText()
                                    .replace("\r", "\\r")
                                    .replace("\n", "\\n")
                                    .replace("\t", "\\t")
                                    .replace("'", "\\'")) :
                    String.format("%s",
                            token.getText()
                                    .replace("\r", "\\r")
                                    .replace("\n", "\\n")
                                    .replace("\t", "\\t"));

            if (parent.label == null) {
              parent.label = tempText;
            }
            else {
              parent.children.add(new SimpleTree.Node(tempText, parent.level + 1, true));
            }
          }
        }
        else {
          SimpleTree.Node node = new SimpleTree.Node(parser.getRuleNames()[((RuleContext)child).getRuleIndex()], parent.level + 1, false);
          parent.children.add(node);
          todo.add(new SimpleTree.ParseTreeParent(child, node));
        }
      }

      for (SimpleTree.ParseTreeParent wrapper : todo) {
        traverse(wrapper.parseTree, wrapper.parent, parser, ignoredTokenTypes, displaySymbolicName);
      }
    }

    @Override
    public String toString() {
      return String.format("{label=%s, level=%s, isTokenNode=%s, children=%s}", label, level, isTokenNode, children);
    }
  }

  private static class ParseTreeParent {

    final ParseTree parseTree;
    final SimpleTree.Node parent;

    private ParseTreeParent(ParseTree parseTree, SimpleTree.Node parent) {
      this.parseTree = parseTree;
      this.parent = parent;
    }
  }
}
{% endraw %}

