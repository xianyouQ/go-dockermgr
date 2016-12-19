app.controller('ManageMentUsersCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) { 


  $scope.users = [];
  $scope.roles = [];
  $scope.userFilter = "";
  $scope.selectedUser = undefined;
  $http.get("/api/auth/user").then(function(resp){
    if (resp.data.status )
    {
      $scope.users = resp.data.data;
    }
    else {
        toaster.pop("error","get user error",resp.data.info);
    }
  });
    $http.get('/api/auth/get').then(function (resp) {
      if (resp.data.status ){
        $scope.roles = resp.data.data;
      }
      else {
        toaster.pop("error","get auth error",resp.data.info);
      } 
  });
  $scope.addUser = function() {
      var modalInstance = $modal.open({
        templateUrl: 'addUserModalContent.html',
        controller: 'addUserModalInstanceCtrl',
        size: 'lg'
      });
 
      modalInstance.result.then(function (newUser) {
        $scope.users.push(newUser);
      }, function () {
        //log error
      });
  };
  $scope.attachUser = function() {

  };
  $scope.selectUser = function(item) {
    $scope.selectedUser = item;
  };
  $scope.isSystem = function() {
    if ($scope.selectedUser == undefined) {
      return false;
    }
    angular.forEach($scope.selectedUser.ServiceAuths,function(item){
      if (item.Name == "SYSTEM") {
        return true;
      }
    });
    return false;
  };
  $scope.commitSystem = function($event) {
    var checkbox = $event.target;  
    var checked = checkbox.checked;
    if(checked != isSystem()) {
      
    }
      
  };
  $scope.commitPassWdReset = function($event){
    var checkbox = $event.target;  
    var checked = checkbox.checked;
    if(checked == true) {

    }
  }
}]);

  app.controller('addUserModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
    $scope.newUser={};
    $scope.ok = function () {
      if ($scope.newUser.Username == ""){
        
      }
      else {
        $http.post('api/auth/user',$scope.newUser).then(function(resp){
          if ( !resp.data.status ) {
            $scope.formError = resp.data.info;
          }else{
            $modalInstance.close(resp.data.data);
          }
        }, function(x) {
          $scope.formError = 'Server Error';
        });
      }
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }])
  ; 