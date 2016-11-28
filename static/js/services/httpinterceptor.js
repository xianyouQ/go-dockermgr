'use strict';

app.factory('myIntercepoter', ['$q','$window','$timeout','authService','toaster',function($q,$window,$timeout,authService,toaster) {
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
      if(0 === err.status) {
        toaster.pop("error","","Server timeout");
      } else if(401 === err.status) {
        authService.logout();
        toaster.pop("error","","session timeout");
        $timeout(function(
        ){
          $window.location.reload();
        },2000);
      } else {
        toaster.pop("error",err.status,"request error");
      }
      return $q.reject(err);
    }
    };
    return myIntercepoter;
}]);