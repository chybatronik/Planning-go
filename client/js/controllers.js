function DirectionSpisokCtrl($scope, $location, Direction, $filter) {
	//
	$scope.list = Direction.query();
	$scope.query_add = ""

	$scope.update_priority = function() {
		//
		for (var i = $scope.list.length - 1; i >= 0; i--) {
			var cur_id = $scope.list[i].Id
			//console.log(cur_id, $scope.list.length - i - 1);
			$scope.list[i].Priority = ($scope.list.length - i - 1) * 100	
			$scope.list[i].$save({id:$scope.list[i].Id})
			//console.log($scope.list[i].Id, $scope.list[i])
		};
	};

	$scope.Add = function() {
		//
		if($scope.query_add!=""){	
			Direction.create($scope.query_add, function(){ $scope.list = Direction.query(); });
		}
		$scope.query_add = ""
	};

	$scope.Del = function(id) {
		//
		Direction.rem({id:id}, function(){ $scope.list = Direction.query();})
	};

	$scope.change_work = function(id, isnotwork) {
		//
		console.log("id", id, isnotwork)
		var user = Direction.get({id:id}, function() {
			console.log("user", user, isnotwork, user.IsWork)
			user.IsWork = isnotwork
			console.log("user", user)
			user.$save({id:id});
			//$scope.list = Tasks_Direction.query({id:$routeParams.id});
		});	
	};

	$scope.change_priority = function(list_id, plus_or_minus) {
		//
		var list = $filter('orderBy')($scope.list,'Priority', true)
		console.log("id", list_id, list)

		var user = Direction.get({id:list_id}, function() {
			console.log("Direction.get", user, user.Priority)
			if(plus_or_minus == "+"){
				var index = -1
				var value = 0
				for (var i = list.length - 1; i >= 0; i--) {
					if (list[i].Id == list_id){
						index = i
					}
				};
				if (index > 0){
					var temp = Direction.get({id:list[index - 1].Id}, function() {
						//
						value = temp.Priority
						temp.Priority = user.Priority
						temp.$save({id:list[index - 1].Id},  function(){ 
							// 
							user.Priority = value
							user.$save({id:list_id},  function(){ 
								//
								$scope.list = Direction.query();
							});
						});
						console.log("temp:::", temp.Priority, list[index - 1].Id)

						
						console.log("user:::", user.Priority, list_id)
						
						
					});
				}
				
			}else{
				var index = -1
				var value = 0
				for (var i = list.length - 1; i >= 0; i--) {
					if (list[i].Id == list_id){
						index = i
					}
				};
				if (index > -1 && index != list.length-1){
					var temp = Direction.get({id:list[index + 1].Id}, function() {
						//
						value = temp.Priority
						temp.Priority = user.Priority
						temp.$save({id:list[index + 1].Id},  function(){ 
							// 
							user.Priority = value
							user.$save({id:list_id},  function(){ 
								//
								$scope.list = Direction.query();
							});
						});
					});
				}
			}
			
		});	

	};

}

function DirectionEditCtrl ($scope, $routeParams, $location, Direction) {
	// body...
	$scope.item = Direction.get({id:$routeParams.id})

	$scope.save = function() {
		//
		$scope.item.$save({id:$routeParams.id},  function(){ $location.path('/direction/task/'); })
	};

	$scope.remove_WhenWork = function(index) {
		$scope.item.WhenWork.splice(index, index+1)
	};

	$scope.Add_WhenWork = function() {
		//
		if($scope.item.WhenWork == null){
			$scope.item.WhenWork = [{Start:'9:00', Stop:'17:00'}]
		}else{
			$scope.item.WhenWork.push({Start:'50:9', Stop:'60:90'});
		}
		
	};
}

