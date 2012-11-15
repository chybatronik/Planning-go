angular.module('admin', ['ui', 'admin.services'])
  .config(['$routeProvider', function($routeProvider) {
    $routeProvider
      .when('/direction', 
      			{templateUrl: 'client/views/direction/spisok.html', controller: DirectionSpisokCtrl})
      .when('/direction/edit/:id', 
      			{templateUrl: 'client/views/direction/edit.html', controller: DirectionEditCtrl})
      .when('/task/direction/:id', 
      			{templateUrl: 'client/views/task/spisok.html', controller: TaskSpisokDirectionCtrl})
      .when('/task', 
            {templateUrl: 'client/views/task/spisok_all.html', controller: TaskSpisokAllCtrl})
      .when('/task/label', 
            {templateUrl: 'client/views/task/spisok_label.html', controller: TaskSpisokLabelCtrl})
      .when('/task/edit/:id', 
      			{templateUrl: 'client/views/task/edit.html', controller: TaskEditCtrl})
      .when('/schedule', 
            {templateUrl: 'client/views/schedule/spisok.html', controller: ScheduleCtrl})
      .otherwise({redirectTo: '/schedule'});
  },
]);