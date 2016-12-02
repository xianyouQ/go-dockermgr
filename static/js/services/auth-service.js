'use strict';

app.factory('authService', function ($sessionStorage) {
 
    var userRole = []; // obtained from backend
    var userRoleRouteMap = {};
    var lastState;
    return {
 
        userHasRole: function (role) {
            if($sessionStorage.user == undefined) {
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
            if($sessionStorage.user == undefined) {
                lastState = route;
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
            return $sessionStorage.user
        },
        login: function(loginuser) {
            $sessionStorage.user = loginuser;
        },
        logout: function() {
            delete $sessionStorage.user;
        },
        getLastState: function() {
            return lastState;
        },
        getServiceAuth: function() {
        }
    };
});