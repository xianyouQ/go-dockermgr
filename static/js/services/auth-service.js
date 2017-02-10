'use strict';

app.factory('authService', function ($sessionStorage) {
 
    var lastState;
    return {
 
        userHasRole: function (role) {
            if($sessionStorage.auths == undefined) {
                return false;
            }
            var yes = false;
            angular.forEach($sessionStorage.auths,function(auth){
                if(String(auth.ServiceAuth.Name).startsWith(role)){
                    yes = true;
                    return false;
                } else if(String(auth.ServiceAuth.Name).endsWith(role)){
                    yes = true;
                    return false;
                }
            });
            return yes;
        },
 
        isUrlAccessibleForUser: function (route) {
            if($sessionStorage.user == undefined) {
                lastState = route;
                return undefined;
            }

            return true;
        },
        returnUser: function () {
            return $sessionStorage.user
        },
        login: function(data) {
            $sessionStorage.user = data.Username;
            $sessionStorage.auths = data.auth;
        },
        logout: function() {
            delete $sessionStorage.user;
            delete $sessionStorage.auths;
        },
        getLastState: function() {
            return lastState;
        },
        getAuths: function(){
            return $sessionStorage.auths;
        }
    };
});