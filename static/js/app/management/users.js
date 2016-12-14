app.controller('ManageMentUsersCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) { 
  $scope.addUser = function() {
      var modalInstance = $modal.open({
        templateUrl: 'addUserModalContent.html',
        controller: 'addUserModalInstanceCtrl',
        size: 'lg'
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
            $modalInstance.close($scope.newUser);
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