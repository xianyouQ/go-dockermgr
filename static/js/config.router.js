'use strict';

/**
 * Config for the router
 */
angular.module('app')
  .run(
    [          '$rootScope', '$state', '$stateParams',
      function ($rootScope,   $state,   $stateParams) {
          $rootScope.$state = $state;
          $rootScope.$stateParams = $stateParams;
        
      }
    ]
  )
  .config(
    [          '$stateProvider', '$urlRouterProvider',
      function ($stateProvider,   $urlRouterProvider) {
          $urlRouterProvider
              .otherwise('/docker/dashboard');
          $stateProvider
              .state('lockme', {
                  url: '/lockme',
                  templateUrl: 'tpl/page_lockme.html'
              })
              .state('access', {
                  url: '/access',
                  template: '<div ui-view class="fade-in-right-big smooth"></div>'
              })
              .state('access.signin', {
                  url: '/signin',
                  templateUrl: 'tpl/page_signin.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/controllers/signin.js'] );
                      }]
                  }
              })
              .state('access.signup', {
                  url: '/signup',
                  templateUrl: 'tpl/page_signup.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/controllers/signup.js'] );
                      }]
                  }
              })
              .state('access.forgotpwd', {
                  url: '/forgotpwd',
                  templateUrl: 'tpl/page_forgotpwd.html'
              })
              .state('access.error', {
                  url: '/error/:errorCode',
                  templateUrl: 'tpl/page_error.html',
                  controller: function($scope,$stateParams) {
                    $scope.errorCode = $stateParams.errorCode;
                  }
              })
              // fullCalendar
                .state('management', {
                  url: '/management',
                  abstract: true,
                  templateUrl: 'tpl/layout.html'
              })
                .state('management.users', {
                  url: '/users',
                  templateUrl: 'tpl/management_users.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/management/users.js'] );
                      }]
                  }
              })
                .state('management.idcs', {
                  url: '/idcs',
                  templateUrl: 'tpl/management_idcs.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/management/idcs.js'] );
                      }]
                  }
              })
                .state('management.services', {
                  url: '/services',
                  templateUrl: 'tpl/management_services.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/management/services.js'] );
                      }]
                  }
              })
                .state('management.auth', {
                  url: '/auth',
                  templateUrl: 'tpl/management_auth.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/management/auth.js'] );
                      }]
                  }
              })
                .state('docker', {
                  url: '/docker',
                  abstract: true,
                  templateUrl: 'tpl/layout.html'
              })
                .state('docker.dashboard', {
                  url: '/dashboard',
                  templateUrl: 'tpl/docker_dashboard.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/docker/dashboard.js'] );
                      }]
                  }
              })
                .state('docker.containers', {
                  url: '/services',
                  templateUrl: 'tpl/docker_containers.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/docker/containers.js'] );
                      }]
                  }
              })
                .state('release', {
                  url: '/release',
                  abstract: true,
                  templateUrl: 'tpl/layout.html'
              })
                .state('release.service', {
                  url: '/services',
                  templateUrl: 'tpl/release_service.html',
                  resolve: {
                      deps: ['uiLoad',
                        function( uiLoad ){
                          return uiLoad.load( ['js/app/release/service.js'] );
                      }]
                  }
              })
      }
    ]
  );