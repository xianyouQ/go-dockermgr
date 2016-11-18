'use strict';

app.factory('authService', function ($http) {
 
    var userRole = []; // obtained from backend
    var userRoleRouteMap = {};
    var user;
 
    return {
 
        userHasRole: function (role) {
            if(user == undefined) {
                return undefined;
            }
            for (var j = 0; j < userRole.length; j++) {
                if (role == userRole[j]) {
                    return true;
                }
            }
            return false;
        },
 
        isUrlAccessibleForUser: function (route) {
            if(user == undefined) {
                return undefined;
            }
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
            return false;
        },
        returnUser: function () {
            return user
        }
    };
});