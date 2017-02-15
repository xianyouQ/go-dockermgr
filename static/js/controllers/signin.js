'use strict';

/* Controllers */
  // signin controller
app.controller('SigninFormController', ['$scope','$rootScope', '$http', '$state','authService' ,function($scope, $rootScope,$http, $state,authService) {
    $scope.user = {};
    $scope.authError = null;
    $scope.login = function() {
      $scope.authError = null;
      // Try to login
      if (authService.returnUser()!== undefined) {
        $scope.authError = "重复登陆";
        $state.go('docker.dashboard');
      }
      $http.post('api/auth/sign', {Username: $scope.user.name, Password: $scope.user.password})
          .then(function(response) {
          if ( !response.data.status ) {
            $scope.authError = response.data.info;
            if (response.data.data != null && "Username" in response.data.data) {
              authService.login(response.data.data);
              if(authService.getLastState() != undefined){
                $state.go(authService.getLastState());
                return;
              }
              $rootScope.Username = authService.returnUser();
              $state.go('docker.dashboard');
            }
          }else{
              authService.login(response.data.data);
              if(authService.getLastState() != undefined){
                $state.go(authService.getLastState());
                return;
              }
              $rootScope.Username = authService.returnUser();
              $state.go('docker.dashboard');
        }
      }, function(x) {
        $scope.authError = 'Server Error';
      });

    };
  }])
;