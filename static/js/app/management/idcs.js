app.controller('ManageMentIDCsCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) {
    $scope.idcs = [];
    $scope.selectedidc = null;

    $http.get('api/idc').then(function (resp) {
      $scope.idcs = resp.data.data;
      $scope.selectedidc = $filter('orderBy')($scope.idcs, 'first')[0];
      $scope.selectedidc.selected = true;
  });

  $scope.selectIDC = function (item) {
    angular.forEach($scope.idcs, function(item) {
      item.selected = false;
    });
    $scope.selectedidc = item;
    $scope.selectedidc.selected = true;
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
  };

  $scope.createCidr = function () {
      var modalInstance = $modal.open({
        templateUrl: 'addCidrModalContent.html',
        controller: 'addCidrModalInstanceCtrl',
        size: 'lg',
      });
      modalInstance.result.then(function (newCidr) {
        $scope.idcs.Cidrs.push(newIdc);
      }, function () {
        //log error
      });
  };

  $scope.editItem = function (){
    $scope.selectedidc.editing = true;
  };
 $scope.doneEditing = function() {
   $scope.selectedidc.editing = false;
 }

}]);

  app.controller('addIDCModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
   
    $scope.newIdc = {"name":"","code":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
        $http.post('api/idc',{IdcName:$scope.newIdc.name,IdcCode:$scope.newIdc.code}).then(function(response) {
          if (response.data.status ){
            $modalInstance.close(response.data.data);
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
  }]); 

    app.controller('addCidrModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
   
    $scope.newCidr = {"Net":"","StartIP":"","EndIP":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
        $http.post('api/idc',{Net:$scope.newCidr.Net,StartIP:$scope.newCidr.StartIP,EndIP:$scope.newCidr.EndIP}).then(function(response) {
          if (response.data.status ){
            $modalInstance.close($scope.newCidr);
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
  }]); 