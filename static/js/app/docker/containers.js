app.controller('DockerContainersCtrl', ['$scope', '$http', '$filter','$modal',function($scope, $http, $filter,$modal) {
  function isObjectValueEqual(a, b) {
   if(a.Code === b.Code){
     return true;
   } 
   else {
     return false;
   }
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
  
  $scope.services = new Map();
  $scope.roles = [] ;
  $scope.filter = new Map();
  $scope.count = [];
  $scope.idcs = [];
  $scope.instances = [];
 //$scope.$watch('services',null,true);



  $http.get("/api/service/count").then(function (resp) {
        if (resp.data.status ){
          for(var i = 0 ;i < resp.data.data ; i++)
          {
            $scope.count.push(i);
            $scope.filter[i]="";
          }
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
            angular.forEach(codeSplit,function(item,index){
              
              if($scope.services[index] == undefined) {
                $scope.services[index] = [];
              }
              if(tempService.Code == "") {
                tempService.Code = item
              } else {
                tempService.Code = tempService.Code + "-" + item
              }

              if(!$scope.services[index].contains(tempService) && index < $scope.count.length - 1) {
                
                var newService = {Code:""};
                newService.Code = tempService.Code;
                $scope.services[index].push(newService)
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

    $http.get('/api/idc').then(function (resp) {
      if (resp.data.status ){
        $scope.idcs = resp.data.data;
      }
      else {
        toaster.pop("error","get idc error",resp.data.info);
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


  $scope.ConfShow = function() {
    if($scope.selectedService == undefined){
      return false;
    }
    var codeSplit = $scope.selectedService.Code.split("-")
    if(codeSplit.length == $scope.count.length){
      return true
    }
    return false
  }
  $scope.selectService = function(item,idx){
    if (idx == $scope.count.length - 1) {
      $scope.selectedService = item;
      return
    } 
    angular.forEach($scope.services, function(item) {
      item.selected = false;
    });
    $scope.selectedService = item;
    $scope.selectedService.selected = true;
    var serviceSplit = $scope.selectedService.Code.split("-")
    $scope.filter[idx] = serviceSplit[idx];
  };

  $scope.returnUpper = function(idx) {
    $scope.filter[idx-1] = "";
    $scope.selectedService = undefined;
  }
  $scope.detail = function(idc) {
    $scope.selectedIdc = idc;
    $scope.instances = [];
    $http.post("/api/docker/list",{"Service":$scope.selectedService,"Idc":$scope.selectedIdc,"Scale":0}).then(function(resp){
      if(resp.data.status){
        $scope.instances = resp.data.data;
      }
      else {
        toaster.pop("error","get container error",resp.data.info);
      } 
    },function(){

    })
  }
  $scope.returnIdc = function() {
    $scope.selectedIdc = undefined;
  }
  $scope.scaleContainer = function() {
      var modalInstance = $modal.open({
        templateUrl: 'scaleContianerModalContent.html',
        controller: 'scaleContianerModalInstanceCtrl',
        size: 'lg',
        resolve: {
          selectedIdc: function () {
            return $scope.selectedIdc;
          },
          selectedService: function () {
            return $scope.selectedService;
          }
        }
      });
 
      modalInstance.result.then(function (newAuth) {
      }, function () {
        //log error
      });
  }
}]);

   app.controller('scaleContianerModalInstanceCtrl', ['$scope', '$modalInstance','$http','selectedIdc','selectedService',function($scope, $modalInstance,$http,$selectedIdc,$selectedService) {
    $scope.formError = null;
    $scope.ContainerScaleForm = {"Service":$selectedService,"Idc":$selectedIdc,"Scale":0};
    $scope.ok = function () {
      $scope.formError = null;
      if (isNaN($scope.ContainerCount) == false){
        $scope.formError = "ContainerCount must be a number";
        return
      }
      $scope.ContainerScaleForm.Scale = Number($scope.ContainerScaleForm.Scale);
        $http.post('/api/docker/scale', $scope.ContainerScaleForm).then(function(response) {
          if (response.data.status ){
            console.log(response.data.data);
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