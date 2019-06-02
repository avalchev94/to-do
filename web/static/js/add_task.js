let tagField = null;

function initAddTaskPage() {
  // labels field, using tagify plugin.
  tagField = $('.add-task input.labels'); 
  tagField.tagify({
    duplicate: false
  })

  // side buttons for scheduling.
  $('.add-task div.buttons').find('.accordion-toggle').click(function(){
    $(this).next().slideToggle('fast');
    $(this).toggleClass('active');
    
    // Hide the other panels
    $('.accordion-content').not($(this).next()).slideUp('fast');
    $('.add-task div.buttons').find('.active').not(this).removeClass("active");
  });

  // change data according to the choosen repeatance type.
  $('.add-task select.repeat-frequency').change(function(){
    $(`.repeatance.${this.value}`).slideDown('fast');
    $('.repeatance').not(`.${this.value}`).slideUp('fast');
  })

  // on add button
  $('.add-task button.add').click(function(){
    let request = addTask();
    request.success(function(){
      $('#modal-cancel').trigger('click');
    })
    request.error(function(data){
      alertBox(data.responseJSON.error)
    })
  })
    
  // on cancel button
  $('.add-task button.cancel').click(function(){
    $('#modal-cancel').trigger('click');
  })
}

function loadAddTaskPage() {
  $("#add-task-page").fadeIn(200);
  
  $.ajax({
    type: "GET",
    url:  "/labels"
  }).success(function(data) {
    tagField.data('tagify').settings.whitelist = data
  })

  return "#add-task-page"
}

function addTask() {
  let getTask = function() {
    let getLabels = function() {
      let labels = []
      $('.add-task input.labels').data('tagify').value.forEach(function(t) {
        labels.push(t.value)
      })
      return labels
    }

    return { 'task': {
      title: $('.add-task input.title').val(),
      body: $('.add-task textarea.body').val(),
      labels: getLabels()
    }}
  }

  let getSchedule = function() {
    const task_type = $('.add-task div.buttons').find('.active').attr('name');
    const content = $(`.add-task button[name=${task_type}]`).next();

    switch (task_type) {
    case 'due_date':
    case 'to_date':
      return {'schedule': {
        type: task_type,
        date: content.children('input[type="date"]').val(),
        time: content.children('input[type="time"]').val()
      }}
    case 'repeat':
      var repeatance = new Object();
      repeatance.type = content.children('select').val();
      
      if (repeatance.type != null) {
        repeatance.hour = content.children('input[type="time"]').val();
        repeatance.days = [];
        content.children(`.repeatance.${repeatance.type}`)
          .find(':checkbox:checked')
          .each(function(i){
            repeatance.days[i] = Number($(this).attr('name'))
          })
      }
      return {'repeatance': repeatance }
    }

    return {'schedule': {type: 'unscheduled'}}
  };

  let task = getTask();
  $.extend(task, getSchedule());

  return $.ajax({
    type: "POST",
    url: "/task",
    data: JSON.stringify(task, function(key, value) {
      if (value == '' || value == null) {
        return undefined
      }
      return value
    })
  })
}