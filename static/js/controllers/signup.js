'use strict';

// signup controller
app.controller('SignupFormController', ['$scope', '$http', '$state', function($scope, $http, $state) {
    $scope.user = {};
    $scope.authError = null;
    $scope.signup = function() {
      $scope.authError = null;
      // Try to create
      if ($scope.user.password !== $scope.user.password1) {
        $scope.authError = "password not match";
        return
      }
      $http.post('api/auth/user', {Username: $scope.user.name, Password: $scope.user.password,Repassword: $scope.user.password1})
      .then(function(response) {
        if ( !response.data.user ) {
          $scope.authError = "new user fail";
        }else{
          $state.go('app.dashboard-v1');
        }
      }, function(x) {
        $scope.authError = 'Server Error';
      });
    };
  }])
 ;