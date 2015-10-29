var httpClient = function httpClient(){
"use strict";
	var Get = function(uri){

		var promise = new Promise( function (resolve, reject) {
			var client = new XMLHttpRequest();

			client.onload=function(e){
				if (this.status==200){
					resolve(this.response);
				}else{
					reject(this.statusText);
				}

			};

			client.open("GET",uri);
			client.send();

		});

		return promise;

	};

	var Post = function(uri,body,options){
		body=body||"";
		options=options||{};
		var promise = new Promise( function (resolve, reject) {
			var client = new XMLHttpRequest();

			if (typeof(options.params)==="object"){
				Object.keys(options.params).forEach(function(key, index){

					if (index===0){
						uri+="?";
					}else{
						uri+="&";
					}
					uri+=key+"="+options.params[key];
				});
			}

			client.open("POST",uri);

			if (typeof(options.headers)==="object"){
				Object.keys(options.headers).forEach(function(key){
					client.setRequestHeader(key,options.headers[key]);
				});
			}

			client.onload=function(e){
				if (this.status==200){
					resolve(this.response);
				}else{
					reject(this.statusText);
				}
			};

			client.send(body);

		});

		return promise;

	};

	return {get:Get,post:Post};
};

module.exports = httpClient;
