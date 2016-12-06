app.controller('ManageMentServicesCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) {
  function isObjectValueEqual(a, b) {
    // Of course, we can do it use for in 
    // Create arrays of property names
    var aProps = Object.getOwnPropertyNames(a);
    var bProps = Object.getOwnPropertyNames(b);
 
    // If number of properties is different,
    // objects are not equivalent
    if (aProps.length != bProps.length) {
        return false;
    }
 
    for (var i = 0; i < aProps.length; i++) {
        var propName = aProps[i];
 
        // If values of same property are not equal,
        // objects are not equivalent
        if (a[propName] !== b[propName]) {
            return false;
        }
    }
 
    // If we made it this far, objects
    // are considered equivalent
    return true;
}

  Array.prototype.contains = function(obj) {
    var i = this.length;
    while (i--) {
        if (isObjectValueEqual(this[i],obj)) {
            return true;
        }
    }
    return false;
 }

 Array.prototype.remove=function(obj){ 
  for(var i =0;i <this.length;i++){ 
    var temp = this[i]; 
    if(!isNaN(obj)){ 
      temp=i; 
    } 
    if(isObjectValueEqual(temp,obj)){ 
      for(var j = i;j <this.length;j++){ 
        this[j]=this[j+1]; 
        } 
      this.length = this.length-1; 
      } 
  } 
  }
  
  $scope.mainbuses = [] ;
  $scope.services = new Map();
  $scope.filter = new Map();
  $scope.count = [];
 //$scope.$watch('services',null,true);
  $http.get("/api/service/count").then(function (resp) {
        if (resp.data.status ){
          for(var i = 0 ;i < resp.data.data ; i++)
          {
            $scope.count.push(i);
            $scope.filter[i]="";
          }
          console.log("count:"+$scope.count)
      }
      else {
        toaster.pop("error","get count error",resp.data.info);
      } 
  });

  $http.get("/api/service/get").then(function (resp) {
        if (resp.data.status ){
          angular.forEach(resp.data.data,function(service){
            var codeSplit = service.Code.split("-")
            if(codeSplit.length != $scope.count.length){
              console.log("invaild service:",service)
              return true
            }
            var tempService = {Code:""};
            angularjs.forEach(codeSplit,function(item,index){
              
              if($scope.services[index] == undefined) {
                $scope.services[index] = [];
              }
              if(tempService.Code == "") {
                tempService.Code = item
              } else {
                tempService.Code = tempService.Code + "-" + item
              }
              if(!$scope.services[index].contains(tempService) && index < $scope.length - 1) {
                $scope.services[index].push(tempService)
              }
            });
          });
          $scope.services[$scope.count.length - 1] = resp.data.data;
          console.log($scope.services)
      }
      else {
        toaster.pop("error","get service error",resp.data.info);
      } 
  });

  $scope.isShow = function(idx) {
    if (idx < 0 ){
      return false;
    }
    if(idx == 0 && ($scope.filter[idx] == undefined||$scope.filter[idx].length == 0)) {
      return true
    } else if (idx > 0 && $scope.filter[idx] == undefined) {
      return false
    }
    else if ($scope.filter[idx].length == 0 && $scope.filter[idx-1].length > 0) {
      return true
    } 
    return false
  };

  $scope.selectService = function(item){    
    angular.forEach($scope.services, function(item) {
      item.selected = false;
    });
    $scope.selectedService = item;
    $scope.selectedService.selected = true;
    var serviceSplit = $scope.selectedService.Code.split("-")
    angular.forEach(serviceSplit,function(item,idx) {
      $scope.filter[idx] = item;
    })
  };

  $scope.returnUpper = function(idx) {
    $scope.filter[idx-1] = ""
  }
  $scope.commitService = function() {
     $http.post('/api/service/Add',$scope.selectedService).then(function(response) {
          if (response.data.status){
            if ($scope.selectedService.Id === 0){
              
            }
            $scope.selectedService = response.data.data
          }
          if  (!response.data.status ) {
            $scope.formError = response.data.info;
          }
        }, function(x) {
        console.log('Server Error')
      });
  };
  $scope.createService = function() {
        var modalInstance = $modal.open({
        templateUrl: 'addServiceModalContent.html',
        controller: 'addServiceModalInstanceCtrl',
        size: 'lg',
        resolve: {
          count: function () {
            return $scope.count;
          }
        }
      });
      modalInstance.result.then(function (newService) {
        $scope.selectedService = newService;
         var codeSplit = newService.Code.split("-")
         var tempService = {Code:""};
        angular.forEach(codeSplit,function(item,index){
            console.log("foreach",item,index);
              if($scope.services[index] == undefined) {
                $scope.services[index] = [];
              }
              if(tempService.Code == "") {
                tempService.Code = item
                console.log("foreach1",item);
              } else {
                tempService.Code = tempService.Code + "-" + item
                console.log("foreach2",item);
              }
              if(!$scope.services[index].contains(tempService) && index < $scope.count.length - 1) {
                console.log("foreach3",tempService);
                $scope.services[index].push(tempService);
              }
            });
             $scope.services[$scope.count.length - 1].push(newService);
             
            console.log($scope.services)
      }, function () {
        //log error
      });
  }
}]);
  app.controller('addServiceModalInstanceCtrl', ['$scope', '$modalInstance','$http','count',function($scope, $modalInstance,$http,$count) {
   
    $scope.newService = {"Name":"","Code":""};
    $scope.formError = null;
    $scope.ok = function () {
      $scope.formError = null;
      if ($scope.newService.Name == "" || $scope.newService.Code == ""){
        return
      }
      var codeSplit = $scope.newService.Code.split("-")
      if (codeSplit.length != $count.length) {
        $scope.formError = "invaild Service Name";
        return
      } 
      $modalInstance.close($scope.newService);
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }]); 