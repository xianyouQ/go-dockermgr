'use strict';

// signup controller
app.controller('PassWdChangeFormController', ['$scope', '$http', '$state','authService',function($scope, $http, $state,authService) {
    $scope.user = {};
    $scope.authError = null;
    $scope.user.Username = authService.returnUser();
    $scope.signup = function() {
      $scope.authError = null;
      // Try to create
      if ($scope.user.Username == undefined) {
        $scope.authError = "请先登录";
        $state.go('access.signin');
      } else 
      if ($scope.user.Password !== $scope.user.Repassword) {
        $scope.authError = "password not match";
        return
      }
      $http.put('/api/auth/passwd/change', $scope.user)
      .then(function(response) {
        if ( !response.data.status ) {
          $scope.authError = response.data.info;
        }else{
          $state.go('docker.dashboard');
        }
      }, function(x) {
        $scope.authError = 'Server Error';
      });
    };
  }])
 ;