function TaskSpisokLabelCtrl ($scope, $routeParams, $location, Tasks_Label, Direction, Task) {
	// body...
	$scope.executed = false
	
	$scope.list = Tasks_Label.query({executed:""+$scope.executed+""});
	$scope.Direction_list = Direction.query();

	$scope.query_add = ""
	$scope.query_add_time = "1h"
	$scope.direction_id = -1

	$scope.show_executed = function() {
		//
		$scope.list = Tasks_Label.query({executed:""+$scope.executed+""});
	};
	
	$scope.completed = function(id, completed) {
		//
		console.log("id", id, completed)
		var user = Task.get({id:id}, function() {
			console.log("user", user, completed, user.IsDone)
			user.IsDone = completed
			console.log("user", user)
			user.$save({id:id});
			//$scope.list = Tasks_Direction.query({id:$routeParams.id});
		});	
	};


	$scope.set_label = function(id) {
		//
		var user = Task.get({id:id}, function() {
			if(user.Label){
				user.Label = false;
			}else{
				user.Label = true;
			}
			
			user.$save({id:id});
			$scope.list = Tasks_Label.query();
		});		
	};

	$scope.Del = function(id) {
		//
		Task.rem({id:id}, function(){ $scope.list = Tasks_Label.query({executed:""+$scope.executed+""});})
	};
}

function TaskSpisokAllCtrl ($scope, $routeParams, $location, Task, Direction) {
	// body...
	$scope.executed = false

	$scope.list = Task.query(function() {
		console.log("asd")
		$('.toggle-button').toggleButtons({
		    onChange: function ($el, status, e) {
		        console.log($el, status, e);
		    }
		})
	});
	$scope.Direction_list = Direction.query();

	$scope.query_add = ""
	$scope.direction_id = 0

	$scope.query_add_time = "1h"

	

	$scope.show_executed = function() {
		//
		$scope.list = Task.query({executed:""+$scope.executed+""});
	};

	$scope.completed = function(id, completed) {
		//
		console.log("id", id, completed)
		var user = Task.get({id:id}, function() {
			console.log("user", user, completed, user.IsDone)
			user.IsDone = completed
			console.log("user", user)
			user.$save({id:id});
			//$scope.list = Tasks_Direction.query({id:$routeParams.id});
		});	
	};

	$scope.set_label = function(id) {
		//
		var user = Task.get({id:id}, function() {
			if(user.Label){
				user.Label = false;
			}else{
				user.Label = true;
			}
			
			user.$save({id:id});
			$scope.list = Task.query();
		});		
	};

	$scope.Del = function(id) {
		//
		Task.rem({id:id}, function(){ $scope.list = Task.query();})
	};

}

