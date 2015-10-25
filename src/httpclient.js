var httpClient = function httpClient(){
"use strict";
	var Get = function(uri){

		var promise = new Promise( function (resolve, reject) {
			var client = new XMLHttpRequest();

			client.onload=function(e){
				if (this.status==200){
					var r = JSON.parse(this.response);
					resolve(r);
				}else{
					reject(this.statusText);
				}

			};

			client.open("GET",uri);
			client.send();

		});

		return promise;

	};

	var Post = function(uri,tempest){
		var promise = new Promise( function (resolve, reject) {

			var client = new XMLHttpRequest();

			client.onload=function(e){
				if (this.status==200){
					var r = JSON.parse(this.response);
					resolve(r);
				}else{
					reject(this.statusText);
				}

			};

			var json = JSON.stringify(tempest);
			client.open("POST",uri);
			client.send(json);

		});

	};

	return {get:Get,post:Post};
};

module.exports = httpClient;
