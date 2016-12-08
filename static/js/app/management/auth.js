app.controller('ManageMentAuthCtrl', ['$scope', '$http', '$filter','$modal','toaster',function($scope, $http, $filter,$modal,toaster) {
    $scope.roles = [];
    $scope.selectedrole = null;
    $scope.nodes = [];
    $http.get('/api/auth/get').then(function (resp) {
      if (resp.data.status ){
        $scope.roles = resp.data.data;
        $scope.selectedrole = $filter('orderBy')($scope.auths, 'first')[0];
        $scope.selectedrole.selected = true;
      }
      else {
        toaster.pop("error","get auth error",resp.data.info);
      } 
  });
    $http.get('/api/node/get').then(function (resp) {
      if (resp.data.status ){
        $scope.nodes = resp.data.data;
      }
      else {
        toaster.pop("error","get node error",resp.data.info);
      } 
  });
  $scope.selectRole = function (item) {
    angular.forEach($scope.roles, function(item) {
      item.selected = false;
    });
    angular.forEach(item.Nodes,function(node){
      angular.forEach($scope.nodes,function(innernode){
        if(innernode.Id == node.Id) {
          innernode.Active = true;
        }
      });
    });
    $scope.selectedrole = item;
    $scope.selectedrole.selected = true;
  };
  $scope.createRole = function () {
      var modalInstance = $modal.open({
        templateUrl: 'addRoleModalContent.html',
        controller: 'addRoleModalInstanceCtrl',
        size: 'lg',
      });
 
      modalInstance.result.then(function (newAuth) {
        $scope.roles.push(newAuth);
      }, function () {
        //log error
      });
  };

  $scope.createNode = function () {
      var modalInstance = $modal.open({
        templateUrl: 'addNodeModalContent.html',
        controller: 'addNodeModalInstanceCtrl',
        size: 'lg',
      });
 
      modalInstance.result.then(function (newNode) {
        $scope.nodes.push(newNode);
      }, function () {
        //log error
      });
  };



}]);

  app.controller('addRoleModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
    $scope.newRole = {"Name":"","Status":0,NeedAddAuth:false};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
      if ($scope.newRole.Name == "" || $scope.newRole.Status == 0){
        return
      }
        $http.post('/api/auth/post',$scope.newRole).then(function(response) {
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

    app.controller('addNodeModalInstanceCtrl', ['$scope', '$modalInstance','$http',function($scope, $modalInstance,$http) {
    $scope.newNode = {"Desc":"","Url":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
      if ($scope.newNode.Desc == "" || $scope.newNode.Url == ""){
        return
      }
        $http.post('/api/node/post',$scope.newNode).then(function(response) {
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