function TaskSpisokDirectionCtrl ($scope, $routeParams, $location, Direction, Task, $filter, $q) {
	// body...
	$scope.executed = false
	//$scope.list = Tasks_Direction.query({id:$routeParams.id, executed:""+$scope.executed+""});
	$scope.list = Task.query({direction:$routeParams.id, executed:""+$scope.executed+""});
	$scope.Direction_list = Direction.query({}, function() {
		for (var i = 0; i < $scope.Direction_list.length; i++) {
			if ($routeParams.id == $scope.Direction_list[i].Id){
				$scope.cur_Direction = $scope.Direction_list[i]
			}
		};
		
	});
	$scope.direction_id = $routeParams.id
	
	$scope.query_add_time = "0.5h"

	$scope.update_priority = function() {
		//
		for (var i = $scope.list.length - 1; i >= 0; i--) {
			var cur_id = $scope.list[i].Id
			//console.log(cur_id, $scope.list.length - i - 1);
			$scope.list[i].PriorityTask = $scope.list.length - i - 1	
			$scope.list[i].$save({id:$scope.list[i].Id})
			//console.log($scope.list[i].Id, $scope.list[i])
		};
	};

	$scope.Time_plus = function() {
		//
		$scope.query_add_time = (parseFloat($scope.query_add_time.split("h")[0]) + 0.5) + "h"
	};

	$scope.Time_minus = function() {
		//
		if (parseFloat($scope.query_add_time.split("h")[0]) > 0.5){
			$scope.query_add_time = (parseFloat($scope.query_add_time.split("h")[0]) - 0.5) + "h"
		}
	};

	$scope.completed = function(id, completed) {
		//
		console.log("id", id, completed)
		var user = Task.get({id:id}, function() {
			console.log("user", user, completed, user.IsDone)
			user.IsDone = completed
			console.log("user", user)
			user.$save({id:id});
			//$scope.list = Tasks_Direction.query({id:$routeParams.id});
		});	
	};

	$scope.set_label = function(id) {
		//
		var user = Task.get({id:id}, function() {
			if(user.Label){
				user.Label = false;
			}else{
				user.Label = true;
			}
			
			user.$save({id:id}, function() {
				$scope.list = Task.query({direction:$routeParams.id, executed:""+$scope.executed+""});
			});
			
		});		
	};

	$scope.Add = function() {
		//
		if($scope.query_add!=""){	
			var id = "0";
			if (typeof $routeParams.id === "undefined"){
				id = "0"
			}
			else{
				id = $routeParams.id
			}
			var list = $filter('orderBy')($scope.list,'Priority', false)
			
			if (list.length > 0){
				Priority = parseInt(list[0].PriorityTask) - 1
			}else{
				Priority = 10
			}
			console.log(list[0], Priority)
			Task.create({name:$scope.query_add, 
							direction_id:parseInt(id), 
							duration:$scope.query_add_time, 
							label:false,
							prioritytask:Priority}, function(){ $scope.list = Task.query({direction:$routeParams.id, executed:""+$scope.executed+""}); });
			$scope.query_add = ""
			$scope.query_add_time = "0.5h"
		}
	};

	$scope.Add_fast = function() {
		//
		if($scope.query_add!=""){	
			var id = "0";
			if (typeof $routeParams.id === "undefined"){
				id = "0"
			}
			else{
				id = $routeParams.id
			}
			var list = $filter('orderBy')($scope.list,'Priority', true)
			
			if (list.length > 0){
				Priority = parseInt(list[0].PriorityTask) + 1
			}else{
				Priority = 10
			}
			console.log(list[0], Priority)
			Task.create({name:$scope.query_add, 
							direction_id:parseInt(id), 
							duration:$scope.query_add_time, 
							label:false,
							prioritytask:Priority}, function(){ $scope.list = Task.query({direction:$routeParams.id, executed:""+$scope.executed+""}); });
			$scope.query_add = ""
			$scope.query_add_time = "0.5h"
		}
	};
	$scope.show_executed = function() {
		//
		$scope.list = Task.query({direction:$routeParams.id, executed:""+$scope.executed+""});
	};

	$scope.Del = function(id) {
		//
		Task.rem({id:id}, function(){ $scope.list = Task.query({direction:$routeParams.id, executed:""+$scope.executed+""});})
	};
}

function TaskEditCtrl ($scope, $location, Task, Direction, $routeParams) {
	// body...
	$scope.item = Task.get({id:$routeParams.id},
		function(){
			if($scope.item.Repeat.WhatWeekyDayRestore!=null ){
				for (var i = $scope.item.Repeat.WhatWeekyDayRestore.length - 1; i >= 0; i--) {
					$scope.item.Repeat.WhatWeekyDayRestore[i]
					$("option[value='" + $scope.item.Repeat.WhatWeekyDayRestore[i]+"']").attr("selected","selected")
					console.log($scope.item.Repeat.WhatWeekyDayRestore[i])
				};
				
			}
			$(".chzn-select").chosen();
		});
	$scope.direction_list = Direction.query();

	$scope.cur = ""

	
	$scope.change_direction = function() {
		//
		$scope.item.Direction_Id = parseInt($scope.cur)
	};

	$scope.save = function() {
		//
		$scope.item.$save({id:$routeParams.id},  
				function(){ $location.path('/task/direction/' + $scope.item.Direction_Id); })
	};
}

function ScheduleCtrl ($scope, $routeParams, Schedule, Statistic) {
	// body...
	$scope.list = Schedule.query();
	$scope.spisok_stat = Statistic.query();

	$scope.completed = function(id, completed) {
		//
		console.log("id", id, completed)
		var user = Schedule.get({id:id}, function() {
			console.log("user", user, completed, user.IsDone)
			user.IsDone = completed
			console.log("user", user)
			user.$save({id:id});
			//$scope.list = Tasks_Direction.query({id:$routeParams.id});
			
		});	
	};
}