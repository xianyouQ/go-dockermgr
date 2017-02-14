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
                var authNameSplits = String(auth.ServiceAuth.Name).split(".");
                angular.forEach(authNameSplits,function(Split){
                    if(Split == role){
                        yes = true;
                        return false;
                    }
                });
                if(yes == true){
                    return false;
                }

            });
            return yes;
        },
 
        isUrlAccessibleForUser: function (route) {
            if($sessionStorage.user == undefined) {
                return undefined;
            }

            return true;
        },
        saveLastState: function(route){
            lastState = route;
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