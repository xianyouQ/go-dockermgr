app.controller('ReleaseServiceCtrl', ['$scope', '$http', '$filter','$modal','toaster',function($scope, $http, $filter,$modal,toaster) {
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
  $scope.releaseConf = {};
  $scope.padderSelect = "list";
  $scope.releases = [];
  $scope.newReleaseTask = {};
  $scope.selectedTask = undefined;



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
      $scope.releaseConf = {};
      $scope.changepadder("list");
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

  $scope.commitDeploySettings = function(){
    var enableidcs = [];
    angular.forEach($scope.idcs,function(idc){
      if (idc.enableRelease) {
        enableidcs.push(idc);
      }
    });
    $scope.releaseConf.Service = $scope.selectedService
    $scope.releaseConf.ReleaseIdc = enableidcs
    $scope.releaseConf.FaultTolerant = Number($scope.releaseConf.FaultTolerant)
    $scope.releaseConf.IdcParalle = Number($scope.releaseConf.IdcParalle)
    $scope.releaseConf.IdcInnerParalle = Number($scope.releaseConf.IdcInnerParalle)
    $scope.releaseConf.TimeOut = Number($scope.releaseConf.TimeOut)
    $http.post("/api/release/conf",$scope.releaseConf).then(function(resp){
      if(resp.data.status) {
        $scope.releaseConf = resp.data.data;
      }
      else {
        toaster.pop("error","update conf error",resp.data.info);
      }
    },function(){
      
    })
  }

  $scope.changepadder = function(selectPadder) {
    $scope.padderSelect = selectPadder;
    if(selectPadder == "list") {
      $http.post("/api/release/task",$scope.selectedService).then(function(resp){
        if(resp.data.status) {
          $scope.releases = resp.data.data;
        }
        else {
          toaster.pop("error","get task error",resp.data.info);
        }
      },function(){

      });
    }
    if(selectPadder == "new" || selectPadder == "config"||selectPadder == "operation") {
      $http.post("/api/release/getconf",$scope.selectedService).then(function(resp){
        if(resp.data.status) {
          $scope.releaseConf = resp.data.data
        }
        else {
          toaster.pop("error","get conf error",resp.data.info);
        }
      },function(){

      });
    }

  }

  $scope.selectTask = function(release) {
    $scope.selectedTask = release;
    $scope.changepadder('operation')
  }
  $scope.returnUpper = function(idx) {
    $scope.filter[idx-1] = "";
    $scope.selectedService = undefined;
  }

  $scope.submitTask = function() {
    $scope.newReleaseTask.ReleaseConf = $scope.releaseConf;
    $scope.newReleaseTask.Service = $scope.selectedService;
    $http.post("/api/release/newtask",$scope.newReleaseTask).then(function(resp){
      if(resp.data.status) {
        $scope.releases.push(resp.data.data)
      }
      else {
        toaster.pop("error","new task error",resp.data.info);
      }
    },function(){

    });


  }

  
}]);
