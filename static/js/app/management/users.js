app.controller('ManageMentUsersCtrl', ['$scope', '$http', '$filter', function($scope, $http, $filter) {

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
  $scope.mainbuses = [] ;
  $scope.people = [] ;
  $scope.mainfilter = '';
  $scope.subfilter = '';
  

  $http.get('js/app/management/bussiness.json').then(function (resp) {
    $scope.subbuses = resp.data.subbuses;
    angular.forEach($scope.subbuses,function(item) {
      newitem ={name: item["mainbus"]};
      if(!$scope.mainbuses.contains(newitem)) {
        $scope.mainbuses.push(newitem);
      }
    })
    $scope.mainbus = $filter('orderBy')($scope.mainbuses, 'name')[0];
    $scope.mainbus.selected = true;
    $scope.mainfilter = $scope.mainbus.name;
  });

  $scope.selectMainBus = function(item){    
    angular.forEach($scope.mainbuses, function(item) {
      item.selected = false;
    });
    $scope.mainbus = item;
    $scope.mainbus.selected = true;
    $scope.mainfilter = item.name;
    console.log(item.name)
  };

  $scope.selectSubBus = function(item){    
    angular.forEach($scope.subbuses, function(item) {
      item.selected = false;
    });
    $scope.subbus = item;
    $scope.subbus.selected = true;
    $scope.subfilter = item.name;
    $http.get('js/app/management/people.json').then(function(resp) {
      $scope.people = resp.data.people ;
      $scope.person = $filter('orderBy')($scope.people, 'name')[0];
      $scope.person.selected = true ;
    })
  };
  $scope.selectPeople = function(item) {
    angular.forEach($scope.people,function(item) {
      item.selected = false;
    });
    $scope.person = item;
    $scope.person.selected = true;
  }
  $scope.createPerson = function() {

  }
}]);