app.controller('ManageMentIDCsCtrl', ['$scope', '$http', '$filter','$modal','toaster',function($scope, $http, $filter,$modal,toaster) {
    $scope.idcs = [];
    $scope.selectedidc = null;
    $scope.padderSelect='conf';

    $http.get('/api/idc/get').then(function (resp) {
      if (resp.data.status ){
        $scope.idcs = resp.data.data;
      }
      else {
        toaster.pop("error","get idc error",resp.data.info);
      } 
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
        resolve: {
          selectedidc: function () {
            return $scope.selectedidc;
          }
        }
      });
      modalInstance.result.then(function (newCidr) {
        $scope.selectedidc.Cidrs.push(newCidr);
      }, function () {
        //log error
      });
  };

 $scope.commitMarathonConf = function () {
   $scope.MarathonformError = null;
    $http.post('/api/marathon/conf',$scope.selectedidc).then(function(response) {
          if (response.data.status ){
            $scope.selectedidc.MarathonSerConf = response.data.data
          }
          if  (!response.data.status ) {
            $scope.MarathonformError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
 };


 $scope.commitRegistryConf = function () {
   $scope.RegistryformError = null;
    $http.post('/api/registry/conf',$scope.selectedidc).then(function(response) {
          if (response.data.status ){
            $scope.selectedidc.RegistryConf = response.data.data
          }
          if  (!response.data.status ) {
            $scope.RegistryformError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
 };

 $scope.delIdc = function(delIdc) {
      var modalInstance = $modal.open({
        templateUrl: 'delIdcConfirmModalContent.html',
        controller: 'delIdcConfirmModalInstanceCtrl',
        size: 'lg',
        resolve: {
          delIdc: function () {
            return delIdc;
          }
        }
      });
       modalInstance.result.then(function (delIdc) {
      $scope.idcs.remove(delIdc);
       }, function () {
        //log error
      });
 }
 $scope.delCidr = function(delCidr) {
      var modalInstance = $modal.open({
        templateUrl: 'delCidrConfirmModalContent.html',
        controller: 'delCidrConfirmModalInstanceCtrl',
        size: 'lg',
        resolve: {
          delCidr: function () {
            return delCidr;
          }
        }
      });
       modalInstance.result.then(function (delCidr) {
        $scope.selectedidc.Cidrs.remove(delCidr);
       }, function () {
        //log error
      });
 }

}]);

  app.controller('addIDCModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
   
    $scope.newIdc = {"IdcName":"","IdcCode":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
      if ($scope.newIdc.IdcName == "" || $scope.newIdc.IdcCode == ""){
        return
      }
        $http.post('/api/idc',$scope.newIdc).then(function(response) {
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

    app.controller('addCidrModalInstanceCtrl', ['$scope', '$modalInstance','$http','selectedidc',function($scope, $modalInstance,$http,$selectedidc) {
   
    $scope.newCidr = {"Net":"","StartIP":"","EndIP":""};
    $scope.newCidr.BelongIdc = $selectedidc
    $scope.formError = null;
    $scope.ok = function () {
    $scope.formError = null;
    if ($scope.newCidr.Net == "" || $scope.newCidr.StartIP == "" || $scope.newCidr.EndIP == ""){
          return
      }
        $http.post('/api/Cidr',$scope.newCidr).then(function(response) {
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

   app.controller('delIdcConfirmModalInstanceCtrl', ['$scope', '$modalInstance','$http','delIdc',function($scope, $modalInstance,$http,$delIdc) {
   
    $scope.formError = null;
    $scope.confirm="delete Idc?";
    $scope.ok = function () {
      $scope.formError = null;
     $http.delete('/api/idc?idcId='+$delIdc.Id).then(function(response) {
          if (response.data.status){
            $modalInstance.close($delIdc);
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }]); 

   app.controller('delCidrConfirmModalInstanceCtrl', ['$scope', '$modalInstance','$http','delCidr',function($scope, $modalInstance,$http,$delCidr) {
   
    $scope.formError = null;
    $scope.confirm="delete Cidr?";
    $scope.ok = function () {
      $scope.formError = null;
     $http.delete('/api/Cidr?cidrId='+$delCidr.Id).then(function(response) {
          if (response.data.status){
            $modalInstance.close($delCidr);
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }]); 