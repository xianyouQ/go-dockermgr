app.controller('ManageMentUsersCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) {

  function isObjectValueEqual(a, b) {
    // Of course, we can do it use for in 
    // Create arrays of property names
    var aProps = Object.getOwnPropertyNames(a);
    var bProps = Object.getOwnPropertyNames(b);
 
    // If number of properties is different,
    // objects are not equivalent
    if (aProps.length != bProps.length) {
        return false;
    }
 
    for (var i = 0; i < aProps.length; i++) {
        var propName = aProps[i];
 
        // If values of same property are not equal,
        // objects are not equivalent
        if (a[propName] !== b[propName]) {
            return false;
        }
    }
 
    // If we made it this far, objects
    // are considered equivalent
    return true;
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
  $scope.services = [];
  $scope.filter = new Map();
  $scope.count = [];
  $scope.people = [] ;
  $scope.rolenames = [];


  $http.get("/api/service/count").then(function (resp) {
        if (resp.data.status ){
          for(var i = 0 ;i < resp.data.data ; i++)
          {
            $scope.count.push(i)
          }
      }
      else {
        toaster.pop("error","get count error",resp.data.info);
      } 
  });

  $http.get("/api/service/get").then(function (resp) {
        if (resp.data.status ){
          angular.forEach($resp.data.data,function(service){
            var serviceSplit = service.split("-")
            if(serviceSplit.length != $scope.count.length){
              console.log("invaild service:",service)
              continue
            }
            angularjs.forEach(serviceSplit,function(item,index){
            });
          });
        $scope.services = resp.data.data;
        $scope.selectedservice = $filter('orderBy')($scope.services, 'first')[0];
        $scope.selectedservice.selected = true;
      }
      else {
        toaster.pop("error","get service error",resp.data.info);
      } 
  });

  $scope.isShow = function(idx) {
    if(idx == 0 && ($scope.filter[idx].length == 0)) {
      return true
    } else if ($scope.filter[idx].length == 0 && $scope.filter[idx-1].length > 0) {
      return true
    } 
    return false
  };

  $scope.selectService = function(item){    
    angular.forEach($scope.services, function(item) {
      item.selected = false;
    });
    $scope.mainbus = item;
    $scope.mainbus.selected = true;
    $scope.mainfilter = item.name;
  };

  $scope.deleteUser = function(selectPerson) {
    //$http()
    console.log(selectPerson);
    $scope.people.remove(selectPerson);
  };
  $scope.addUser = function() {
      var modalInstance = $modal.open({
        templateUrl: 'addUserModalContent.html',
        controller: 'addUserModalInstanceCtrl',
        size: 'lg',
        resolve: {
          rolenames: function () {
            return $scope.rolenames;
          }
        }
      });
 
      modalInstance.result.then(function (newUser) {
        newUser.docker = $scope.roles[newUser.role].docker;
        newUser.releaseNew = $scope.roles[newUser.role].releaseNew;
        newUser.releaseVerify = $scope.roles[newUser.role].releaseVerify;
        newUser.releaseOperation = $scope.roles[newUser.role].releaseOperation;
        $scope.people.push(newUser);
      }, function () {
        //log error
      });
  };
  $scope.attachUser = function() {

  };
  $scope.returnUpper = function(idx) {
    $scope.filter[idx-1] = ""
  }
}]);

  app.controller('addUserModalInstanceCtrl', ['$scope', '$modalInstance','rolenames',function($scope, $modalInstance,$rolenames) {
    $scope.rolenames = $rolenames
    $scope.newUser = {"name":"","role":""};
    $scope.ok = function () {
      if ($scope.newUser.name == "" || $scope.newUser.role == ""){
        
      }
      else {
      $modalInstance.close($scope.newUser);
      }
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }])
  ; 