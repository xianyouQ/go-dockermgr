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
  $scope.people = [] ;
  $scope.mainfilter = '';
  $scope.subfilter = '';
  $scope.rolenames = [];

  $http.get('js/app/management/roles.json').then(function (resp) {
    $scope.roles = resp.data.roles;
    $scope.rolenames = resp.data.names;
  });

  $http.get('js/app/management/bussiness.json').then(function (resp) {
    $scope.subbuses = resp.data.subbuses;
    angular.forEach($scope.subbuses,function(item) {
      newitem ={name: item["mainbus"]};
      if(!$scope.mainbuses.contains(newitem)) {
        $scope.mainbuses.push(newitem);
      }
    })
    $scope.mainbus = $filter('orderBy')($scope.mainbuses, 'name')[0];
    $scope.mainbus.selected = true;
  });
    $http.get('js/app/management/people.json').then(function(resp) {
      angular.forEach(resp.data.people,function(item) {
        if ($scope.roles.length == 0 ) {
          return ;
        }
        item.docker = $scope.roles[item.role].docker;
        item.releaseNew = $scope.roles[item.role].releaseNew;
        item.releaseVerify = $scope.roles[item.role].releaseVerify;
        item.releaseOperation = $scope.roles[item.role].releaseOperation;
        $scope.people.push(item)
      });
       $scope.person = $filter('orderBy')($scope.people, 'name')[0];
    })
  $scope.selectMainBus = function(item){    
    angular.forEach($scope.mainbuses, function(item) {
      item.selected = false;
    });
    $scope.mainbus = item;
    $scope.mainbus.selected = true;
    $scope.mainfilter = item.name;
  };

  $scope.selectSubBus = function(item){    
    angular.forEach($scope.subbuses, function(item) {
      item.selected = false;
    });
    $scope.subbus = item;
    $scope.subbus.selected = true;
    $scope.subfilter = item.name;

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
  $scope.returnMain = function() {
    $scope.mainfilter = '';
    $scope.subfilter = '';
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