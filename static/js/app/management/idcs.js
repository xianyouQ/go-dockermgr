app.controller('ManageMentIDCsCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) {
    $scope.idcs = [];

    $http.get('api/idc').then(function (resp) {
    console.log(resp);
  });

  $scope.selectIDC = function (item) {

  };
  $scope.createIDC = function () {
      var modalInstance = $modal.open({
        templateUrl: 'addIDCModalContent.html',
        controller: 'addIDCModalInstanceCtrl',
        size: 'lg',
      });
 
      modalInstance.result.then(function (newIdc) {
        $scope.idcs.push(newIdc);
      }, function () {
        //log error
      });
  }
}]);

  app.controller('addIDCModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
   
    $scope.newIdc = {"name":"","code":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
        $http.post('api/idc',{IdcName:$scope.newIdc.name,IdcCode:$scope.newIdc.code}).then(function(response) {
          if (response.data.status ){
            $modalInstance.close($scope.newIdc);
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        $scope.formError = 'Server Error';
      });
      
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }])
  ; 