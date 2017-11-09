app.controller('DockerContainersCtrl', ['$scope', '$http', '$filter','$modal','$interval','$q','myCache','toaster',function($scope, $http, $filter,$modal,$interval,$q,myCache,toaster) {
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
  $scope.instances = new Map();
  $scope.instancekey = undefined;
  $scope.scaleTask = new Map();
  $scope.scaleTask["container_num_new"] = 0;
 //$scope.$watch('services',null,true);
  $scope.initData = function(){
    var serviceCount = myCache.getCount();
    console.log(serviceCount);
    if(serviceCount == null){
      return;
    }
    for(var i = 0 ;i < serviceCount ; i++)
    {
      $scope.count.push(i);
      $scope.filter[i]="";
    }
    var tmpservices = myCache.getServices();
    if(tmpservices == null){
      return;
    }
    angular.forEach(tmpservices,function(service){
      var codeSplit = service.Code.split("-")
      if(codeSplit.length != serviceCount){
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
    $scope.services[$scope.count.length - 1] = tmpservices;
    $scope.idcs = myCache.getIdcs();
    if($scope.idcs == null){
      return;
    }
  }
  var timer = function(){
    return $q(function(resolve,reject){
      myCache.fresh();
      var count = 0;
      var wait = $interval(function() {
        console.log(count);
        if(myCache.dataOk() == true){
          resolve();
          $interval.cancel(wait);
        }
        else {
          count = count + 1;
          if (count > 5){
            reject("timeout");
            $interval.cancel(wait);
          }
        }
      },200);
    })
  }
  timer().then(function(){
     $scope.initData();
  },function(){

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
      $scope.returnIdc();
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
    console.log("return idc");
    $scope.filter[idx-1] = "";
    $scope.selectedService = undefined;
  }
  $scope.detail = function(idc) {
    console.log("set idc");
    $scope.selectedIdc = idc;
    $scope.instancekey=$scope.selectedService.Id+ "_"+$scope.selectedIdc.Id;
    if($scope.instances[$scope.instancekey] != undefined){
       return;
    }
    $http.get("/api/docker/list?serviceId="+$scope.selectedService.Id+"&idcId="+$scope.selectedIdc.Id).then(function(resp){
      if(resp.data.status){
        $scope.instances[$scope.instancekey] = resp.data.data;
      }
      else {
        toaster.pop("error","get container error",resp.data.info);
      } 
    },function(){

    })
  }
  $scope.returnIdc = function() {
    console.log("return idc");
    $scope.selectedIdc = undefined;
  }
  $scope.scaleContainer = function() {
    console.log($scope.scaleTask)
   if(isNaN($scope.scaleTask.container_num_new) == true){
     toaster.pop("error","input error","ContainerCount must be a number")
     return;
   }
   console.log($scope.selectedService);
   console.log($scope.selectedIdc);
   $http.get('/api/docker/scale?serviceId='+$scope.selectedService.Id+"&idcId="+$scope.selectedIdc.Id+'&scaleCount='+$scope.scaleTask.container_num_new).then(function(response) {
    if (response.data.status ){
      console.log(response.data.data);
    }
    if  (!response.data.status ) {
      toaster.pop("error","get idc error",resp.data.info);
    }
  }, function(x) {
    toaster.pop("error","get idc error",null);
});
  }
}]);

