var https = require("https");
var http = require("http");
var cheerio = require("cheerio");
var urlparse = require("url");


var server = http.createServer(function(req, res) {

  readBody(req,res).then(function(url){

    return getContentType(url).then(function(contentType){

    if (contentType.indexOf("html")>=0){
      return getHtmlContent(url)
    }else{
      return Promise.resolve("<body>" + contentType + "</body>")
    }

  });

  }).then(getTextContent).then(function(text){

    console.log(text);
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.write(text);
    res.end();

  }).catch(function(err){
    console.error(err);
  });


});

server.listen("8765", "127.0.0.1", function(e){


});

function readBody(req,res){
  return new Promise(function(resolve,reject){

    var body ="";

    req.on("data",function(data){
        body+=data;
    });

    req.on("end",function(e){
      resolve(body);
    });

  });
}

function getContentType(url){

  return new Promise(function(resolve,reject){

    var options = urlparse.parse(url);
    options.method="HEAD";
    var client;
    if (options.protocol==="https:"){
      client = https;
    }else{
      client = http;
    }

    var req = client.request(options,function(res){
      var contentType = res.headers['content-type'];
      resolve(contentType);
    });

    req.end();
  });

}


function getHtmlContent(url){

  return new Promise(function(resolve,reject){

    var options = urlparse.parse(url);
    options.method="GET";
    var client;

    if (options.protocol==="https:"){
      client = https;
    }else{
      client = http;
    }

    var req = client.request(url,function(res){

      var body = "";
      res.on('data', function(chunk) {
          body += chunk.toString();
      });
      res.on('end', function() {
          resolve(body);
      });
    });

    req.end();

  });
}

function getTextContent(body){

  return new Promise(function(resolve,reject){

    console.log(body);

    var $ = cheerio.load(body);
    var articles = [];

    var contentDivs = [];

    $("div").each(function(i,div){
      var classname = $(div).attr('class');

      if ((typeof(classname) !== "undefined")&&classname.indexOf("content")>=0){
        contentDivs.push(div);
      }

    });

    if (contentDivs.length>0){

      contentDivs.forEach(function(div){

        console.log($(div).text());

        $(div).find("p").each(function(i,article){
          var text = $($(article).contents()[0]).text();
          var text = $(article).text();
          console.log(text);
          articles.push(text);
        });
      });

    }else{
      $('p').each(function (i, article) {

        var text = $($(article).contents()[0]).text();
        articles.push(text);
        //articles.push($(article).text());
      });
    }
/*
    $('p').each(function (i, article) {

      //var text = $($(article).contents()[0]).text();
      //articles.push(text);
      var text = $(article).text();
      console.log(text);
      articles.push(text);
    });
*/
    //articles = articles.filter(x=>(x.search(/[{}<>\[\]]/))===-1)

    articles = articles.reduce(function(acc,cur){
      var lines = cur.split("\n");
      lines.forEach(line=>acc.push(line));
      return acc;
    },[]).map(x=>x.trim()).filter(x=>x.length>11);

    resolve(articles.slice(0,6).join("\n").substr(0,500));
  });


}


function removeTags(str){

    var result = [];
    var contents = /\b\w+[.?!]?\b/g;

    while ((results = contents.exec(str)) !== null) {
        result.push(results[0]);
    }
    return result.join(" ");
}


function everythingButTags(str){
    return str.split(/<[\w/][^<]*?>/).join(" ");
}
