app.controller('DockerDashBoardCtrl', ['$scope', '$http', '$filter','$modal','toaster',function($scope, $http, $filter,$modal,toaster) {
    $scope.idcs = [];

    $http.get('/api/docker/dashboard').then(function (resp) {
      if (resp.data.status ){
        $scope.idcs = resp.data.data;
      }
      else {
        toaster.pop("error","get idc error",resp.data.info);
      } 
  });


}]);
