app.controller('ReleaseServiceCtrl', ['$scope', '$http', '$filter','$modal','$interval','$q','toaster','myCache',function($scope, $http, $filter,$modal,$interval,$q,toaster,myCache) {
  function isObjectValueEqual(a, b) {
   if(a.Id === b.Id){
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
  $scope.refreshStatus = undefined;


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
    $http.post("/api/release/conf?serviceId="+$scope.selectedService.Id,$scope.releaseConf).then(function(resp){
      if(resp.data.status) {
        $scope.releaseConf = resp.data.data;
      }
      else {
        toaster.pop("error","update conf error",resp.data.info);
      }
    },function(){
      
    })
  }

  $scope.stopRefreshStatus = function() {
    console.log("stop refresh");
    if(angular.isDefined($scope.refreshStatus)){
      $interval.cancel($scope.refreshStatus)
      $scope.refreshStatus = undefined;
    }
  }
  $scope.changepadder = function(selectPadder) {
    $scope.padderSelect = selectPadder;
    if(selectPadder == "list") {
      $scope.stopRefreshStatus();
      $http.get("/api/release/task?serviceId="+$scope.selectedService.Id).then(function(resp){
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
      $scope.stopRefreshStatus();
      $http.get("/api/release/getconf?serviceId="+$scope.selectedService.Id).then(function(resp){
        if(resp.data.status) {
          $scope.releaseConf = resp.data.data
          angular.forEach($scope.releaseConf.ReleaseIdc,function(outIdc){
            angular.forEach($scope.idcs,function(idc){
               if(idc.Id == outIdc.Id) {
                 idc.enableRelease = true;
               }
            });
          });
          
        }
        else {
          toaster.pop("error","get conf error",resp.data.info);
        }
      },function(){

      });
    }
    if(selectPadder == "detail") {
      $scope.refreshStatus = $interval(function() {
        $scope.selectedTask.ReleaseResult = undefined;
        $http.get("/api/release/status?serviceId="+$scope.selectedService.Id+"&taskId="+$scope.selectedTask.Id).then(function(resp){
          if(resp.data.status) {
            //$scope.selectedTask = resp.data.data;
            try{
              $scope.selectedTask.ReleaseResult = JSON.parse(resp.data.data.ReleaseResult);
            } catch(e) {

            }
            $scope.selectedTask.TaskStatus = resp.data.data.TaskStatus;
            if ($scope.selectedTask.TaskStatus >= 4) {
              console.log("stop refresh before");
              $scope.stopRefreshStatus();
            }
          }
          else{
            toaster.pop("error","get detail error",resp.data.info);
          }
        });
      },1000);
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


  $scope.abandonTask = function(taskId) {
    
  }
  $scope.submitTask = function() {
    $scope.newReleaseTask.ReleaseConf = $scope.releaseConf;
    $scope.newReleaseTask.Service = $scope.selectedService;
    $http.post("/api/release/newtask?serviceId="+ $scope.selectedService.Id,$scope.newReleaseTask).then(function(resp){
      if(resp.data.status) {
        $scope.releases.push(resp.data.data)
      }
      else {
        toaster.pop("error","new task error",resp.data.info);
      }
    },function(){

    });
  }
  $scope.reviewTask = function() {
    $http.get("/api/release/review?serviceId="+$scope.selectedService.Id+"&taskId="+$scope.selectedTask.Id).then(function(resp){
      if(resp.data.status) {
        $scope.selectedTask.TaskStatus = resp.data.data.TaskStatus;
        $scope.selectedTask.ReviewUser = resp.data.data.ReviewUser;
        $scope.selectedTask.ReviewTime = resp.data.data.ReviewTime;
      }
      else {
        toaster.pop("error","review task error",resp.data.info);
      }
    },function(){

    });
  };
  $scope.operationTask = function() {
    //$scope.selectedTask.ReleaseConf = $scope.releaseConf;
    $http.get("/api/release/operate?serviceId="+$scope.selectedService.Id+"&taskId="+$scope.selectedTask.Id).then(function(resp){
      if(resp.data.status) {
        $scope.releases.remove($scope.selectedTask);
        $scope.selectedTask = resp.data.data;
        $scope.releases.push(resp.data.data);
        $scope.changepadder("detail");
      }
      else {
        toaster.pop("error","operate task error",resp.data.info);
      }
    },function(){
      
    });
  }
  
}]);
