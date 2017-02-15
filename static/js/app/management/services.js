app.controller('ManageMentServicesCtrl', ['$scope', '$http', '$filter','$modal','$q','$interval','myCache',function($scope, $http, $filter,$modal,$q,$interval,myCache) {
  function isObjectValueEqual(a, b) {
   if(a.Id === b.Id){
     return true;
   } 
   else {
     return false;
   }
}

  Array.prototype.contains = function(obj) {
    var i = this.length;
    while (i--) {
        if (isObjectValueEqual(this[i],obj)) {
            return true;
        }
    }
    return false;
 }

 Array.prototype.remove=function(obj){ 
  for(var i =0;i <this.length;i++){ 
    var temp = this[i]; 
    if(!isNaN(obj)){ 
      temp=i; 
    } 
    if(isObjectValueEqual(temp,obj)){ 
      for(var j = i;j <this.length;j++){ 
        this[j]=this[j+1]; 
        } 
      this.length = this.length-1; 
      } 
  } 
  }
  
  $scope.mainbuses = [] ;
  $scope.services = new Map();
  $scope.roles = [] ;
  $scope.filter = new Map();
  $scope.count = [];
  $scope.Users = [];
  $scope.padderSelect='conf';
 //$scope.$watch('services',null,true);

    $http.get('/api/auth/get').then(function (resp) {
      if (resp.data.status ){
        angular.forEach(resp.data.data,function(role){
          if(role.NeedAddAuth == true) {
            $scope.roles.push(role);
          }
        });
      }
      else {
        toaster.pop("error","get auth error",resp.data.info);
      } 
  });

  $scope.initData = function(){
    var serviceCount = myCache.getCount();
    console.log(serviceCount);
    if(serviceCount == null){
      return;
    }
    for(var i = 0 ;i < serviceCount ; i++)
    {
      $scope.count.push(i);
      $scope.filter[i]="";
    }
    var tmpservices = myCache.getServices();
    if(tmpservices == null){
      return;
    }
    angular.forEach(tmpservices,function(service){
      var codeSplit = service.Code.split("-")
      if(codeSplit.length != serviceCount){
        console.log("invaild service:",service)
        return true
      }
      var tempService = {Code:""};
      angular.forEach(codeSplit,function(item,index){        
        if($scope.services[index] == undefined) {
          $scope.services[index] = [];
        }
        if(tempService.Code == "") {
          tempService.Code = item
        } else {
          tempService.Code = tempService.Code + "-" + item
        }

        if(!$scope.services[index].contains(tempService) && index < $scope.count.length - 1) {        
          var newService = {Code:""};
          newService.Code = tempService.Code;
          $scope.services[index].push(newService)
        }
      });
    });
    $scope.services[$scope.count.length - 1] = tmpservices;

  }
  var timer = function(){
    return $q(function(resolve,reject){
      myCache.fresh();
      var count = 0;
      var wait = $interval(function() {
        console.log(count);
        if(myCache.dataOk() == true){
          resolve();
          $interval.cancel(wait);
        }
        else {
          count = count + 1;
          if (count > 5){
            reject("timeout");
            $interval.cancel(wait);
          }
        }
      },200);
    })
  }
  timer().then(function(){
     $scope.initData();
  },function(){
  });
  $scope.isShow = function(idx) {
    if (idx < 0 ){
      return false;
    }
    if(idx == 0 && ($scope.filter[idx] == undefined||$scope.filter[idx].length == 0)) {
      return true
    } else if (idx > 0 && $scope.filter[idx] == undefined) {
      return false
    }
    else if ($scope.filter[idx].length == 0 && $scope.filter[idx-1].length > 0) {
      return true
    } 
    return false
  };


  $scope.ConfShow = function() {
    if($scope.selectedService == undefined){
      return false;
    }
    var codeSplit = $scope.selectedService.Code.split("-")
    if(codeSplit.length == $scope.count.length){
      return true
    }
    return false
  }
  $scope.selectService = function(item,idx){
    if (idx == $scope.count.length - 1) {
      $scope.selectedService = item;
      $http.get("/api/auth/auths?serviceId="+$scope.selectedService.Id).then(function(resp){
        $scope.Users = resp.data.data;
      });
      return
    } 
    angular.forEach($scope.services, function(item) {
      item.selected = false;
    });
    $scope.selectedService = item;
    $scope.selectedService.selected = true;
    var serviceSplit = $scope.selectedService.Code.split("-")
    $scope.filter[idx] = serviceSplit[idx];
  };

  $scope.returnUpper = function(idx) {
    $scope.filter[idx-1] = "";
  }
  $scope.commitService = function() {
     $http.post('/api/service/Add',$scope.selectedService).then(function(response) {
          if (response.data.status){
            $scope.selectedService = response.data.data
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
  };
  $scope.createService = function() {
        var modalInstance = $modal.open({
        templateUrl: 'addServiceModalContent.html',
        controller: 'addServiceModalInstanceCtrl',
        size: 'lg',
        resolve: {
          count: function () {
            return $scope.count;
          }
        }
      });
      modalInstance.result.then(function (newService) {
        $scope.selectedService = newService;
         var codeSplit = newService.Code.split("-")
         var tempService = {Code:""};
        angular.forEach(codeSplit,function(item,index){
              if($scope.services[index] == undefined) {
                $scope.services[index] = [];
              }
              if(tempService.Code == "") {
                tempService.Code = item;
              } else {
                tempService.Code = tempService.Code + "-" + item
              }
              if(!$scope.services[index].contains(tempService) && index < $scope.count.length - 1) {
                var newTempService = {Code:""};
                newTempService.Code = tempService.Code;
                $scope.services[index].push(newTempService)
              }
            });
             $scope.services[$scope.count.length - 1].push(newService);
      }, function () {
        //log error
      });
  }
  $scope.deleteService = function(item) {
      var modalInstance = $modal.open({
        templateUrl: 'delConfirmModalContent.html',
        controller: 'delConfirmModalInstanceCtrl',
        size: 'lg',
        resolve: {
          selectedService: function () {
            return $scope.selectedService;
          }
        }
      });
       modalInstance.result.then(function (delService) {
      $scope.services[$scope.count.length -1].remove(delService);
       }, function () {
        //log error
      });
  }

  $scope.addUser = function() {
      var modalInstance = $modal.open({
        templateUrl: 'addAuthModalContent.html',
        controller: 'addAuthModalInstanceCtrl',
        size: 'lg',
        resolve: {
          roles: function () {
            return $scope.roles;
          },
          service: function() {
            return $scope.selectedService;
          },
          othsusers: function() {
            return $scope.Users;
          }
        }
      });
       modalInstance.result.then(function (Users) {
         $scope.Users.push(Users);
       }, function () {
        //log error
      });
  };
  $scope.delUserAuth = function(delUser,serviceAuth) {
      var modalInstance = $modal.open({
        templateUrl: 'delUserAuthConfirmModalContent.html',
        controller: 'delUserAuthConfirmModalInstanceCtrl',
        size: 'lg',
        resolve: {
          delUser: function () {
            return delUser;
          },
          serviceAuth: function(){
            return serviceAuth;
          }
        }
      });
       modalInstance.result.then(function (delUser) {
      $scope.Users.remove(delUser);
       }, function () {
        //log error
      });
  };
}]);
  app.controller('addServiceModalInstanceCtrl', ['$scope', '$modalInstance','$http','count',function($scope, $modalInstance,$http,$count) {
   
    $scope.newService = {"Name":"","Code":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
      if ($scope.newService.Name == "" || $scope.newService.Code == ""){
        return
      }
      var codeSplit = $scope.newService.Code.split("-")
      if (codeSplit.length != $count.length) {
        $scope.formError = "invaild Service Code";
        return
      } 
     $http.post('/api/service/Add',$scope.newService).then(function(response) {
          if (response.data.status){
            $scope.newService = response.data.data
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
      $modalInstance.close($scope.newService);
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }]); 

  app.controller('delUserAuthConfirmModalInstanceCtrl', ['$scope', '$modalInstance','$http','delUser','serviceAuth',function($scope, $modalInstance,$http,$delUser,$serviceAuth) {
   
    $scope.formError = null;
    $scope.confirm="delete User's auth?";
    $scope.ok = function () {
      $scope.formError = null;
     $http.post('/api/auth/delete?serviceId='+$serviceAuth.Service.Id,{User:$delUser,ServiceAuth:$serviceAuth}).then(function(response) {
          if (response.data.status){
            $modalInstance.close($delUser);
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }]); 


  app.controller('delConfirmModalInstanceCtrl', ['$scope', '$modalInstance','$http','selectedService',function($scope, $modalInstance,$http,$selectedService) {
   
    $scope.formError = null;
    $scope.confirm="delete service?";
    $scope.ok = function () {
      $scope.formError = null;
     $http.post('/api/service/Delete',$selectedService).then(function(response) {
          if (response.data.status){
            $modalInstance.close($selectedService);
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }]); 

    app.controller('addAuthModalInstanceCtrl', ['$scope', '$modalInstance','$http','roles','service','othsusers',function($scope, $modalInstance,$http,$roles,$service,$othsusers) {
    $scope.formError = null;
    $scope.Users = [];
    $scope.roles = $roles;
    $scope.selected = {
      Users:[],
      Role: {}
    };
    $http.get("/api/auth/user/get?serviceId="+$service.Id).then(function(resp){
        if(resp.data.status) {
          angular.forEach(resp.data.data,function(item){
            var isSkip = false;
            angular.forEach($othsusers,function(inner){
              if (item.Id == inner.Id) {
                isSkip = true;
                return;
              }
            });
            if (isSkip == false) {
              $scope.Users.push(item);
            }
          })
        } 
        else {
          $scope.formError = resp.data.info;
        }
      });
    $scope.ok = function () {
      $scope.formError = null;
      $http.post('/api/auth/new?serviceId='+$service.Id,{Users:$scope.selected.Users,Service:$service,Role:$scope.selected.Role}).then(function(response) {
          if (response.data.status){
            $modalInstance.close(response.data.data);
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
    };

    $scope.cancel = function () {
      console.log($scope.selected);
      console.log($scope.roles);
      $modalInstance.dismiss('cancel');
    };
  }]); 