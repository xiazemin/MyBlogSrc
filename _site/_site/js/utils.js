const Query = {
  parse(search = window.location.search) {
    if (!search) return {}
    const queryString = search[0] === '?' ? search.substring(1) : search
    const query = {}
    queryString.split('&')
      .forEach(queryStr => {
        const [key, value] = queryStr.split('=')
        if (key) query[key] = value
      })

    return query
  },
  stringify(query, prefix = '?') {
    const queryString = Object.keys(query)
      .map(key => `${key}=${encodeURIComponent(query[key] || '')}`)
      .join('&')
    return queryString ? prefix + queryString : ''
  },
}

function ajaxFactory(method) {
  return function(apiPath, data = {}, base = 'https://api.github.com') {
    const req = new XMLHttpRequest()
    const token = localStorage.getItem(LS_ACCESS_TOKEN_KEY)

    let url = `${base}${apiPath}`
    let body = null
    if (method === 'GET' || method === 'DELETE') {
      url += Query.stringify(data)
    }

    const p = new Promise((resolve, reject) => {
      req.addEventListener('load', () => {
        const contentType = req.getResponseHeader('content-type')
        const res = req.responseText
        if (!/json/.test(contentType)) {
          resolve(res)
          return
        }
        const data = req.responseText ? JSON.parse(res) : {}
        if (data.message) {
          reject(new Error(data.message))
        } else {
          resolve(data)
        }
      })
      req.addEventListener('error', error => reject(error))
    })
    req.open(method, url, true)

    req.setRequestHeader('Accept', 'application/vnd.github.squirrel-girl-preview, application/vnd.github.html+json')
    if (token) {
      req.setRequestHeader('Authorization', `token ${token}`)
    }
    if (method !== 'GET' && method !== 'DELETE') {
      body = JSON.stringify(data)
      req.setRequestHeader('Content-Type', 'application/json')
    }

    req.send(body)
    return p
  }
}

function createXHR() {
        if(typeof XMLHttpRequest != "undefined"){
            return new XMLHttpRequest();
        }else if(typeof ActiveXObject != "undefined"){
            if(typeof arguments.callee.activeXString != "string"){
                var versions = ["MSXML2.XMLHttp.6.0", "MSXML2.XMLHttp.3.0", "MSXML2.XMLHttp"];
                for(var i=0, len=versions.length;i < len; i++){
                    try{
                        var xhr = new ActiveXObject(versions[i]);
                        arguments.callee.activeXString = versions[i];
                        return xhr;
                    }catch(ex){//跳过
                            console.log(ex);
                    }
                }
            }
            return new ActiveXObject(arguments.callee.activeXString);
        }else{
            throw new Error("NO XHR object available");
        }
}

function postRequest(){
var xhr = createXHR();
xhr.onreadystatechange = function () {
    console.log(xhr);
    if(xhr.readyState === 4){
        if((xhr.status >= 200 && xhr.status < 300) || xhr.status === 304){
            alert(xhr.responseText);
        }else{
            alert("Response wa unsuccessful: " + xhr.status);
        }
    }
};
xhr.open('POST', "https://github.com/login/oauth/authorize?client_id=981ba8c916c262631ea0", true); //异步请求
xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");//post请求增加的
xhr.send("id=1");
}
