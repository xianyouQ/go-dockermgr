'use strict';

app.factory('cacheStorage', function ($localStorage,$http) {

    return {
        get: function(url) {
           var vaule = $localStorage[url]
            if ( vaule === nul) {
                $http.get(url).then(function(resp){
                    if (resp.data.status ){
                        $localStorage[url] = resp.data.data
                        return resp.data.data
                    }
                });
            }
            else {

            }
        }
    }
});