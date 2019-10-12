---
title: radix tree router
layout: post
category: golang
author: 夏泽民
---
golang的web框架echo和gin都使用了radix tree作为路由查找的算法
https://github.com/labstack/echo
https://github.com/gin-gonic/gin
<!-- more -->
在gin的路由中，每一个Http Method(GET, PUT, POST…)都对应了一棵 radix tree

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
    // ...
    
    // 获取method对应的树，如果没有就创建
    root := engine.trees.get(method)
    if root == nil {
        // 创建radix tree，只有根节点
        root = new(node)
        engine.trees = append(engine.trees, methodTree{method: method, root: root})
    }
    root.addRoute(path, handlers)
}

radix tree可以被认为是一棵简洁版的前缀树。拥有共同前缀的节点也共享同一个父节点。下面是一个GET方法对应的路由树的结构：

Priority   Path             Handle
9          \                *<1>
3          ├s               nil
2          |├earch\         *<2>
1          |└upport\        *<3>
2          ├blog\           *<4>
1          |    └:post      nil
1          |         └\     *<5>
2          ├about-us\       *<6>
1          |        └team\  *<7>
1          └contact\        *<8>
*<num>是方法（handler）对应的指针，从根节点遍历到叶子节点我们就能得到完整的路由表，图中的示例实现了以下路由：

GET("/", func1)
GET("/search/", func2)
GET("/support/", func3)
GET("/blog/", func4)
GET("/blog/:post/", func5)
GET("/about-us/", func6)
GET("/about-us/team/", func7)
GET("/contact/", func8)
:post是真实的post name的一个占位符（就是一个参数）。这里体现了radix tree相较于hash-map的一个优点，树结构允许我们的路径中存在动态的部分（参数）,因为我们匹配的是路由的模式而不是hash值

为了更具扩展性，每一层的节点按照priority排序，priority是节点的子节点（儿子节点，孙子节点等）注册的handler的数量，这样做有两个好处：

被最多路径包含的节点会被最先评估。这样可以让尽量多的路由快速被定位。
有点像成本补偿。最长的路径可以被最先评估，补偿体现在最长的路径需要花费更长的时间来定位，如果最长路径的节点能被优先评估（即每次拿子节点都命中），那么所花时间不一定比短路径的路由长。下面展示了节点（每个-可以看做一个节点）评估的路径：从左到右，从上到下
├------------
├---------
├-----
├----
├--
├--
└-
节点数据结构

节点的数据结构如下：

