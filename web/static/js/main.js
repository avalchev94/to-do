let activeModalPage = null;

$(document).ready(function(){
  initTopMenu();
  initSideMenu();

  initModalPages();
  initAddTaskPage();
})

function initTopMenu() {
  $("#add-task").click(function(e){
    activeModalPage = loadAddTaskPage();
  })
}

function initSideMenu() {
  let toggled = false
  $("#menu-toggle").click(function(e){
    let display = toggled ? "none" : "inline"
    let width = toggled ? "50px" : "200px"
    $("#left-bar").animate({width: width}, 400)
    $("#left-bar span").css("display", display)
    
    toggled = !toggled  
  })
}

function initModalPages() {
  $("#modal-cancel").click(function(){
    $(activeModalPage).fadeOut(200);
    activeModalPage = null;
  })
}

function alertBox(msg) {
  alert(msg)
}