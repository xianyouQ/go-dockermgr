'use strict';

app.factory('myCache', function ($http,toaster,authService) {
    var services = undefined;
    var count = undefined;
    var idc = undefined;
    //var rolenode = undefined;
    return {

      fresh: function(){
          if(count != undefined && services != undefined && idc != undefined){
            return;
          }
          var auths = authService.getAuths()
          if(auths == undefined){
            return;
          }
          services = [];
          $http.get("/api/service/get").then(function (resp){
            if (resp.data.status){
              if(authService.userHasRole("SYSTEM")){
                services = resp.data.data;
              }
              angular.forEach(resp.data.data,function(service){
                angular.forEach(auths,function(auth){
                  if(String(auth.ServiceAuth.Name).startsWith(service.Code)){
                      services.push(service);
                      return false;
                  } 
                });
              });

            }
            else {
              toaster.pop("error","get service error",resp.data.info);
              services = null;
            }
          });
            $http.get("/api/service/count").then(function (resp) {
                if (resp.data.status ){
                    count = resp.data.data;
                }
                else {
                  toaster.pop("error","get count error",resp.data.info);
                  count = null;
                } 
          });
            $http.get('/api/idc/get').then(function (resp) {
              if (resp.data.status ){
                idc = resp.data.data;
              }
              else {
                toaster.pop("error","get idc error",resp.data.info);
                idc = null;
              } 
          });
      },
      getServices: function(){
        return services;
      },
      getCount: function() {
        return count;
      },
      getIdcs: function() {
        return idc;
      },
      dataOk: function(){
        if(count != undefined && services != undefined && idc != undefined){
          return true;
        }
        return false;
      }
    }
});