type node struct {
    // 节点路径，比如上面的s，earch，和upport
    path      string
    // 节点是否是参数节点，比如上面的:post
    wildChild bool
    // 节点类型，包括static, root, param, catchAll
    // static: 静态节点，比如上面的s，earch等节点
    // root: 树的根节点
    // catchAll: 有*匹配的节点
    // param: 参数节点
    nType     nodeType
    // 路径上最大参数个数
    maxParams uint8
    // 和children字段对应, 保存的是分裂的分支的第一个字符
    // 例如search和support, 那么s节点的indices对应的"eu"
    // 代表有两个分支, 分支的首字母分别是e和u
    indices   string
    // 儿子节点
    children  []*node
    // 处理函数
    handlers  HandlersChain
    // 优先级，子节点注册的handler数量
    priority  uint32
}
添加路由
func (n *node) addRoute(path string, handlers HandlersChain) {
    fullPath := path
    n.priority++
    numParams := countParams(path)
    // non-empty tree
    if len(n.path) > 0 || len(n.children) > 0 {
    walk:
        for {
            // Update maxParams of the current node
            if numParams > n.maxParams {
                n.maxParams = numParams
            }
            // Find the longest common prefix.
            // This also implies that the common prefix contains no ':' or '*'
            // since the existing key can't contain those chars.
            i := 0
            max := min(len(path), len(n.path))
            for i < max && path[i] == n.path[i] {
                i++
            }
            // Split edge
            // 开始分裂，比如一开始path是search，新来了support，s是他们匹配的部分，
            // 那么会将s拿出来作为parent节点，增加earch和upport作为child节点
            if i < len(n.path) {
                child := node{
                    path:      n.path[i:],  // 不匹配的部分作为child节点
                    wildChild: n.wildChild,
                    indices:   n.indices,
                    children:  n.children,
                    handlers:  n.handlers,
                    priority:  n.priority - 1,  // 降级成子节点，priority减1
                }
                // Update maxParams (max of all children)
                for i := range child.children {
                    if child.children[i].maxParams > child.maxParams {
                        child.maxParams = child.children[i].maxParams
                    }
                }
                
                // 当前节点的子节点变成刚刚分裂的出来的节点
                n.children = []*node{&child}
                // []byte for proper unicode char conversion, see #65
                n.indices = string([]byte{n.path[i]})
                n.path = path[:i]
                n.handlers = nil
                n.wildChild = false
            }
            // Make new node a child of this node
            // 将新来的节点插入新的parent节点作为子节点
            if i < len(path) {
                path = path[i:]
					// 如果是参数节点（包含:或*）
                if n.wildChild {
                    n = n.children[0]
                    n.priority++
                    // Update maxParams of the child node
                    if numParams > n.maxParams {
                        n.maxParams = numParams
                    }
                    numParams--
                    // Check if the wildcard matches
                    // 例如：/blog/:pp 和 /blog/:ppp，需要检查更长的通配符
                    if len(path) >= len(n.path) && n.path == path[:len(n.path)] {
                        // check for longer wildcard, e.g. :name and :names
                        if len(n.path) >= len(path) || path[len(n.path)] == '/' {
                            continue walk
                        }
                    }
                    panic("path segment '" + path +
                        "' conflicts with existing wildcard '" + n.path +
                        "' in path '" + fullPath + "'")
                }
					// 首字母，用来与indices做比较
                c := path[0]
                // slash after param
                if n.nType == param && c == '/' && len(n.children) == 1 {
                    n = n.children[0]
                    n.priority++
                    continue walk
                }
                // Check if a child with the next path byte exists
                // 判断子节点中是否有和当前path有匹配的，只需要查看子节点path的第一个字母即可，即indices
                // 比如s的子节点现在是earch和upport，indices为eu
                // 如果新来的路由为super，那么就是和upport有匹配的部分u，将继续分类现在的upport节点
                for i := 0; i < len(n.indices); i++ {
                    if c == n.indices[i] {
                        i = n.incrementChildPrio(i)
                        n = n.children[i]
                        continue walk
                    }
                }
                // Otherwise insert it
                if c != ':' && c != '*' {
                    // []byte for proper unicode char conversion, see #65
                    // 记录第一个字符，放在indices中
                    n.indices += string([]byte{c})
                    child := &node{
                        maxParams: numParams,
                    }
                    // 增加子节点
                    n.children = append(n.children, child)
                    n.incrementChildPrio(len(n.indices) - 1)
                    n = child
                }
                n.insertChild(numParams, path, fullPath, handlers)
                return
            } else if i == len(path) { // Make node a (in-path) leaf
                // 路径相同，如果已有handler就报错，没有就赋值
                if n.handlers != nil {
                    panic("handlers are already registered for path ''" + fullPath + "'")
                }
                n.handlers = handlers
            }
            return
        }
    } else { // Empty tree，空树，插入节点，节点种类是root
        n.insertChild(numParams, path, fullPath, handlers)
        n.nType = root
    }
}
此函数的主要目的是找到插入节点的位置，如果和现有节点存在相同的前缀，那么要将现有节点进行分裂，然后再插入，下面是insertChild函数

插入子节点
// @1: 参数个数
// @2: 路径
// @3: 完整路径
// @4: 处理函数
func (n *node) insertChild(numParams uint8, path string, fullPath string, handlers HandlersChain) {
    var offset int // already handled bytes of the path
    // find prefix until first wildcard (beginning with ':'' or '*'')
    // 找到前缀，只要匹配到wildcard
    for i, max := 0, len(path); numParams > 0; i++ {
        c := path[i]
        if c != ':' && c != '*' {
            continue
        }
        // find wildcard end (either '/' or path end)
        end := i + 1
        for end < max && path[end] != '/' {
            switch path[end] {
            // the wildcard name must not contain ':' and '*'
            case ':', '*':
                panic("only one wildcard per path segment is allowed, has: '" +
                    path[i:] + "' in path '" + fullPath + "'")
            default:
                end++
            }
        }
        // check if this Node existing children which would be
        // unreachable if we insert the wildcard here
        if len(n.children) > 0 {
            panic("wildcard route '" + path[i:end] +
                "' conflicts with existing children in path '" + fullPath + "'")
        }
        // check if the wildcard has a name
        if end-i < 2 {
            panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
        }
        if c == ':' { // param
            // split path at the beginning of the wildcard
            if i > 0 {
                n.path = path[offset:i]
                offset = i
            }
            child := &node{
                nType:     param,
                maxParams: numParams,
            }
            n.children = []*node{child}
            n.wildChild = true
            n = child
            n.priority++
            numParams--
            // if the path doesn't end with the wildcard, then there
            // will be another non-wildcard subpath starting with '/'
            if end < max {
                n.path = path[offset:end]
                offset = end
                
                child := &node{
                    maxParams: numParams,
                    priority:  1,
                }
                n.children = []*node{child}
                // 下次循环这个新的child节点
                n = child
            }
        } else { // catchAll
            if end != max || numParams > 1 {
                panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
            }
            if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
                panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
            }
            // currently fixed width 1 for '/'
            i--
            if path[i] != '/' {
                panic("no / before catch-all in path '" + fullPath + "'")
            }
            n.path = path[offset:i]
            // first node: catchAll node with empty path
            child := &node{
                wildChild: true,
                nType:     catchAll,
                maxParams: 1,
            }
            n.children = []*node{child}
            n.indices = string(path[i])
            n = child
            n.priority++
            // second node: node holding the variable
            child = &node{
                path:      path[i:],
                nType:     catchAll,
                maxParams: 1,
                handlers:  handlers,
                priority:  1,
            }
            n.children = []*node{child}
            return
        }
    }
    // insert remaining path part and handle to the leaf
    n.path = path[offset:]
    n.handlers = handlers
}
insertChild函数是根据path本身进行分割, 将/分开的部分分别作为节点保存, 形成一棵树结构. 注意参数匹配中的:和*的区别, 前者是匹配一个字段, 后者是匹配后面所有的路径

