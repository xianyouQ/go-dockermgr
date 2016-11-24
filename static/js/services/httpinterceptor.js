'use strict';

app.factory('myIntercepoter', ['$q','authService','toaster',function($q,authService,toaster) {
    var myIntercepoter = {
    request: function(config){
      return config;
    },
    requestError: function(err){
      return $q.reject(err);
    },
    response: function(res){
      return res;
    },
    responseError: function(err){
      if(-1 === err.status) {
        toaster.pop("error","","Server timeout");
      } else if(401 === err.status) {
        authService.logout();
        toaster.pop("error","","session timeout");
      } else {
        toaster.pop("error",err.status,"request error");
      }
      return $q.reject(err);
    }
    };
    return myIntercepoter;
}]);