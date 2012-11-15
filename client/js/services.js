angular.module('admin.services', ['ngResource'])
  .factory('Direction', function($resource){
    return $resource('direction/:id', {}, {
      create: {method:'PUT'},
      rem: {method:'DELETE'},
      saveData: {method:'POST'}
    })})
  .factory('Task', function($resource){
    return $resource('task/:id?direction=:direction;executed=:executed', {}, {
      get: {method:'GET', isArray: false},
      create: {method:'PUT'},
      rem: {method:'DELETE'},
      saveData: {method:'POST'}
    });
  })
  .factory('Tasks_Direction', function($resource){
    return $resource('task/direction/:id?executed=:executed', {}, {
      //get: {method:'GET', isArray: true},
      create: {method:'PUT'},
      rem: {method:'DELETE'},
      saveData: {method:'POST'}
    });
  })
.factory('Tasks_Label', function($resource){
    return $resource('task/label/:id?executed=:executed', {}, {
      create: {method:'PUT'},
      rem: {method:'DELETE'},
      saveData: {method:'POST'}
    });
  })
  .factory('Schedule', function($resource){
    return $resource('schedule/:id', {}, {
      saveData: {method:'POST'}
    });
  })
  .factory('Statistic', function($resource){
    return $resource('statistic/:id', {}, {
      saveData: {method:'POST'}
    });
  });