路径查找
匹配每个children的path，最长匹配

// Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getValue(path string, po Params, unescape bool) (handlers HandlersChain, p Params, tsr bool) {
    p = po
walk: // Outer loop for walking the tree
    for {
        // 尚未到达path的终点
        if len(path) > len(n.path) {
            // 前面一段需要一致
            if path[:len(n.path)] == n.path {
                path = path[len(n.path):]
                // If this node does not have a wildcard (param or catchAll)
                // child,  we can just look up the next child node and continue
                // to walk down the tree
                if !n.wildChild {
                    c := path[0]
                    for i := 0; i < len(n.indices); i++ {
                        if c == n.indices[i] {
                            n = n.children[i]
                            continue walk
                        }
                    }
                    // Nothing found.
                    // We can recommend to redirect to the same URL without a
                    // trailing slash if a leaf exists for that path.
                    tsr = (path == "/" && n.handlers != nil)
                    return
                }
                // handle wildcard child
                n = n.children[0]
                switch n.nType {
                case param:
                    // find param end (either '/' or path end)
                    end := 0
                    for end < len(path) && path[end] != '/' {
                        end++
                    }
                    // save param value
                    if cap(p) < int(n.maxParams) {
                        p = make(Params, 0, n.maxParams)
                    }
                    i := len(p)
                    p = p[:i+1] // expand slice within preallocated capacity
                    p[i].Key = n.path[1:]
                    val := path[:end]
                    if unescape {
                        var err error
                        if p[i].Value, err = url.QueryUnescape(val); err != nil {
                            p[i].Value = val // fallback, in case of error
                        }
                    } else {
                        p[i].Value = val
                    }
                    // we need to go deeper!
                                        if end < len(path) {
                        if len(n.children) > 0 {
                            path = path[end:]
                            n = n.children[0]
                            continue walk
                        }
                        // ... but we can't
                        tsr = (len(path) == end+1)
                        return
                    }
                    if handlers = n.handlers; handlers != nil {
                        return
                    }
                    if len(n.children) == 1 {
                        // No handle found. Check if a handle for this path + a
                        // trailing slash exists for TSR recommendation
                        n = n.children[0]
                        tsr = (n.path == "/" && n.handlers != nil)
                    }
                    return
                case catchAll:
                    // save param value
                    if cap(p) < int(n.maxParams) {
                        p = make(Params, 0, n.maxParams)
                    }
                    i := len(p)
                    p = p[:i+1] // expand slice within preallocated capacity
                    p[i].Key = n.path[2:]
                    if unescape {
                        var err error
                        if p[i].Value, err = url.QueryUnescape(path); err != nil {
                            p[i].Value = path // fallback, in case of error
                        }
                    } else {
                        p[i].Value = path
                    }
                    handlers = n.handlers
                    return
                default:
                    panic("invalid node type")
                }
            }
        } else if path == n.path {
            // We should have reached the node containing the handle.
            // Check if this node has a handle registered.
            if handlers = n.handlers; handlers != nil {
                return
            }
            if path == "/" && n.wildChild && n.nType != root {
                tsr = true
                return
            }
            // No handle found. Check if a handle for this path + a
            // trailing slash exists for trailing slash recommendation
            for i := 0; i < len(n.indices); i++ {
                if n.indices[i] == '/' {
                    n = n.children[i]
                    tsr = (len(n.path) == 1 && n.handlers != nil) ||
                        (n.nType == catchAll && n.children[0].handlers != nil)
                    return
                }
            }
            return
        }
        // Nothing found. We can recommend to redirect to the same URL with an
        // extra trailing slash if a leaf exists for that path
        tsr = (path == "/") ||
            (len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
                path == n.path[:len(n.path)-1] && n.handlers != nil)
        return
    }
}
