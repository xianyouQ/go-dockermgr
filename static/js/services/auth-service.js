'use strict';

app.factory('authService', function ($http,$localStorage) {
 
    var userRole = []; // obtained from backend
    var userRoleRouteMap = {};
    $localStorage.$default({
        user:undefined
    });
 
    return {
 
        userHasRole: function (role) {
            if($localStorage.user == undefined) {
                return undefined;
            }
            /*
            for (var j = 0; j < userRole.length; j++) {
                if (role == userRole[j]) {
                    return true;
                }
            }*/
            return false;
        },
 
        isUrlAccessibleForUser: function (route) {
            if($localStorage.user == undefined) {
                return undefined;
            }
            /*
            for (var i = 0; i < userRole.length; i++) {
                var role = userRole[i];
                var validUrlsForRole = userRoleRouteMap[role];
                if (validUrlsForRole) {
                    for (var j = 0; j < validUrlsForRole.length; j++) {
                        if (validUrlsForRole[j] == route)
                            return true;
                    }
                }
            }
            */
            return true;
        },
        returnUser: function () {
            return $localStorage.user
        },
        login: function(loginuser) {
             console.log("add auth");
            $localStorage.user = loginuser;
        },
        logout: function() {
            console.log("removing auth");
            $localStorage.user = undefined;
        }
    };
});