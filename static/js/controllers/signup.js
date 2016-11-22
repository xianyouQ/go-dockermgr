'use strict';

// signup controller
app.controller('SignupFormController', ['$scope', '$http', '$state','authService',function($scope, $http, $state,authService) {
    $scope.user = {};
    $scope.authError = null;
    $scope.signup = function() {
      $scope.authError = null;
      // Try to create
      if (authService.returnUser()!== undefined) {
        $scope.authError = "重复登陆";
        $state.go('app.dashboard-v1');
      }
      if ($scope.user.password !== $scope.user.password1) {
        $scope.authError = "password not match";
        return
      }
      $http.post('api/auth/user', {Username: $scope.user.name, Password: $scope.user.password,Repassword: $scope.user.password1})
      .then(function(response) {
        if ( !response.data.status ) {
          $scope.authError = response.data.info;
        }else{
          $state.go('app.dashboard-v1');
        }
      }, function(x) {
        $scope.authError = 'Server Error';
      });
    };
  }])
 ;